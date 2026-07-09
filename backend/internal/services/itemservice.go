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
// The DB stores "None" for active; empty string is also treated as active for domain flexibility.
func (s *ItemService) ValidateEntryDeviceTransition(currentDefunctReason, newDefunctReason string) error {
	if !isActiveReason(currentDefunctReason) && isActiveReason(newDefunctReason) {
		return errors.New("cannot reactivate a defunct entry device")
	}
	return nil
}

func isActiveReason(reason string) bool {
	return reason == "" || reason == "None"
}
