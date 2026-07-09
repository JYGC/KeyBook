package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestPropertyApplicationService_CRUD(t *testing.T) {
	app := newApp(t)
	svc := application.NewPropertyApplicationService(
		services.NewPropertyService(),
		repositories.NewPropertyRepository(app),
	)

	// Create — valid
	created, err := svc.CreateProperty("1 Main St")
	if err != nil {
		t.Fatalf("CreateProperty: %v", err)
	}
	if created.Address != "1 Main St" {
		t.Errorf("Address: want 1 Main St, got %s", created.Address)
	}

	// Create — validation failure
	if _, err := svc.CreateProperty(""); err == nil {
		t.Error("expected error for empty address, got nil")
	}

	// Update — valid
	if err := svc.UpdateProperty(created.Id, "2 Updated Rd"); err != nil {
		t.Fatalf("UpdateProperty: %v", err)
	}

	// Update — validation failure
	if err := svc.UpdateProperty(created.Id, ""); err == nil {
		t.Error("expected error for empty address on update, got nil")
	}

	// Delete
	if err := svc.DeleteProperty(created.Id); err != nil {
		t.Fatalf("DeleteProperty: %v", err)
	}
}
