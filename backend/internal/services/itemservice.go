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

func (s *ItemService) ValidateEntryDeviceTransition(currentDefunctReason, newDefunctReason string) error {
	if !defunctReasonMeansActive(currentDefunctReason) && defunctReasonMeansActive(newDefunctReason) {
		return errors.New("cannot reactivate a defunct entry device")
	}
	return nil
}

// The DB stores the literal "None" for an active device; the empty string is
// also accepted so callers that omit the field are not treated as defunct.
func defunctReasonMeansActive(defunctReason string) bool {
	return defunctReason == "" || defunctReason == "None"
}
