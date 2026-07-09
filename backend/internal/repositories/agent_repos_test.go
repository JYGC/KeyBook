package repositories_test

import (
	"testing"

	"keybook/backend/internal/repositories"
)

// ── AgentRepository ───────────────────────────────────────────────────────────

func TestAgentRepository_CRUD(t *testing.T) {
	app := newApp(t)

	person := createRecord(t, app, "persons", map[string]any{"name": "Dave", "DOB": "1985-07-10 00:00:00.000Z"})
	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Agency Co"})

	repo := repositories.NewAgentRepository(app)

	// Create
	created, err := repo.CreateAgent(person.GetId(), cobrand.GetId())
	if err != nil {
		t.Fatalf("CreateAgent: %v", err)
	}
	if created.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), created.Person)
	}
	if created.Cobrand != cobrand.GetId() {
		t.Errorf("Cobrand: want %s, got %s", cobrand.GetId(), created.Cobrand)
	}

	// Get by ID
	got, err := repo.GetAgentById(created.Id)
	if err != nil {
		t.Fatalf("GetAgentById: %v", err)
	}
	if got.Person != person.GetId() {
		t.Errorf("Person: want %s, got %s", person.GetId(), got.Person)
	}

	// Get by cobrand ID
	agents, err := repo.GetAgentsByCobrandId(cobrand.GetId())
	if err != nil {
		t.Fatalf("GetAgentsByCobrandId: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("want 1 agent, got %d", len(agents))
	}

	// Delete
	if err := repo.DeleteAgent(created.Id); err != nil {
		t.Fatalf("DeleteAgent: %v", err)
	}
	if _, err := repo.GetAgentById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

// ── PropertyAgentRepository ───────────────────────────────────────────────────

func TestPropertyAgentRepository_CRUD(t *testing.T) {
	app := newApp(t)

	person := createRecord(t, app, "persons", map[string]any{"name": "Dave", "DOB": "1985-07-10 00:00:00.000Z"})
	cobrand := createRecord(t, app, "cobrands", map[string]any{"name": "Agency Co"})
	agent := createRecord(t, app, "agents", map[string]any{"person": person.GetId(), "cobrand": cobrand.GetId()})
	property := createRecord(t, app, "properties", map[string]any{"address": "1 Test St"})

	repo := repositories.NewPropertyAgentRepository(app)

	// Create
	created, err := repo.CreatePropertyAgent(agent.GetId(), property.GetId())
	if err != nil {
		t.Fatalf("CreatePropertyAgent: %v", err)
	}
	if created.Agent != agent.GetId() {
		t.Errorf("Agent: want %s, got %s", agent.GetId(), created.Agent)
	}
	if created.Property != property.GetId() {
		t.Errorf("Property: want %s, got %s", property.GetId(), created.Property)
	}

	// Get by ID
	got, err := repo.GetPropertyAgentById(created.Id)
	if err != nil {
		t.Fatalf("GetPropertyAgentById: %v", err)
	}
	if got.Agent != agent.GetId() {
		t.Errorf("Agent: want %s, got %s", agent.GetId(), got.Agent)
	}

	// Get by property ID
	agents, err := repo.GetPropertyAgentsByPropertyId(property.GetId())
	if err != nil {
		t.Fatalf("GetPropertyAgentsByPropertyId: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("want 1 agent, got %d", len(agents))
	}

	// Delete
	if err := repo.DeletePropertyAgent(created.Id); err != nil {
		t.Fatalf("DeletePropertyAgent: %v", err)
	}
	if _, err := repo.GetPropertyAgentById(created.Id); err == nil {
		t.Error("expected error after delete, got nil")
	}
}
