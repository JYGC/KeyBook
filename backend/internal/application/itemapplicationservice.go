package application

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

type IItemApplicationService interface {
	CreateItem(name, description string) (dtos.ItemDto, error)
	UpdateItem(id, name, description string) error
	DeleteItem(id string) error
	CreateEntryDevice(itemId, deviceType, identifier, defunctReason string) (dtos.EntryDeviceDto, error)
	UpdateEntryDevice(id, deviceType, identifier, defunctReason string) error
	DeleteEntryDevice(id string) error
	AddPropertyItem(itemId, propertyId string) (dtos.PropertyItemDto, error)
	RemovePropertyItem(id string) error
	AddPersonItem(personId, itemId string) (dtos.PersonItemDto, error)
	RemovePersonItem(id string) error
}

type ItemApplicationService struct {
	itemService            services.IItemService
	itemRepository         repositories.IItemRepository
	entryDeviceRepository  repositories.IEntryDeviceRepository
	propertyItemRepository repositories.IPropertyItemRepository
	personItemRepository   repositories.IPersonItemRepository
}

func NewItemApplicationService(
	itemService services.IItemService,
	itemRepository repositories.IItemRepository,
	entryDeviceRepository repositories.IEntryDeviceRepository,
	propertyItemRepository repositories.IPropertyItemRepository,
	personItemRepository repositories.IPersonItemRepository,
) IItemApplicationService {
	return &ItemApplicationService{itemService, itemRepository, entryDeviceRepository, propertyItemRepository, personItemRepository}
}

func (s *ItemApplicationService) CreateItem(name, description string) (dtos.ItemDto, error) {
	if err := s.itemService.ValidateItem(name); err != nil {
		return dtos.ItemDto{}, err
	}
	return s.itemRepository.CreateItem(name, description)
}

func (s *ItemApplicationService) UpdateItem(id, name, description string) error {
	if err := s.itemService.ValidateItem(name); err != nil {
		return err
	}
	return s.itemRepository.UpdateItem(id, name, description)
}

func (s *ItemApplicationService) DeleteItem(id string) error {
	return s.itemRepository.DeleteItem(id)
}

func (s *ItemApplicationService) CreateEntryDevice(itemId, deviceType, identifier, defunctReason string) (dtos.EntryDeviceDto, error) {
	return s.entryDeviceRepository.CreateEntryDevice(itemId, deviceType, identifier, defunctReason)
}

func (s *ItemApplicationService) UpdateEntryDevice(id, deviceType, identifier, defunctReason string) error {
	currentEntryDevice, err := s.entryDeviceRepository.GetEntryDeviceById(id)
	if err != nil {
		return err
	}
	if err := s.itemService.ValidateEntryDeviceTransition(currentEntryDevice.DefunctReason, defunctReason); err != nil {
		return err
	}
	return s.entryDeviceRepository.UpdateEntryDevice(id, deviceType, identifier, defunctReason)
}

func (s *ItemApplicationService) DeleteEntryDevice(id string) error {
	return s.entryDeviceRepository.DeleteEntryDevice(id)
}

func (s *ItemApplicationService) AddPropertyItem(itemId, propertyId string) (dtos.PropertyItemDto, error) {
	return s.propertyItemRepository.CreatePropertyItem(itemId, propertyId)
}

func (s *ItemApplicationService) RemovePropertyItem(id string) error {
	return s.propertyItemRepository.DeletePropertyItem(id)
}

func (s *ItemApplicationService) AddPersonItem(personId, itemId string) (dtos.PersonItemDto, error) {
	return s.personItemRepository.CreatePersonItem(personId, itemId)
}

func (s *ItemApplicationService) RemovePersonItem(id string) error {
	return s.personItemRepository.DeletePersonItem(id)
}
