package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

// ── HouseholdRepository ───────────────────────────────────────────────────────

func TestHouseholdRepository_CRUD(t *testing.T) {
	app := newApp(t)

	person := createRecord(t, app, "persons", map[string]any{"name": "Carol", "DOB": "1988-03-20 00:00:00.000Z"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	repo := repositories.NewHouseholdRepository(app.PocketBase)

	// Create
	created, err := repo.CreateHousehold(person.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	// Get by ID
	got, err := repo.GetHouseholdById(created.Id)
	if err != nil {
		t.Fatalf("GetHouseholdById: %v", err)
	}
	if got.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), got.Person)
	}

	// Get by property ID
	members, err := repo.GetHouseholdsByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetHouseholdsByPropertyId: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("want 1 member, got %d", len(members))
	}

	// Delete
	if err := repo.DeleteHousehold(created.Id); err != nil {
		t.Fatalf("DeleteHousehold: %v", err)
	}
	if _, err := repo.GetHouseholdById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── TenantRepository ──────────────────────────────────────────────────────────

func TestTenantRepository_CRUD(t *testing.T) {
	app := newApp(t)

	person := createRecord(t, app, "persons", map[string]any{"name": "Bob", "DOB": "1992-05-15 00:00:00.000Z"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	repo := repositories.NewTenantRepository(app.PocketBase)

	// Create
	created, err := repo.CreateTenant(person.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreateTenant: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	// Get by ID
	got, err := repo.GetTenantById(created.Id)
	if err != nil {
		t.Fatalf("GetTenantById: %v", err)
	}
	if got.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), got.Property)
	}

	// Get by property ID
	tenants, err := repo.GetTenantsByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetTenantsByPropertyId: %v", err)
	}
	if len(tenants) != 1 {
		t.Fatalf("want 1 tenant, got %d", len(tenants))
	}

	// Delete
	if err := repo.DeleteTenant(created.Id); err != nil {
		t.Fatalf("DeleteTenant: %v", err)
	}
	if _, err := repo.GetTenantById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
