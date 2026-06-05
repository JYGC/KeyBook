package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

// ── CobrandRepository ─────────────────────────────────────────────────────────

func TestCobrandRepository_CRUD(t *testing.T) {
	app := newApp(t)
	repo := repositories.NewCobrandRepository(app.PocketBase)

	// Create
	created, err := repo.CreateCobrand("Acme Co")
	if err != nil {
		t.Fatalf("CreateCobrand: %v", err)
	}
	if created.Id == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Name != "Acme Co" {
		t.Errorf("Name: want Acme Co, got %s", created.Name)
	}

	// Get by ID
	got, err := repo.GetCobrandById(created.Id)
	if err != nil {
		t.Fatalf("GetCobrandById: %v", err)
	}
	if got.Name != "Acme Co" {
		t.Errorf("Name: want Acme Co, got %s", got.Name)
	}

	// Update
	if err := repo.UpdateCobrand(created.Id, "Acme LLC"); err != nil {
		t.Fatalf("UpdateCobrand: %v", err)
	}
	updated, _ := repo.GetCobrandById(created.Id)
	if updated.Name != "Acme LLC" {
		t.Errorf("after update Name: want Acme LLC, got %s", updated.Name)
	}

	// Delete
	if err := repo.DeleteCobrand(created.Id); err != nil {
		t.Fatalf("DeleteCobrand: %v", err)
	}
	if _, err := repo.GetCobrandById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── CobrandAdminRepository ────────────────────────────────────────────────────

func TestCobrandAdminRepository_CRUD(t *testing.T) {
	app := newApp(t)

	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Test Co"})
	// Use a synthetic user ID — PocketBase DAO does not enforce FK constraints.
	const userID = "testuser00000001"

	repo := repositories.NewCobrandAdminRepository(app.PocketBase)

	// Create
	created, err := repo.CreateCobrandAdmin(userID, cobrand.GetId())
	if err != nil {
		t.Fatalf("CreateCobrandAdmin: %v", err)
	}
	if created.User != userID {
		t.Errorf("User: want %s, got %s", userID, created.User)
	}
	if created.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), created.Cobrand)
	}

	// Get by ID
	got, err := repo.GetCobrandAdminById(created.Id)
	if err != nil {
		t.Fatalf("GetCobrandAdminById: %v", err)
	}
	if got.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), got.Cobrand)
	}

	// Get by cobrand ID
	admins, err := repo.GetCobrandAdminsByCobrandId(cobrand.GetId())
	if err != nil {
		t.Fatalf("GetCobrandAdminsByCobrandId: %v", err)
	}
	if len(admins) != 1 {
		t.Fatalf("want 1 admin, got %d", len(admins))
	}

	// Delete
	if err := repo.DeleteCobrandAdmin(created.Id); err != nil {
		t.Fatalf("DeleteCobrandAdmin: %v", err)
	}
	if _, err := repo.GetCobrandAdminById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
