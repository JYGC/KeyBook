package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

func TestPersonRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)
	personRepository := repositories.NewPersonRepository(app)

	created, err := personRepository.CreatePerson("Alice", "1990-01-01 00:00:00.000Z", "")
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

	got, err := personRepository.GetPersonById(created.Id)
	if err != nil {
		t.Fatalf("GetPersonById: %v", err)
	}
	if got.Name != "Alice" {
		t.Errorf("Name: want Alice, got %s", got.Name)
	}

	if err := personRepository.UpdatePerson(created.Id, "Alice Updated", "1990-06-01 00:00:00.000Z"); err != nil {
		t.Fatalf("UpdatePerson: %v", err)
	}
	updated, err := personRepository.GetPersonById(created.Id)
	if err != nil {
		t.Fatalf("GetPersonById after update: %v", err)
	}
	if updated.Name != "Alice Updated" {
		t.Errorf("after update Name: want Alice Updated, got %s", updated.Name)
	}

	if err := personRepository.DeletePerson(created.Id); err != nil {
		t.Fatalf("DeletePerson: %v", err)
	}
	if _, err := personRepository.GetPersonById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestPropertyRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)
	propertyRepository := repositories.NewPropertyRepository(app)

	created, err := propertyRepository.CreateProperty("1 Example Street")
	if err != nil {
		t.Fatalf("CreateProperty: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Address != "1 Example Street" {
		t.Errorf("Address: want 1 Example Street, got %s", created.Address)
	}

	got, err := propertyRepository.GetPropertyById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyById: %v", err)
	}
	if got.Address != "1 Example Street" {
		t.Errorf("Address: want 1 Example Street, got %s", got.Address)
	}

	if err := propertyRepository.UpdateProperty(created.Id, "2 Updated Street"); err != nil {
		t.Fatalf("UpdateProperty: %v", err)
	}
	updated, err := propertyRepository.GetPropertyById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyById after update: %v", err)
	}
	if updated.Address != "2 Updated Street" {
		t.Errorf("after update Address: want 2 Updated Street, got %s", updated.Address)
	}

	if err := propertyRepository.DeleteProperty(created.Id); err != nil {
		t.Fatalf("DeleteProperty: %v", err)
	}
	if _, err := propertyRepository.GetPropertyById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
