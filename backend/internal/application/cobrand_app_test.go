package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestCobrandApplicationService_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)
	cobrandApplicationService := application.NewCobrandApplicationService(
		services.NewCobrandService(),
		repositories.NewCobrandRepository(app),
		repositories.NewCobrandAdminRepository(app),
		repositories.NewCobrandPropertyManagerRepository(app),
	)

	created, err := cobrandApplicationService.CreateCobrand("Acme Co")
	if err != nil {
		t.Fatalf("CreateCobrand: %v", err)
	}
	if created.Name != "Acme Co" {
		t.Errorf("Name: want Acme Co, got %s", created.Name)
	}

	if _, err := cobrandApplicationService.CreateCobrand(""); err == nil {
		t.Error("expected error for empty name, got nil")
	}

	if err := cobrandApplicationService.UpdateCobrand(created.Id, "Acme LLC"); err != nil {
		t.Fatalf("UpdateCobrand: %v", err)
	}

	if err := cobrandApplicationService.DeleteCobrand(created.Id); err != nil {
		t.Fatalf("DeleteCobrand: %v", err)
	}
}

func TestCobrandApplicationService_AdminUniqueness(t *testing.T) {
	app := newTestAppWithMigrations(t)
	cobrandApplicationService := application.NewCobrandApplicationService(
		services.NewCobrandService(),
		repositories.NewCobrandRepository(app),
		repositories.NewCobrandAdminRepository(app),
		repositories.NewCobrandPropertyManagerRepository(app),
	)

	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Test Co"})
	const userID = "testuser00000001"

	// Add admin — first time succeeds (bootstrap: inviter == invitee)
	admin, err := cobrandApplicationService.AddCobrandAdmin(userID, cobrand.GetId(), userID)
	if err != nil {
		t.Fatalf("AddCobrandAdmin: %v", err)
	}
	if admin.User != userID {
		t.Errorf("User: want %s, got %s", userID, admin.User)
	}

	if _, err := cobrandApplicationService.AddCobrandAdmin(userID, cobrand.GetId(), userID); err == nil {
		t.Error("expected error for duplicate admin, got nil")
	}

	if err := cobrandApplicationService.RemoveCobrandAdmin(admin.Id); err != nil {
		t.Fatalf("RemoveCobrandAdmin: %v", err)
	}
}

func TestCobrandApplicationService_InviterMustBeApprovedAdmin(t *testing.T) {
	app := newTestAppWithMigrations(t)
	cobrandAdminRepository := repositories.NewCobrandAdminRepository(app)
	cobrandApplicationService := application.NewCobrandApplicationService(
		services.NewCobrandService(),
		repositories.NewCobrandRepository(app),
		cobrandAdminRepository,
		repositories.NewCobrandPropertyManagerRepository(app),
	)

	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Approval Co"})
	const bootstrapUserID = "testuser00000002"

	// Bootstrap: no existing admins, any inviter (here, self) succeeds
	// regardless of approval — approval cannot exist yet for a brand-new
	// cobrand.
	bootstrapAdmin, err := cobrandApplicationService.AddCobrandAdmin(bootstrapUserID, cobrand.GetId(), bootstrapUserID)
	if err != nil {
		t.Fatalf("bootstrap AddCobrandAdmin: %v", err)
	}

	// The bootstrap admin defaults to unapproved (approval is a KeyBook
	// staff-only action), so they cannot yet invite a second admin.
	const secondUserID = "testuser00000003"
	if _, err := cobrandApplicationService.AddCobrandAdmin(secondUserID, cobrand.GetId(), bootstrapUserID); err == nil {
		t.Error("expected error inviting via an unapproved admin, got nil")
	}

	// Approve the bootstrap admin directly via the DAO (simulating a KeyBook
	// staff action, since there is no app-facing update path).
	record, err := app.Dao().FindRecordById("cobrandAdmins", bootstrapAdmin.Id)
	if err != nil {
		t.Fatalf("find bootstrap admin record: %v", err)
	}
	record.Set("approved", true)
	if err := app.Dao().SaveRecord(record); err != nil {
		t.Fatalf("approve bootstrap admin: %v", err)
	}

	if _, err := cobrandApplicationService.AddCobrandAdmin(secondUserID, cobrand.GetId(), bootstrapUserID); err != nil {
		t.Errorf("expected approved admin to invite successfully, got error: %v", err)
	}

	// A non-admin (unrelated user) cannot invite anyone.
	const thirdUserID = "testuser00000004"
	const unrelatedUserID = "testuser00000005"
	if _, err := cobrandApplicationService.AddCobrandAdmin(thirdUserID, cobrand.GetId(), unrelatedUserID); err == nil {
		t.Error("expected error inviting via a non-admin, got nil")
	}
}

func TestCobrandApplicationService_PropertyManager(t *testing.T) {
	app := newTestAppWithMigrations(t)
	cobrandApplicationService := application.NewCobrandApplicationService(
		services.NewCobrandService(),
		repositories.NewCobrandRepository(app),
		repositories.NewCobrandAdminRepository(app),
		repositories.NewCobrandPropertyManagerRepository(app),
	)

	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Mgr Co"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	mgr, err := cobrandApplicationService.AddPropertyManager(cobrand.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("AddPropertyManager: %v", err)
	}
	if mgr.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), mgr.Cobrand)
	}

	if err := cobrandApplicationService.RemovePropertyManager(mgr.Id); err != nil {
		t.Fatalf("RemovePropertyManager: %v", err)
	}
}
