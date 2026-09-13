package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestPropertyApplicationService_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)
	propertyApplicationService := application.NewPropertyApplicationService(
		services.NewPropertyService(),
		repositories.NewPropertyRepository(app),
	)

	created, err := propertyApplicationService.CreateProperty("1 Main St")
	if err != nil {
		t.Fatalf("CreateProperty: %v", err)
	}
	if created.Address != "1 Main St" {
		t.Errorf("Address: want 1 Main St, got %s", created.Address)
	}

	if _, err := propertyApplicationService.CreateProperty(""); err == nil {
		t.Error("expected error for empty address, got nil")
	}

	if err := propertyApplicationService.UpdateProperty(created.Id, "2 Updated Rd"); err != nil {
		t.Fatalf("UpdateProperty: %v", err)
	}

	if err := propertyApplicationService.UpdateProperty(created.Id, ""); err == nil {
		t.Error("expected error for empty address on update, got nil")
	}

	if err := propertyApplicationService.DeleteProperty(created.Id); err != nil {
		t.Fatalf("DeleteProperty: %v", err)
	}
}
