package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

// ── PropertyOwnerRepository ───────────────────────────────────────────────────

func TestPropertyOwnerRepository_CRUD(t *testing.T) {
	app := newApp(t)

	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	repo := repositories.NewPropertyOwnerRepository(app)

	// Create
	created, err := repo.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	// Get by ID
	got, err := repo.GetPropertyOwnerById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyOwnerById: %v", err)
	}
	if got.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), got.Property)
	}

	// Get by property ID
	owners, err := repo.GetPropertyOwnersByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetPropertyOwnersByPropertyId: %v", err)
	}
	if len(owners) != 1 {
		t.Fatalf("want 1 owner, got %d", len(owners))
	}

	// Delete
	if err := repo.DeletePropertyOwner(created.Id); err != nil {
		t.Fatalf("DeletePropertyOwner: %v", err)
	}
	if _, err := repo.GetPropertyOwnerById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── PersonPropertyOwnerRepository ────────────────────────────────────────────

func TestPersonPropertyOwnerRepository_CRUD(t *testing.T) {
	app := newApp(t)

	person := createRecord(t, app, "persons", map[string]any{"name": "Alice", "DOB": "1990-01-01 00:00:00.000Z"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})
	propertyOwner := createRecord(t, app, "propertyOwners", map[string]any{"property": property.GetId()})

	repo := repositories.NewPersonPropertyOwnerRepository(app)

	// Create
	created, err := repo.CreatePersonPropertyOwner(person.GetId(), propertyOwner.GetId())
	if err != nil {
		t.Fatalf("CreatePersonPropertyOwner: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.PropertyOwner != propertyOwner.GetId() {
		t.Errorf("PropertyOwner: want %s, got %s", propertyOwner.GetId(), created.PropertyOwner)
	}

	// Get by ID
	got, err := repo.GetPersonPropertyOwnerById(created.Id)
	if err != nil {
		t.Fatalf("GetPersonPropertyOwnerById: %v", err)
	}
	if got.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), got.Person)
	}

	// Get by propertyOwner ID
	owners, err := repo.GetPersonPropertyOwnersByPropertyOwnerId(propertyOwner.GetId())
	if err != nil {
		t.Fatalf("GetPersonPropertyOwnersByPropertyOwnerId: %v", err)
	}
	if len(owners) != 1 {
		t.Fatalf("want 1 owner, got %d", len(owners))
	}

	// Delete
	if err := repo.DeletePersonPropertyOwner(created.Id); err != nil {
		t.Fatalf("DeletePersonPropertyOwner: %v", err)
	}
	if _, err := repo.GetPersonPropertyOwnerById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── CobrandPropertyOwnerRepository ───────────────────────────────────────────

func TestCobrandPropertyOwnerRepository_CRUD(t *testing.T) {
	app := newApp(t)

	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Test Co"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})
	propertyOwner := createRecord(t, app, "propertyOwners", map[string]any{"property": property.GetId()})

	repo := repositories.NewCobrandPropertyOwnerRepository(app)

	// Create
	created, err := repo.CreateCobrandPropertyOwner(cobrand.GetId(), propertyOwner.GetId())
	if err != nil {
		t.Fatalf("CreateCobrandPropertyOwner: %v", err)
	}
	if created.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), created.Cobrand)
	}
	if created.PropertyOwner != propertyOwner.GetId() {
		t.Errorf("PropertyOwner: want %s, got %s", propertyOwner.GetId(), created.PropertyOwner)
	}

	// Get by ID
	got, err := repo.GetCobrandPropertyOwnerById(created.Id)
	if err != nil {
		t.Fatalf("GetCobrandPropertyOwnerById: %v", err)
	}
	if got.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), got.Cobrand)
	}

	// Get by propertyOwner ID
	owners, err := repo.GetCobrandPropertyOwnersByPropertyOwnerId(propertyOwner.GetId())
	if err != nil {
		t.Fatalf("GetCobrandPropertyOwnersByPropertyOwnerId: %v", err)
	}
	if len(owners) != 1 {
		t.Fatalf("want 1 owner, got %d", len(owners))
	}

	// Delete
	if err := repo.DeleteCobrandPropertyOwner(created.Id); err != nil {
		t.Fatalf("DeleteCobrandPropertyOwner: %v", err)
	}
	if _, err := repo.GetCobrandPropertyOwnerById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── CobrandPropertyManagerRepository ─────────────────────────────────────────

func TestCobrandPropertyManagerRepository_CRUD(t *testing.T) {
	app := newApp(t)

	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Manager Co"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	repo := repositories.NewCobrandPropertyManagerRepository(app)

	// Create
	created, err := repo.CreateCobrandPropertyManager(cobrand.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreateCobrandPropertyManager: %v", err)
	}
	if created.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), created.Cobrand)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	// Get by ID
	got, err := repo.GetCobrandPropertyManagerById(created.Id)
	if err != nil {
		t.Fatalf("GetCobrandPropertyManagerById: %v", err)
	}
	if got.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), got.Property)
	}

	// Get by property ID
	managers, err := repo.GetCobrandPropertyManagersByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetCobrandPropertyManagersByPropertyId: %v", err)
	}
	if len(managers) != 1 {
		t.Fatalf("want 1 manager, got %d", len(managers))
	}

	// Delete
	if err := repo.DeleteCobrandPropertyManager(created.Id); err != nil {
		t.Fatalf("DeleteCobrandPropertyManager: %v", err)
	}
	if _, err := repo.GetCobrandPropertyManagerById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
