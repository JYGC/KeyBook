package services

import (
	"keybook/backend/internal/repositories"

	"github.com/pocketbase/pocketbase/models"
)

type IPropertyHistoryServices interface {
	AddPropertyHistoryDueToCreatePropertyHook(
		propertyModel models.Model,
	) error
	AddPropertyHistoryDueToUpdatePropertyHook(
		propertyAfterUpdateModel models.Model,
	) error
}

type PropertyHistoryServices struct {
	propertyRepository        repositories.IPropertyRepository
	propertyHistoryRepository repositories.IPropertyHistoryRepository
}

func (phs PropertyHistoryServices) AddPropertyHistoryDueToCreatePropertyHook(
	propertyModel models.Model,
) error {
	return nil
}

func (phs PropertyHistoryServices) AddPropertyHistoryDueToUpdatePropertyHook(
	propertyAfterUpdateModel models.Model,
) error {
	return nil
}

func NewPropertyHistoryServices(
	propertyRepository repositories.IPropertyRepository,
	propertyHistoryRepository repositories.IPropertyHistoryRepository,
) IPropertyHistoryServices {
	return PropertyHistoryServices{
		propertyRepository,
		propertyHistoryRepository,
	}
}
