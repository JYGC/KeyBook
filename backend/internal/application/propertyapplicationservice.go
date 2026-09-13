package application

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

type IPropertyApplicationService interface {
	CreateProperty(address string) (dtos.PropertyDto, error)
	UpdateProperty(id, address string) error
	DeleteProperty(id string) error
}

type PropertyApplicationService struct {
	propertyService    services.IPropertyService
	propertyRepository repositories.IPropertyRepository
}

func NewPropertyApplicationService(
	propertyService services.IPropertyService,
	propertyRepository repositories.IPropertyRepository,
) IPropertyApplicationService {
	return &PropertyApplicationService{propertyService, propertyRepository}
}

func (s *PropertyApplicationService) CreateProperty(address string) (dtos.PropertyDto, error) {
	if err := s.propertyService.ValidateProperty(address); err != nil {
		return dtos.PropertyDto{}, err
	}
	return s.propertyRepository.CreateProperty(address)
}

func (s *PropertyApplicationService) UpdateProperty(id, address string) error {
	if err := s.propertyService.ValidateProperty(address); err != nil {
		return err
	}
	return s.propertyRepository.UpdateProperty(id, address)
}

func (s *PropertyApplicationService) DeleteProperty(id string) error {
	return s.propertyRepository.DeleteProperty(id)
}
