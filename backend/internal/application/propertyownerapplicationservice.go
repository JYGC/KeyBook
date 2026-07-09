package application

import (
	"errors"

	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

type IPropertyOwnerApplicationService interface {
	CreatePropertyOwner(propertyId string) (dtos.PropertyOwnerDto, error)
	DeletePropertyOwner(id string) error
	AddPersonOwner(propertyOwnerId, personId string) (dtos.PersonPropertyOwnerDto, error)
	RemovePersonOwner(id string) error
	AddCobrandOwner(propertyOwnerId, cobrandId string) (dtos.CobrandPropertyOwnerDto, error)
	RemoveCobrandOwner(id string) error
}

type PropertyOwnerApplicationService struct {
	propertyOwnerService    services.IPropertyOwnerService
	propertyOwnerRepo       repositories.IPropertyOwnerRepository
	personPropertyOwnerRepo repositories.IPersonPropertyOwnerRepository
	cobrandPropertyOwnerRepo repositories.ICobrandPropertyOwnerRepository
}

func NewPropertyOwnerApplicationService(
	propertyOwnerService services.IPropertyOwnerService,
	propertyOwnerRepo repositories.IPropertyOwnerRepository,
	personPropertyOwnerRepo repositories.IPersonPropertyOwnerRepository,
	cobrandPropertyOwnerRepo repositories.ICobrandPropertyOwnerRepository,
) IPropertyOwnerApplicationService {
	return &PropertyOwnerApplicationService{
		propertyOwnerService,
		propertyOwnerRepo,
		personPropertyOwnerRepo,
		cobrandPropertyOwnerRepo,
	}
}

func (s *PropertyOwnerApplicationService) CreatePropertyOwner(propertyId string) (dtos.PropertyOwnerDto, error) {
	return s.propertyOwnerRepo.CreatePropertyOwner(propertyId)
}

func (s *PropertyOwnerApplicationService) DeletePropertyOwner(id string) error {
	owner, err := s.propertyOwnerRepo.GetPropertyOwnerById(id)
	if err != nil {
		return err
	}
	owners, err := s.propertyOwnerRepo.GetPropertyOwnersByPropertyId(owner.Property)
	if err != nil {
		return err
	}
	if len(owners) <= 1 {
		return errors.New("cannot delete the last property owner")
	}
	return s.propertyOwnerRepo.DeletePropertyOwner(id)
}

func (s *PropertyOwnerApplicationService) AddPersonOwner(propertyOwnerId, personId string) (dtos.PersonPropertyOwnerDto, error) {
	existing, err := s.personPropertyOwnerRepo.GetPersonPropertyOwnersByPropertyOwnerId(propertyOwnerId)
	if err != nil {
		return dtos.PersonPropertyOwnerDto{}, err
	}
	if err := s.propertyOwnerService.EnsureNoDuplicatePersonOwner(existing, personId); err != nil {
		return dtos.PersonPropertyOwnerDto{}, err
	}
	return s.personPropertyOwnerRepo.CreatePersonPropertyOwner(personId, propertyOwnerId)
}

func (s *PropertyOwnerApplicationService) RemovePersonOwner(id string) error {
	return s.personPropertyOwnerRepo.DeletePersonPropertyOwner(id)
}

func (s *PropertyOwnerApplicationService) AddCobrandOwner(propertyOwnerId, cobrandId string) (dtos.CobrandPropertyOwnerDto, error) {
	existing, err := s.cobrandPropertyOwnerRepo.GetCobrandPropertyOwnersByPropertyOwnerId(propertyOwnerId)
	if err != nil {
		return dtos.CobrandPropertyOwnerDto{}, err
	}
	if err := s.propertyOwnerService.EnsureNoDuplicateCobrandOwner(existing, cobrandId); err != nil {
		return dtos.CobrandPropertyOwnerDto{}, err
	}
	return s.cobrandPropertyOwnerRepo.CreateCobrandPropertyOwner(cobrandId, propertyOwnerId)
}

func (s *PropertyOwnerApplicationService) RemoveCobrandOwner(id string) error {
	return s.cobrandPropertyOwnerRepo.DeleteCobrandPropertyOwner(id)
}
