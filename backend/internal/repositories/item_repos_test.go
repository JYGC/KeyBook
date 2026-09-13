package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

func TestItemRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)
	itemRepository := repositories.NewItemRepository(app)

	created, err := itemRepository.CreateItem("Front Door Key", "Main entrance key")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Name != "Front Door Key" {
		t.Errorf("Name: want Front Door Key, got %s", created.Name)
	}
	if created.Description != "Main entrance key" {
		t.Errorf("Description: want Main entrance key, got %s", created.Description)
	}

	got, err := itemRepository.GetItemById(created.Id)
	if err != nil {
		t.Fatalf("GetItemById: %v", err)
	}
	if got.Name != "Front Door Key" {
		t.Errorf("Name: want Front Door Key, got %s", got.Name)
	}

	if err := itemRepository.UpdateItem(created.Id, "Gate Key", "Side gate key"); err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}
	updated, err := itemRepository.GetItemById(created.Id)
	if err != nil {
		t.Fatalf("GetItemById after update: %v", err)
	}
	if updated.Name != "Gate Key" {
		t.Errorf("after update Name: want Gate Key, got %s", updated.Name)
	}

	if err := itemRepository.DeleteItem(created.Id); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}
	if _, err := itemRepository.GetItemById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestPropertyItemRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	item := createRecordBypassingAccessRules(t, app, "items", map[string]any{"name": "Key", "description": "test"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	propertyItemRepository := repositories.NewPropertyItemRepository(app)

	created, err := propertyItemRepository.CreatePropertyItem(item.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyItem: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Item != item.GetId() {
		t.Errorf("Item: want %s, got %s", item.GetId(), created.Item)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	got, err := propertyItemRepository.GetPropertyItemById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyItemById: %v", err)
	}
	if got.Item != item.GetId() {
		t.Errorf("Item: want %s, got %s", item.GetId(), got.Item)
	}

	items, err := propertyItemRepository.GetPropertyItemsByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetPropertyItemsByPropertyId: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1 item, got %d", len(items))
	}
	if items[0].Id != created.Id {
		t.Errorf("filtered item ID: want %s, got %s", created.Id, items[0].Id)
	}

	if err := propertyItemRepository.DeletePropertyItem(created.Id); err != nil {
		t.Fatalf("DeletePropertyItem: %v", err)
	}
	if _, err := propertyItemRepository.GetPropertyItemById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestPersonItemRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Alice", "DOB": "1990-01-01 00:00:00.000Z"})
	item := createRecordBypassingAccessRules(t, app, "items", map[string]any{"name": "Key", "description": "test"})

	personItemRepository := repositories.NewPersonItemRepository(app)

	created, err := personItemRepository.CreatePersonItem(person.GetId(), item.GetId())
	if err != nil {
		t.Fatalf("CreatePersonItem: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.Item != item.GetId() {
		t.Errorf("Item: want %s, got %s", item.GetId(), created.Item)
	}

	got, err := personItemRepository.GetPersonItemById(created.Id)
	if err != nil {
		t.Fatalf("GetPersonItemById: %v", err)
	}
	if got.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), got.Person)
	}

	pItems, err := personItemRepository.GetPersonItemsByPersonId(person.GetId())
	if err != nil {
		t.Fatalf("GetPersonItemsByPersonId: %v", err)
	}
	if len(pItems) != 1 {
		t.Fatalf("want 1 item, got %d", len(pItems))
	}

	if err := personItemRepository.DeletePersonItem(created.Id); err != nil {
		t.Fatalf("DeletePersonItem: %v", err)
	}
	if _, err := personItemRepository.GetPersonItemById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestEntryDeviceRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	item := createRecordBypassingAccessRules(t, app, "items", map[string]any{"name": "Key", "description": "test"})

	entryDeviceRepository := repositories.NewEntryDeviceRepository(app)

	created, err := entryDeviceRepository.CreateEntryDevice(item.GetId(), "Key", "KEY-001", "None")
	if err != nil {
		t.Fatalf("CreateEntryDevice: %v", err)
	}
	if created.Item != item.GetId() {
		t.Errorf("Item: want %s, got %s", item.GetId(), created.Item)
	}
	if created.DeviceType != "Key" {
		t.Errorf("DeviceType: want Key, got %s", created.DeviceType)
	}
	if created.Identifier != "KEY-001" {
		t.Errorf("Identifier: want KEY-001, got %s", created.Identifier)
	}
	if created.DefunctReason != "None" {
		t.Errorf("DefunctReason: want None, got %s", created.DefunctReason)
	}

	got, err := entryDeviceRepository.GetEntryDeviceById(created.Id)
	if err != nil {
		t.Fatalf("GetEntryDeviceById: %v", err)
	}
	if got.DeviceType != "Key" {
		t.Errorf("DeviceType: want Key, got %s", got.DeviceType)
	}

	if err := entryDeviceRepository.UpdateEntryDevice(created.Id, "Fob", "FOB-001", "Lost"); err != nil {
		t.Fatalf("UpdateEntryDevice: %v", err)
	}
	updated, _ := entryDeviceRepository.GetEntryDeviceById(created.Id)
	if updated.DeviceType != "Fob" {
		t.Errorf("after update DeviceType: want Fob, got %s", updated.DeviceType)
	}
	if updated.DefunctReason != "Lost" {
		t.Errorf("after update DefunctReason: want Lost, got %s", updated.DefunctReason)
	}

	if err := entryDeviceRepository.DeleteEntryDevice(created.Id); err != nil {
		t.Fatalf("DeleteEntryDevice: %v", err)
	}
	if _, err := entryDeviceRepository.GetEntryDeviceById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
