package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestPropertyOwnerApplicationService_PersonOwner(t *testing.T) {
	app := newTestAppWithMigrations(t)
	propertyOwnerApplicationService := application.NewPropertyOwnerApplicationService(
		services.NewPropertyOwnerService(),
		repositories.NewPropertyOwnerRepository(app),
		repositories.NewPersonPropertyOwnerRepository(app),
		repositories.NewCobrandPropertyOwnerRepository(app),
	)

	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})
	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Alice", "DOB": "1990-01-01 00:00:00.000Z"})

	propertyOwner, err := propertyOwnerApplicationService.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner: %v", err)
	}

	personPropertyOwner, err := propertyOwnerApplicationService.AddPersonOwner(propertyOwner.Id, person.GetId())
	if err != nil {
		t.Fatalf("AddPersonOwner: %v", err)
	}
	if personPropertyOwner.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), personPropertyOwner.Person)
	}

	if _, err := propertyOwnerApplicationService.AddPersonOwner(propertyOwner.Id, person.GetId()); err == nil {
		t.Error("expected error for duplicate person owner, got nil")
	}

	if err := propertyOwnerApplicationService.RemovePersonOwner(personPropertyOwner.Id); err != nil {
		t.Fatalf("RemovePersonOwner: %v", err)
	}
}

func TestPropertyOwnerApplicationService_CobrandOwner(t *testing.T) {
	app := newTestAppWithMigrations(t)
	propertyOwnerApplicationService := application.NewPropertyOwnerApplicationService(
		services.NewPropertyOwnerService(),
		repositories.NewPropertyOwnerRepository(app),
		repositories.NewPersonPropertyOwnerRepository(app),
		repositories.NewCobrandPropertyOwnerRepository(app),
	)

	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})
	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Acme Co"})

	propertyOwner, err := propertyOwnerApplicationService.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner: %v", err)
	}

	cobrandPropertyOwner, err := propertyOwnerApplicationService.AddCobrandOwner(propertyOwner.Id, cobrand.GetId())
	if err != nil {
		t.Fatalf("AddCobrandOwner: %v", err)
	}
	if cobrandPropertyOwner.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), cobrandPropertyOwner.Cobrand)
	}

	if _, err := propertyOwnerApplicationService.AddCobrandOwner(propertyOwner.Id, cobrand.GetId()); err == nil {
		t.Error("expected error for duplicate cobrand owner, got nil")
	}

	if err := propertyOwnerApplicationService.RemoveCobrandOwner(cobrandPropertyOwner.Id); err != nil {
		t.Fatalf("RemoveCobrandOwner: %v", err)
	}
}

func TestPropertyOwnerApplicationService_LastOwnerGuard(t *testing.T) {
	app := newTestAppWithMigrations(t)
	propertyOwnerApplicationService := application.NewPropertyOwnerApplicationService(
		services.NewPropertyOwnerService(),
		repositories.NewPropertyOwnerRepository(app),
		repositories.NewPersonPropertyOwnerRepository(app),
		repositories.NewCobrandPropertyOwnerRepository(app),
	)

	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	po1, err := propertyOwnerApplicationService.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner: %v", err)
	}

	if err := propertyOwnerApplicationService.DeletePropertyOwner(po1.Id); err == nil {
		t.Error("expected error when deleting last property owner, got nil")
	}

	po2, err := propertyOwnerApplicationService.CreatePropertyOwner(property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyOwner (second): %v", err)
	}

	if err := propertyOwnerApplicationService.DeletePropertyOwner(po1.Id); err != nil {
		t.Fatalf("DeletePropertyOwner with multiple owners: %v", err)
	}
	_ = po2
}
