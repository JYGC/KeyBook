package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

func TestCobrandRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)
	cobrandRepository := repositories.NewCobrandRepository(app)

	created, err := cobrandRepository.CreateCobrand("Acme Co")
	if err != nil {
		t.Fatalf("CreateCobrand: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Name != "Acme Co" {
		t.Errorf("Name: want Acme Co, got %s", created.Name)
	}

	got, err := cobrandRepository.GetCobrandById(created.Id)
	if err != nil {
		t.Fatalf("GetCobrandById: %v", err)
	}
	if got.Name != "Acme Co" {
		t.Errorf("Name: want Acme Co, got %s", got.Name)
	}

	if err := cobrandRepository.UpdateCobrand(created.Id, "Acme LLC"); err != nil {
		t.Fatalf("UpdateCobrand: %v", err)
	}
	updated, _ := cobrandRepository.GetCobrandById(created.Id)
	if updated.Name != "Acme LLC" {
		t.Errorf("after update Name: want Acme LLC, got %s", updated.Name)
	}

	if err := cobrandRepository.DeleteCobrand(created.Id); err != nil {
		t.Fatalf("DeleteCobrand: %v", err)
	}
	if _, err := cobrandRepository.GetCobrandById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestCobrandAdminRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Test Co"})
	// Use a synthetic user ID — PocketBase DAO does not enforce FK constraints.
	const userID = "testuser00000001"

	cobrandAdminRepository := repositories.NewCobrandAdminRepository(app)

	created, err := cobrandAdminRepository.CreateCobrandAdmin(userID, cobrand.GetId())
	if err != nil {
		t.Fatalf("CreateCobrandAdmin: %v", err)
	}
	if created.User != userID {
		t.Errorf("User: want %s, got %s", userID, created.User)
	}
	if created.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), created.Cobrand)
	}
	if created.Approved {
		t.Error("Approved: want false on creation, got true")
	}

	got, err := cobrandAdminRepository.GetCobrandAdminById(created.Id)
	if err != nil {
		t.Fatalf("GetCobrandAdminById: %v", err)
	}
	if got.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), got.Cobrand)
	}
	if got.Approved {
		t.Error("Approved: want false, got true")
	}

	admins, err := cobrandAdminRepository.GetCobrandAdminsByCobrandId(cobrand.GetId())
	if err != nil {
		t.Fatalf("GetCobrandAdminsByCobrandId: %v", err)
	}
	if len(admins) != 1 {
		t.Fatalf("want 1 admin, got %d", len(admins))
	}

	if err := cobrandAdminRepository.DeleteCobrandAdmin(created.Id); err != nil {
		t.Fatalf("DeleteCobrandAdmin: %v", err)
	}
	if _, err := cobrandAdminRepository.GetCobrandAdminById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
