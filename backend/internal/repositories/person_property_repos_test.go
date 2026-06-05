package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

// ── PersonRepository ──────────────────────────────────────────────────────────

func TestPersonRepository_CRUD(t *testing.T) {
	app := newApp(t)
	repo := repositories.NewPersonRepository(app)

	// Create
	created, err := repo.CreatePerson("Alice", "1990-01-01 00:00:00.000Z", "")
	if err != nil {
		t.Fatalf("CreatePerson: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Name != "Alice" {
		t.Errorf("Name: want Alice, got %s", created.Name)
	}
	if created.DOB != "1990-01-01 00:00:00.000Z" {
		t.Errorf("DOB: want 1990-01-01 00:00:00.000Z, got %s", created.DOB)
	}

	// Get by ID
	got, err := repo.GetPersonById(created.Id)
	if err != nil {
		t.Fatalf("GetPersonById: %v", err)
	}
	if got.Name != "Alice" {
		t.Errorf("Name: want Alice, got %s", got.Name)
	}

	// Update
	if err := repo.UpdatePerson(created.Id, "Alice Updated", "1990-06-01 00:00:00.000Z"); err != nil {
		t.Fatalf("UpdatePerson: %v", err)
	}
	updated, err := repo.GetPersonById(created.Id)
	if err != nil {
		t.Fatalf("GetPersonById after update: %v", err)
	}
	if updated.Name != "Alice Updated" {
		t.Errorf("after update Name: want Alice Updated, got %s", updated.Name)
	}

	// Delete
	if err := repo.DeletePerson(created.Id); err != nil {
		t.Fatalf("DeletePerson: %v", err)
	}
	if _, err := repo.GetPersonById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── PropertyRepository ────────────────────────────────────────────────────────

func TestPropertyRepository_CRUD(t *testing.T) {
	app := newApp(t)
	repo := repositories.NewPropertyRepository(app)

	// Create
	created, err := repo.CreateProperty("1 Example Street")
	if err != nil {
		t.Fatalf("CreateProperty: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Address != "1 Example Street" {
		t.Errorf("Address: want 1 Example Street, got %s", created.Address)
	}

	// Get by ID
	got, err := repo.GetPropertyById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyById: %v", err)
	}
	if got.Address != "1 Example Street" {
		t.Errorf("Address: want 1 Example Street, got %s", got.Address)
	}

	// Update
	if err := repo.UpdateProperty(created.Id, "2 Updated Street"); err != nil {
		t.Fatalf("UpdateProperty: %v", err)
	}
	updated, err := repo.GetPropertyById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyById after update: %v", err)
	}
	if updated.Address != "2 Updated Street" {
		t.Errorf("after update Address: want 2 Updated Street, got %s", updated.Address)
	}

	// Delete
	if err := repo.DeleteProperty(created.Id); err != nil {
		t.Fatalf("DeleteProperty: %v", err)
	}
	if _, err := repo.GetPropertyById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
