package services

import "errors"

type IItemService interface {
	ValidateItem(name string) error
	ValidateEntryDeviceTransition(currentDefunctReason, newDefunctReason string) error
}

type ItemService struct{}

func NewItemService() IItemService {
	return &ItemService{}
}

func (s *ItemService) ValidateItem(name string) error {
	if name == "" {
		return errors.New("item name is required")
	}
	return nil
}

// ValidateEntryDeviceTransition enforces that a defunct entry device cannot be reactivated.
func (s *ItemService) ValidateEntryDeviceTransition(currentDefunctReason, newDefunctReason string) error {
	if currentDefunctReason != "" && newDefunctReason == "" {
		return errors.New("cannot reactivate a defunct entry device")
	}
	return nil
}
