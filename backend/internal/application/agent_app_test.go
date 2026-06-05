package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestAgentApplicationService_RegisterAndAssign(t *testing.T) {
	app := newApp(t)
	svc := application.NewAgentApplicationService(
		services.NewAgentService(),
		repositories.NewAgentRepository(app),
		repositories.NewPropertyAgentRepository(app),
	)

	person := createRecord(t, app, "persons", map[string]any{"name": "Dave", "DOB": "1985-07-10 00:00:00.000Z"})
	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Agency Co"})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	// Register agent
	agent, err := svc.RegisterAgent(person.GetId(), cobrand.GetId())
	if err != nil {
		t.Fatalf("RegisterAgent: %v", err)
	}
	if agent.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), agent.Person)
	}

	// Assign to property — first time succeeds
	pa, err := svc.AssignToProperty(agent.Id, property.GetId())
	if err != nil {
		t.Fatalf("AssignToProperty: %v", err)
	}
	if pa.Agent != agent.Id {
		t.Errorf("Agent: want %s, got %s", agent.Id, pa.Agent)
	}

	// Assign to same property again — duplicate, should fail
	if _, err := svc.AssignToProperty(agent.Id, property.GetId()); err == nil {
		t.Error("expected error for duplicate property assignment, got nil")
	}

	// Resign from property
	if err := svc.ResignFromProperty(pa.Id); err != nil {
		t.Fatalf("ResignFromProperty: %v", err)
	}

	// Unregister agent
	if err := svc.UnregisterAgent(agent.Id); err != nil {
		t.Fatalf("UnregisterAgent: %v", err)
	}
}
