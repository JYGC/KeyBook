package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

// ── ItemRepository ────────────────────────────────────────────────────────────

func TestItemRepository_CRUD(t *testing.T) {
	app := newApp(t)
	repo := repositories.NewItemRepository(app.PocketBase)

	// Create
	created, err := repo.CreateItem("Front Door Key", "Main entrance key")
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

	// Get by ID
	got, err := repo.GetItemById(created.Id)
	if err != nil {
		t.Fatalf("GetItemById: %v", err)
	}
	if got.Name != "Front Door Key" {
		t.Errorf("Name: want Front Door Key, got %s", got.Name)
	}

	// Update
	if err := repo.UpdateItem(created.Id, "Gate Key", "Side gate key"); err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}
	updated, err := repo.GetItemById(created.Id)
	if err != nil {
		t.Fatalf("GetItemById after update: %v", err)
	}
	if updated.Name != "Gate Key" {
		t.Errorf("after update Name: want Gate Key, got %s", updated.Name)
	}

	// Delete
	if err := repo.DeleteItem(created.Id); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}
	if _, err := repo.GetItemById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── PropertyItemRepository ────────────────────────────────────────────────────

func TestPropertyItemRepository_CRUD(t *testing.T) {
	app := newApp(t)

	item := createRecord(t, app, "items", map[string]any{"name": "Key", "description": "test"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	repo := repositories.NewPropertyItemRepository(app.PocketBase)

	// Create
	created, err := repo.CreatePropertyItem(item.GetId(), property.GetId())
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

	// Get by ID
	got, err := repo.GetPropertyItemById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyItemById: %v", err)
	}
	if got.Item != item.GetId() {
		t.Errorf("Item: want %s, got %s", item.GetId(), got.Item)
	}

	// Get by property ID
	items, err := repo.GetPropertyItemsByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetPropertyItemsByPropertyId: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1 item, got %d", len(items))
	}
	if items[0].Id != created.Id {
		t.Errorf("filtered item ID: want %s, got %s", created.Id, items[0].Id)
	}

	// Delete
	if err := repo.DeletePropertyItem(created.Id); err != nil {
		t.Fatalf("DeletePropertyItem: %v", err)
	}
	if _, err := repo.GetPropertyItemById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── PersonItemRepository ──────────────────────────────────────────────────────

func TestPersonItemRepository_CRUD(t *testing.T) {
	app := newApp(t)

	person := createRecord(t, app, "persons", map[string]any{"name": "Alice", "DOB": "1990-01-01 00:00:00.000Z"})
	item := createRecord(t, app, "items", map[string]any{"name": "Key", "description": "test"})

	repo := repositories.NewPersonItemRepository(app.PocketBase)

	// Create
	created, err := repo.CreatePersonItem(person.GetId(), item.GetId())
	if err != nil {
		t.Fatalf("CreatePersonItem: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.Item != item.GetId() {
		t.Errorf("Item: want %s, got %s", item.GetId(), created.Item)
	}

	// Get by ID
	got, err := repo.GetPersonItemById(created.Id)
	if err != nil {
		t.Fatalf("GetPersonItemById: %v", err)
	}
	if got.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), got.Person)
	}

	// Get by person ID
	pItems, err := repo.GetPersonItemsByPersonId(person.GetId())
	if err != nil {
		t.Fatalf("GetPersonItemsByPersonId: %v", err)
	}
	if len(pItems) != 1 {
		t.Fatalf("want 1 item, got %d", len(pItems))
	}

	// Delete
	if err := repo.DeletePersonItem(created.Id); err != nil {
		t.Fatalf("DeletePersonItem: %v", err)
	}
	if _, err := repo.GetPersonItemById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── EntryDeviceRepository ─────────────────────────────────────────────────────

func TestEntryDeviceRepository_CRUD(t *testing.T) {
	app := newApp(t)

	item := createRecord(t, app, "items", map[string]any{"name": "Key", "description": "test"})

	repo := repositories.NewEntryDeviceRepository(app.PocketBase)

	// Create
	created, err := repo.CreateEntryDevice(item.GetId(), "Key", "KEY-001", "None")
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

	// Get by ID
	got, err := repo.GetEntryDeviceById(created.Id)
	if err != nil {
		t.Fatalf("GetEntryDeviceById: %v", err)
	}
	if got.DeviceType != "Key" {
		t.Errorf("DeviceType: want Key, got %s", got.DeviceType)
	}

	// Update
	if err := repo.UpdateEntryDevice(created.Id, "Fob", "FOB-001", "Lost"); err != nil {
		t.Fatalf("UpdateEntryDevice: %v", err)
	}
	updated, _ := repo.GetEntryDeviceById(created.Id)
	if updated.DeviceType != "Fob" {
		t.Errorf("after update DeviceType: want Fob, got %s", updated.DeviceType)
	}
	if updated.DefunctReason != "Lost" {
		t.Errorf("after update DefunctReason: want Lost, got %s", updated.DefunctReason)
	}

	// Delete
	if err := repo.DeleteEntryDevice(created.Id); err != nil {
		t.Fatalf("DeleteEntryDevice: %v", err)
	}
	if _, err := repo.GetEntryDeviceById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
