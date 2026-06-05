package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestItemApplicationService_CRUD(t *testing.T) {
	app := newApp(t)
	svc := application.NewItemApplicationService(
		services.NewItemService(),
		repositories.NewItemRepository(app),
		repositories.NewEntryDeviceRepository(app),
		repositories.NewPropertyItemRepository(app),
		repositories.NewPersonItemRepository(app),
	)

	// Create — valid
	created, err := svc.CreateItem("Front Door Key", "Main entrance")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	if created.Name != "Front Door Key" {
		t.Errorf("Name: want Front Door Key, got %s", created.Name)
	}

	// Create — validation failure
	if _, err := svc.CreateItem("", "desc"); err == nil {
		t.Error("expected error for empty name, got nil")
	}

	// Update — valid
	if err := svc.UpdateItem(created.Id, "Gate Key", "Side gate"); err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}

	// Update — validation failure
	if err := svc.UpdateItem(created.Id, "", "desc"); err == nil {
		t.Error("expected error for empty name on update, got nil")
	}

	// Delete
	if err := svc.DeleteItem(created.Id); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}
}

func TestItemApplicationService_EntryDevice(t *testing.T) {
	app := newApp(t)
	svc := application.NewItemApplicationService(
		services.NewItemService(),
		repositories.NewItemRepository(app),
		repositories.NewEntryDeviceRepository(app),
		repositories.NewPropertyItemRepository(app),
		repositories.NewPersonItemRepository(app),
	)

	item := createRecord(t, app, "items", map[string]any{"name": "Key", "description": "test"})

	// Create entry device (active — "None")
	ed, err := svc.CreateEntryDevice(item.GetId(), "Key", "KEY-001", "None")
	if err != nil {
		t.Fatalf("CreateEntryDevice: %v", err)
	}
	if ed.DefunctReason != "None" {
		t.Errorf("DefunctReason: want None, got %s", ed.DefunctReason)
	}

	// Update — valid transition: active → defunct
	if err := svc.UpdateEntryDevice(ed.Id, "Key", "KEY-001", "Lost"); err != nil {
		t.Fatalf("UpdateEntryDevice active→defunct: %v", err)
	}

	// Update — invalid transition: defunct → active (reactivation)
	if err := svc.UpdateEntryDevice(ed.Id, "Key", "KEY-001", "None"); err == nil {
		t.Error("expected error for reactivation, got nil")
	}

	// Delete entry device
	if err := svc.DeleteEntryDevice(ed.Id); err != nil {
		t.Fatalf("DeleteEntryDevice: %v", err)
	}
}

func TestItemApplicationService_Associations(t *testing.T) {
	app := newApp(t)
	svc := application.NewItemApplicationService(
		services.NewItemService(),
		repositories.NewItemRepository(app),
		repositories.NewEntryDeviceRepository(app),
		repositories.NewPropertyItemRepository(app),
		repositories.NewPersonItemRepository(app),
	)

	item := createRecord(t, app, "items", map[string]any{"name": "Key", "description": "test"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})
	person := createRecord(t, app, "persons", map[string]any{"name": "Alice", "DOB": "1990-01-01 00:00:00.000Z"})

	// Add and remove property item
	pi, err := svc.AddPropertyItem(item.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("AddPropertyItem: %v", err)
	}
	if err := svc.RemovePropertyItem(pi.Id); err != nil {
		t.Fatalf("RemovePropertyItem: %v", err)
	}

	// Add and remove person item
	psi, err := svc.AddPersonItem(person.GetId(), item.GetId())
	if err != nil {
		t.Fatalf("AddPersonItem: %v", err)
	}
	if err := svc.RemovePersonItem(psi.Id); err != nil {
		t.Fatalf("RemovePersonItem: %v", err)
	}
}
