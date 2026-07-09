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
	agentService      services.IAgentService
	agentRepo         repositories.IAgentRepository
	propertyAgentRepo repositories.IPropertyAgentRepository
}

func NewAgentApplicationService(
	agentService services.IAgentService,
	agentRepo repositories.IAgentRepository,
	propertyAgentRepo repositories.IPropertyAgentRepository,
) IAgentApplicationService {
	return &AgentApplicationService{agentService, agentRepo, propertyAgentRepo}
}

func (s *AgentApplicationService) RegisterAgent(personId, cobrandId string) (dtos.AgentDto, error) {
	return s.agentRepo.CreateAgent(personId, cobrandId)
}

func (s *AgentApplicationService) UnregisterAgent(id string) error {
	return s.agentRepo.DeleteAgent(id)
}

func (s *AgentApplicationService) AssignToProperty(agentId, propertyId string) (dtos.PropertyAgentDto, error) {
	existing, err := s.propertyAgentRepo.GetPropertyAgentsByPropertyId(propertyId)
	if err != nil {
		return dtos.PropertyAgentDto{}, err
	}
	if err := s.agentService.EnsureNoDuplicatePropertyAgent(existing, agentId); err != nil {
		return dtos.PropertyAgentDto{}, err
	}
	return s.propertyAgentRepo.CreatePropertyAgent(agentId, propertyId)
}

func (s *AgentApplicationService) ResignFromProperty(id string) error {
	return s.propertyAgentRepo.DeletePropertyAgent(id)
}
