package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestPersonApplicationService_CRUD(t *testing.T) {
	app := newApp(t)
	svc := application.NewPersonApplicationService(
		services.NewPersonService(),
		repositories.NewPersonRepository(app),
	)

	// Create — valid
	created, err := svc.CreatePerson("Alice", "1990-01-01 00:00:00.000Z", "")
	if err != nil {
		t.Fatalf("CreatePerson: %v", err)
	}
	if created.Name != "Alice" {
		t.Errorf("Name: want Alice, got %s", created.Name)
	}

	// Create — validation failure
	if _, err := svc.CreatePerson("", "1990-01-01 00:00:00.000Z", ""); err == nil {
		t.Error("expected error for empty name, got nil")
	}
	if _, err := svc.CreatePerson("Bob", "", ""); err == nil {
		t.Error("expected error for empty DOB, got nil")
	}

	// Update — valid
	if err := svc.UpdatePerson(created.Id, "Alice Updated", "1991-06-15 00:00:00.000Z"); err != nil {
		t.Fatalf("UpdatePerson: %v", err)
	}

	// Update — validation failure
	if err := svc.UpdatePerson(created.Id, "", "1991-06-15 00:00:00.000Z"); err == nil {
		t.Error("expected error for empty name on update, got nil")
	}

	// Delete
	if err := svc.DeletePerson(created.Id); err != nil {
		t.Fatalf("DeletePerson: %v", err)
	}
}
