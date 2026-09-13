package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

func TestAgentRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Dave", "DOB": "1985-07-10 00:00:00.000Z"})
	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Agency Co"})

	agentRepository := repositories.NewAgentRepository(app)

	created, err := agentRepository.CreateAgent(person.GetId(), cobrand.GetId())
	if err != nil {
		t.Fatalf("CreateAgent: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), created.Cobrand)
	}

	got, err := agentRepository.GetAgentById(created.Id)
	if err != nil {
		t.Fatalf("GetAgentById: %v", err)
	}
	if got.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), got.Person)
	}

	agents, err := agentRepository.GetAgentsByCobrandId(cobrand.GetId())
	if err != nil {
		t.Fatalf("GetAgentsByCobrandId: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("want 1 agent, got %d", len(agents))
	}

	if err := agentRepository.DeleteAgent(created.Id); err != nil {
		t.Fatalf("DeleteAgent: %v", err)
	}
	if _, err := agentRepository.GetAgentById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestPropertyAgentRepository_CRUD(t *testing.T) {
	app := newTestAppWithMigrations(t)

	person := createRecordBypassingAccessRules(t, app, "persons", map[string]any{"name": "Dave", "DOB": "1985-07-10 00:00:00.000Z"})
	cobrand := createRecordBypassingAccessRules(t, app, "cobrands", map[string]any{"name": "Agency Co"})
	agent := createRecordBypassingAccessRules(t, app, "agents", map[string]any{"person": person.GetId(), "cobrand": cobrand.GetId()})
	property := createRecordBypassingAccessRules(t, app, "properties", map[string]any{"address": "1 Test St"})

	propertyAgentRepository := repositories.NewPropertyAgentRepository(app)

	created, err := propertyAgentRepository.CreatePropertyAgent(agent.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyAgent: %v", err)
	}
	if created.Agent != agent.GetId() {
		t.Errorf("Agent: want %s, got %s", agent.GetId(), created.Agent)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	got, err := propertyAgentRepository.GetPropertyAgentById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyAgentById: %v", err)
	}
	if got.Agent != agent.GetId() {
		t.Errorf("Agent: want %s, got %s", agent.GetId(), got.Agent)
	}

	agents, err := propertyAgentRepository.GetPropertyAgentsByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetPropertyAgentsByPropertyId: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("want 1 agent, got %d", len(agents))
	}

	if err := propertyAgentRepository.DeletePropertyAgent(created.Id); err != nil {
		t.Fatalf("DeletePropertyAgent: %v", err)
	}
	if _, err := propertyAgentRepository.GetPropertyAgentById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
