package application

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

type IAgentApplicationService interface {
	RegisterAgent(personId, cobrandId string) (dtos.AgentDto, error)
	UnregisterAgent(id string) error
	AssignToProperty(agentId, propertyId string) (dtos.PropertyAgentDto, error)
	ResignFromProperty(id string) error
}

type AgentApplicationService struct {
	agentService            services.IAgentService
	agentRepository         repositories.IAgentRepository
	propertyAgentRepository repositories.IPropertyAgentRepository
}

func NewAgentApplicationService(
	agentService services.IAgentService,
	agentRepository repositories.IAgentRepository,
	propertyAgentRepository repositories.IPropertyAgentRepository,
) IAgentApplicationService {
	return &AgentApplicationService{agentService, agentRepository, propertyAgentRepository}
}

func (s *AgentApplicationService) RegisterAgent(personId, cobrandId string) (dtos.AgentDto, error) {
	return s.agentRepository.CreateAgent(personId, cobrandId)
}

func (s *AgentApplicationService) UnregisterAgent(id string) error {
	return s.agentRepository.DeleteAgent(id)
}

func (s *AgentApplicationService) AssignToProperty(agentId, propertyId string) (dtos.PropertyAgentDto, error) {
	existingPropertyAgents, err := s.propertyAgentRepository.GetPropertyAgentsByPropertyId(propertyId)
	if err != nil {
		return dtos.PropertyAgentDto{}, err
	}
	if err := s.agentService.EnsureNoDuplicatePropertyAgent(existingPropertyAgents, agentId); err != nil {
		return dtos.PropertyAgentDto{}, err
	}
	return s.propertyAgentRepository.CreatePropertyAgent(agentId, propertyId)
}

func (s *AgentApplicationService) ResignFromProperty(id string) error {
	return s.propertyAgentRepository.DeletePropertyAgent(id)
}
