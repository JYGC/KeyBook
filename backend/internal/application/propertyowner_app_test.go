package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestPropertyOwnerApplicationService_PersonOwner(t *testing.T) {
	app := newApp(t)
	svc := application.NewPropertyOwnerApplicationService(
		services.NewPropertyOwnerService(),
		repositories.NewPropertyOwnerRepository(app),
		repositories.NewPersonPropertyOwnerRepository(app),
		repositories.NewCobrandPropertyOwnerRepository(app),
	)

	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})
	person := createRecord(t, app, "persons", map[string]any{"name": "Alice", "DOB": "1990-01-01 00:00:00.000Z"})

	// Create property owner
	po, err := svc.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner: %v", err)
	}

	// Add person owner — valid
	ppo, err := svc.AddPersonOwner(po.Id, person.GetId())
	if err != nil {
		t.Fatalf("AddPersonOwner: %v", err)
	}
	if ppo.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), ppo.Person)
	}

	// Add person owner — duplicate, should fail
	if _, err := svc.AddPersonOwner(po.Id, person.GetId()); err == nil {
		t.Error("expected error for duplicate person owner, got nil")
	}

	// Remove person owner
	if err := svc.RemovePersonOwner(ppo.Id); err != nil {
		t.Fatalf("RemovePersonOwner: %v", err)
	}
}

func TestPropertyOwnerApplicationService_CobrandOwner(t *testing.T) {
	app := newApp(t)
	svc := application.NewPropertyOwnerApplicationService(
		services.NewPropertyOwnerService(),
		repositories.NewPropertyOwnerRepository(app),
		repositories.NewPersonPropertyOwnerRepository(app),
		repositories.NewCobrandPropertyOwnerRepository(app),
	)

	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})
	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Acme Co"})

	po, err := svc.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner: %v", err)
	}

	// Add cobrand owner — valid
	cpo, err := svc.AddCobrandOwner(po.Id, cobrand.GetId())
	if err != nil {
		t.Fatalf("AddCobrandOwner: %v", err)
	}
	if cpo.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), cpo.Cobrand)
	}

	// Add cobrand owner — duplicate, should fail
	if _, err := svc.AddCobrandOwner(po.Id, cobrand.GetId()); err == nil {
		t.Error("expected error for duplicate cobrand owner, got nil")
	}

	// Remove cobrand owner
	if err := svc.RemoveCobrandOwner(cpo.Id); err != nil {
		t.Fatalf("RemoveCobrandOwner: %v", err)
	}
}

func TestPropertyOwnerApplicationService_LastOwnerGuard(t *testing.T) {
	app := newApp(t)
	svc := application.NewPropertyOwnerApplicationService(
		services.NewPropertyOwnerService(),
		repositories.NewPropertyOwnerRepository(app),
		repositories.NewPersonPropertyOwnerRepository(app),
		repositories.NewCobrandPropertyOwnerRepository(app),
	)

	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	// Create a single property owner
	po1, err := svc.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner: %v", err)
	}

	// Attempt to delete the only owner — should fail
	if err := svc.DeletePropertyOwner(po1.Id); err == nil {
		t.Error("expected error when deleting last property owner, got nil")
	}

	// Create a second owner
	po2, err := svc.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner (second): %v", err)
	}

	// Now deleting one should succeed
	if err := svc.DeletePropertyOwner(po1.Id); err != nil {
		t.Fatalf("DeletePropertyOwner with multiple owners: %v", err)
	}
	_ = po2
}
