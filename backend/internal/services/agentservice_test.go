package services_test

import (
	"testing"

	"keybook/backend/internal/dtos"
	"keybook/backend/internal/services"
)

func TestAgentService_EnsureNoDuplicatePropertyAgent(t *testing.T) {
	svc := services.NewAgentService()

	existing := []dtos.PropertyAgentDto{
		{Id: "pa1", Agent: "agent1", Property: "prop1"},
	}

	tests := []struct {
		name        string
		assignments []dtos.PropertyAgentDto
		agentId     string
		wantErr     bool
	}{
		{"first assignment", []dtos.PropertyAgentDto{}, "agent1", false},
		{"different agent", existing, "agent2", false},
		{"duplicate assignment", existing, "agent1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.EnsureNoDuplicatePropertyAgent(tt.assignments, tt.agentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureNoDuplicatePropertyAgent(_, %q) error = %v, wantErr %v", tt.agentId, err, tt.wantErr)
			}
		})
	}
}
