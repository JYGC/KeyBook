package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

func TestPropertyOwnerRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	propertyOwnerRepository := repositories.NewPropertyOwnerRepository(app)

	created, err := propertyOwnerRepository.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	got, err := propertyOwnerRepository.GetPropertyOwnerById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyOwnerById: %v", err)
	}
	if got.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), got.Property)
	}

	owners, err := propertyOwnerRepository.GetPropertyOwnersByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetPropertyOwnersByPropertyId: %v", err)
	}
	if len(owners) != 1 {
		t.Fatalf("want 1 owner, got %d", len(owners))
	}

	if err := propertyOwnerRepository.DeletePropertyOwner(created.Id); err != nil {
		t.Fatalf("DeletePropertyOwner: %v", err)
	}
	if _, err := propertyOwnerRepository.GetPropertyOwnerById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestPersonPropertyOwnerRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Alice", "DOB": "1990-01-01 00:00:00.000Z"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})
	propertyOwner := createRecordBypassingAccessRules(t, app, "propertyOwners", map[string]any{"property": property.GetId()})

	personPropertyOwnerRepository := repositories.NewPersonPropertyOwnerRepository(app)

	created, err := personPropertyOwnerRepository.CreatePersonPropertyOwner(person.GetId(), propertyOwner.GetId())
	if err != nil {
		t.Fatalf("CreatePersonPropertyOwner: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.PropertyOwner != propertyOwner.GetId() {
		t.Errorf("PropertyOwner: want %s, got %s", propertyOwner.GetId(), created.PropertyOwner)
	}

	got, err := personPropertyOwnerRepository.GetPersonPropertyOwnerById(created.Id)
	if err != nil {
		t.Fatalf("GetPersonPropertyOwnerById: %v", err)
	}
	if got.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), got.Person)
	}

	owners, err := personPropertyOwnerRepository.GetPersonPropertyOwnersByPropertyOwnerId(propertyOwner.GetId())
	if err != nil {
		t.Fatalf("GetPersonPropertyOwnersByPropertyOwnerId: %v", err)
	}
	if len(owners) != 1 {
		t.Fatalf("want 1 owner, got %d", len(owners))
	}

	if err := personPropertyOwnerRepository.DeletePersonPropertyOwner(created.Id); err != nil {
		t.Fatalf("DeletePersonPropertyOwner: %v", err)
	}
	if _, err := personPropertyOwnerRepository.GetPersonPropertyOwnerById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestCobrandPropertyOwnerRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Test Co"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})
	propertyOwner := createRecordBypassingAccessRules(t, app, "propertyOwners", map[string]any{"property": property.GetId()})

	cobrandPropertyOwnerRepository := repositories.NewCobrandPropertyOwnerRepository(app)

	created, err := cobrandPropertyOwnerRepository.CreateCobrandPropertyOwner(cobrand.GetId(), propertyOwner.GetId())
	if err != nil {
		t.Fatalf("CreateCobrandPropertyOwner: %v", err)
	}
	if created.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), created.Cobrand)
	}
	if created.PropertyOwner != propertyOwner.GetId() {
		t.Errorf("PropertyOwner: want %s, got %s", propertyOwner.GetId(), created.PropertyOwner)
	}

	got, err := cobrandPropertyOwnerRepository.GetCobrandPropertyOwnerById(created.Id)
	if err != nil {
		t.Fatalf("GetCobrandPropertyOwnerById: %v", err)
	}
	if got.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), got.Cobrand)
	}

	owners, err := cobrandPropertyOwnerRepository.GetCobrandPropertyOwnersByPropertyOwnerId(propertyOwner.GetId())
	if err != nil {
		t.Fatalf("GetCobrandPropertyOwnersByPropertyOwnerId: %v", err)
	}
	if len(owners) != 1 {
		t.Fatalf("want 1 owner, got %d", len(owners))
	}

	if err := cobrandPropertyOwnerRepository.DeleteCobrandPropertyOwner(created.Id); err != nil {
		t.Fatalf("DeleteCobrandPropertyOwner: %v", err)
	}
	if _, err := cobrandPropertyOwnerRepository.GetCobrandPropertyOwnerById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestCobrandPropertyManagerRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Manager Co"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	cobrandPropertyManagerRepository := repositories.NewCobrandPropertyManagerRepository(app)

	created, err := cobrandPropertyManagerRepository.CreateCobrandPropertyManager(cobrand.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreateCobrandPropertyManager: %v", err)
	}
	if created.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), created.Cobrand)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	got, err := cobrandPropertyManagerRepository.GetCobrandPropertyManagerById(created.Id)
	if err != nil {
		t.Fatalf("GetCobrandPropertyManagerById: %v", err)
	}
	if got.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), got.Property)
	}

	managers, err := cobrandPropertyManagerRepository.GetCobrandPropertyManagersByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetCobrandPropertyManagersByPropertyId: %v", err)
	}
	if len(managers) != 1 {
		t.Fatalf("want 1 manager, got %d", len(managers))
	}

	if err := cobrandPropertyManagerRepository.DeleteCobrandPropertyManager(created.Id); err != nil {
		t.Fatalf("DeleteCobrandPropertyManager: %v", err)
	}
	if _, err := cobrandPropertyManagerRepository.GetCobrandPropertyManagerById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
