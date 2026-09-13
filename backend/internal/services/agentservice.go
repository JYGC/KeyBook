package services

import (
	"errors"

	"keybook/backend/internal/dtos"
)

type IAgentService interface {
	EnsureNoDuplicatePropertyAgent(existing []dtos.PropertyAgentDto, agentId string) error
}

type AgentService struct{}

func NewAgentService() IAgentService {
	return &AgentService{}
}

func (s *AgentService) EnsureNoDuplicatePropertyAgent(existing []dtos.PropertyAgentDto, agentId string) error {
	for _, existingPropertyAgent := range existing {
		if existingPropertyAgent.Agent == agentId {
			return errors.New("agent is already assigned to this property")
		}
	}
	return nil
}
