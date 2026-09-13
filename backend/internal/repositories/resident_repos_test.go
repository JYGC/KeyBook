package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

func TestHouseholdRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Carol", "DOB": "1988-03-20 00:00:00.000Z"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	householdRepository := repositories.NewHouseholdRepository(app)

	created, err := householdRepository.CreateHousehold(person.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreateHousehold: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	got, err := householdRepository.GetHouseholdById(created.Id)
	if err != nil {
		t.Fatalf("GetHouseholdById: %v", err)
	}
	if got.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), got.Person)
	}

	members, err := householdRepository.GetHouseholdsByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetHouseholdsByPropertyId: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("want 1 member, got %d", len(members))
	}

	if err := householdRepository.DeleteHousehold(created.Id); err != nil {
		t.Fatalf("DeleteHousehold: %v", err)
	}
	if _, err := householdRepository.GetHouseholdById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestTenantRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Bob", "DOB": "1992-05-15 00:00:00.000Z"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	tenantRepository := repositories.NewTenantRepository(app)

	created, err := tenantRepository.CreateTenant(person.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreateTenant: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	got, err := tenantRepository.GetTenantById(created.Id)
	if err != nil {
		t.Fatalf("GetTenantById: %v", err)
	}
	if got.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), got.Property)
	}

	tenants, err := tenantRepository.GetTenantsByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetTenantsByPropertyId: %v", err)
	}
	if len(tenants) != 1 {
		t.Fatalf("want 1 tenant, got %d", len(tenants))
	}

	if err := tenantRepository.DeleteTenant(created.Id); err != nil {
		t.Fatalf("DeleteTenant: %v", err)
	}
	if _, err := tenantRepository.GetTenantById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
