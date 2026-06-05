package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestCobrandApplicationService_CRUD(t *testing.T) {
	app := newApp(t)
	svc := application.NewCobrandApplicationService(
		services.NewCobrandService(),
		repositories.NewCobrandRepository(app),
		repositories.NewCobrandAdminRepository(app),
		repositories.NewCobrandPropertyManagerRepository(app),
	)

	// Create — valid
	created, err := svc.CreateCobrand("Acme Co")
	if err != nil {
		t.Fatalf("CreateCobrand: %v", err)
	}
	if created.Name != "Acme Co" {
		t.Errorf("Name: want Acme Co, got %s", created.Name)
	}

	// Create — validation failure
	if _, err := svc.CreateCobrand(""); err == nil {
		t.Error("expected error for empty name, got nil")
	}

	// Update — valid
	if err := svc.UpdateCobrand(created.Id, "Acme LLC"); err != nil {
		t.Fatalf("UpdateCobrand: %v", err)
	}

	// Delete
	if err := svc.DeleteCobrand(created.Id); err != nil {
		t.Fatalf("DeleteCobrand: %v", err)
	}
}

func TestCobrandApplicationService_AdminUniqueness(t *testing.T) {
	app := newApp(t)
	svc := application.NewCobrandApplicationService(
		services.NewCobrandService(),
		repositories.NewCobrandRepository(app),
		repositories.NewCobrandAdminRepository(app),
		repositories.NewCobrandPropertyManagerRepository(app),
	)

	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Test Co"})
	const userID = "testuser00000001"

	// Add admin — first time succeeds
	admin, err := svc.AddCobrandAdmin(userID, cobrand.GetId())
	if err != nil {
		t.Fatalf("AddCobrandAdmin: %v", err)
	}
	if admin.User != userID {
		t.Errorf("User: want %s, got %s", userID, admin.User)
	}

	// Add admin — duplicate fails
	if _, err := svc.AddCobrandAdmin(userID, cobrand.GetId()); err == nil {
		t.Error("expected error for duplicate admin, got nil")
	}

	// Remove admin
	if err := svc.RemoveCobrandAdmin(admin.Id); err != nil {
		t.Fatalf("RemoveCobrandAdmin: %v", err)
	}
}

func TestCobrandApplicationService_PropertyManager(t *testing.T) {
	app := newApp(t)
	svc := application.NewCobrandApplicationService(
		services.NewCobrandService(),
		repositories.NewCobrandRepository(app),
		repositories.NewCobrandAdminRepository(app),
		repositories.NewCobrandPropertyManagerRepository(app),
	)

	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Mgr Co"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	mgr, err := svc.AddPropertyManager(cobrand.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("AddPropertyManager: %v", err)
	}
	if mgr.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), mgr.Cobrand)
	}

	if err := svc.RemovePropertyManager(mgr.Id); err != nil {
		t.Fatalf("RemovePropertyManager: %v", err)
	}
}
