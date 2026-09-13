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
	propertyOwnerService           services.IPropertyOwnerService
	propertyOwnerRepository        repositories.IPropertyOwnerRepository
	personPropertyOwnerRepository  repositories.IPersonPropertyOwnerRepository
	cobrandPropertyOwnerRepository repositories.ICobrandPropertyOwnerRepository
}

func NewPropertyOwnerApplicationService(
	propertyOwnerService services.IPropertyOwnerService,
	propertyOwnerRepository repositories.IPropertyOwnerRepository,
	personPropertyOwnerRepository repositories.IPersonPropertyOwnerRepository,
	cobrandPropertyOwnerRepository repositories.ICobrandPropertyOwnerRepository,
) IPropertyOwnerApplicationService {
	return &PropertyOwnerApplicationService{
		propertyOwnerService,
		propertyOwnerRepository,
		personPropertyOwnerRepository,
		cobrandPropertyOwnerRepository,
	}
}

func (s *PropertyOwnerApplicationService) CreatePropertyOwner(propertyId string) (dtos.PropertyOwnerDto, error) {
	return s.propertyOwnerRepository.CreatePropertyOwner(propertyId)
}

func (s *PropertyOwnerApplicationService) DeletePropertyOwner(id string) error {
	ownerToDelete, err := s.propertyOwnerRepository.GetPropertyOwnerById(id)
	if err != nil {
		return err
	}
	ownersOfSameProperty, err := s.propertyOwnerRepository.GetPropertyOwnersByPropertyId(ownerToDelete.Property)
	if err != nil {
		return err
	}
	if len(ownersOfSameProperty) <= 1 {
		return errors.New("cannot delete the last property owner")
	}
	return s.propertyOwnerRepository.DeletePropertyOwner(id)
}

func (s *PropertyOwnerApplicationService) AddPersonOwner(propertyOwnerId, personId string) (dtos.PersonPropertyOwnerDto, error) {
	existingPersonOwners, err := s.personPropertyOwnerRepository.GetPersonPropertyOwnersByPropertyOwnerId(propertyOwnerId)
	if err != nil {
		return dtos.PersonPropertyOwnerDto{}, err
	}
	if err := s.propertyOwnerService.EnsureNoDuplicatePersonOwner(existingPersonOwners, personId); err != nil {
		return dtos.PersonPropertyOwnerDto{}, err
	}
	return s.personPropertyOwnerRepository.CreatePersonPropertyOwner(personId, propertyOwnerId)
}

func (s *PropertyOwnerApplicationService) RemovePersonOwner(id string) error {
	return s.personPropertyOwnerRepository.DeletePersonPropertyOwner(id)
}

func (s *PropertyOwnerApplicationService) AddCobrandOwner(propertyOwnerId, cobrandId string) (dtos.CobrandPropertyOwnerDto, error) {
	existingCobrandOwners, err := s.cobrandPropertyOwnerRepository.GetCobrandPropertyOwnersByPropertyOwnerId(propertyOwnerId)
	if err != nil {
		return dtos.CobrandPropertyOwnerDto{}, err
	}
	if err := s.propertyOwnerService.EnsureNoDuplicateCobrandOwner(existingCobrandOwners, cobrandId); err != nil {
		return dtos.CobrandPropertyOwnerDto{}, err
	}
	return s.cobrandPropertyOwnerRepository.CreateCobrandPropertyOwner(cobrandId, propertyOwnerId)
}

func (s *PropertyOwnerApplicationService) RemoveCobrandOwner(id string) error {
	return s.cobrandPropertyOwnerRepository.DeleteCobrandPropertyOwner(id)
}
