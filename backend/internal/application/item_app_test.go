package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestItemApplicationService_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)
	itemApplicationService := application.NewItemApplicationService(
		services.NewItemService(),
		repositories.NewItemRepository(app),
		repositories.NewEntryDeviceRepository(app),
		repositories.NewPropertyItemRepository(app),
		repositories.NewPersonItemRepository(app),
	)

	created, err := itemApplicationService.CreateItem("Front Door Key", "Main entrance")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	if created.Name != "Front Door Key" {
		t.Errorf("Name: want Front Door Key, got %s", created.Name)
	}

	if _, err := itemApplicationService.CreateItem("", "desc"); err == nil {
		t.Error("expected error for empty name, got nil")
	}

	if err := itemApplicationService.UpdateItem(created.Id, "Gate Key", "Side gate"); err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}

	if err := itemApplicationService.UpdateItem(created.Id, "", "desc"); err == nil {
		t.Error("expected error for empty name on update, got nil")
	}

	if err := itemApplicationService.DeleteItem(created.Id); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}
}

func TestItemApplicationService_EntryDevice(t *testing.T) {
	app := newTestAppWithMigrations(t)
	itemApplicationService := application.NewItemApplicationService(
		services.NewItemService(),
		repositories.NewItemRepository(app),
		repositories.NewEntryDeviceRepository(app),
		repositories.NewPropertyItemRepository(app),
		repositories.NewPersonItemRepository(app),
	)

	item := createRecordBypassingAccessRules(t, app, "items", map[string]any{"name": "Key", "description": "test"})

	entryDevice, err := itemApplicationService.CreateEntryDevice(item.GetId(), "Key", "KEY-001", "None")
	if err != nil {
		t.Fatalf("CreateEntryDevice: %v", err)
	}
	if entryDevice.DefunctReason != "None" {
		t.Errorf("DefunctReason: want None, got %s", entryDevice.DefunctReason)
	}

	if err := itemApplicationService.UpdateEntryDevice(entryDevice.Id, "Key", "KEY-001", "Lost"); err != nil {
		t.Fatalf("UpdateEntryDevice active→defunct: %v", err)
	}

	if err := itemApplicationService.UpdateEntryDevice(entryDevice.Id, "Key", "KEY-001", "None"); err == nil {
		t.Error("expected error for reactivation, got nil")
	}

	if err := itemApplicationService.DeleteEntryDevice(entryDevice.Id); err != nil {
		t.Fatalf("DeleteEntryDevice: %v", err)
	}
}

func TestItemApplicationService_Associations(t *testing.T) {
	app := newTestAppWithMigrations(t)
	itemApplicationService := application.NewItemApplicationService(
		services.NewItemService(),
		repositories.NewItemRepository(app),
		repositories.NewEntryDeviceRepository(app),
		repositories.NewPropertyItemRepository(app),
		repositories.NewPersonItemRepository(app),
	)

	item := createRecordBypassingAccessRules(t, app, "items", map[string]any{"name": "Key", "description": "test"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})
	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Alice", "DOB": "1990-01-01 00:00:00.000Z"})

	pi, err := itemApplicationService.AddPropertyItem(item.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("AddPropertyItem: %v", err)
	}
	if err := itemApplicationService.RemovePropertyItem(pi.Id); err != nil {
		t.Fatalf("RemovePropertyItem: %v", err)
	}

	personItem, err := itemApplicationService.AddPersonItem(person.GetId(), item.GetId())
	if err != nil {
		t.Fatalf("AddPersonItem: %v", err)
	}
	if err := itemApplicationService.RemovePersonItem(personItem.Id); err != nil {
		t.Fatalf("RemovePersonItem: %v", err)
	}
}
