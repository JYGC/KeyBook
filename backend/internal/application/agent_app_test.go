package application_test

import (
	"testing"

	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

func TestAgentApplicationService_RegisterAndAssign(t *testing.T) {
	app := newTestAppWithMigrations(t)
	agentApplicationService := application.NewAgentApplicationService(
		services.NewAgentService(),
		repositories.NewAgentRepository(app),
		repositories.NewPropertyAgentRepository(app),
	)

	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Dave", "DOB": "1985-07-10 00:00:00.000Z"})
	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Agency Co"})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	agent, err := agentApplicationService.RegisterAgent(person.GetId(), cobrand.GetId())
	if err != nil {
		t.Fatalf("RegisterAgent: %v", err)
	}
	if agent.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), agent.Person)
	}

	pa, err := agentApplicationService.AssignToProperty(agent.Id, property.GetId())
	if err != nil {
		t.Fatalf("AssignToProperty: %v", err)
	}
	if pa.Agent != agent.Id {
		t.Errorf("Agent: want %s, got %s", agent.Id, pa.Agent)
	}

	if _, err := agentApplicationService.AssignToProperty(agent.Id, property.GetId()); err == nil {
		t.Error("expected error for duplicate property assignment, got nil")
	}

	if err := agentApplicationService.ResignFromProperty(pa.Id); err != nil {
		t.Fatalf("ResignFromProperty: %v", err)
	}

	if err := agentApplicationService.UnregisterAgent(agent.Id); err != nil {
		t.Fatalf("UnregisterAgent: %v", err)
	}
}
