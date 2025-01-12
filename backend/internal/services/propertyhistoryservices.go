package services

import (
	"encoding/json"
	"fmt"
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"time"

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
	return phs.propertyHistoryRepository.AddPropertyHistoryFromModel(
		propertyModel,
		"New property record in system",
		time.Now(),
	)
}

func (phs PropertyHistoryServices) AddPropertyHistoryDueToUpdatePropertyHook(
	propertyAfterUpdateModel models.Model,
) error {
	propertyBeforeUpdate, getPropertyErr := phs.propertyRepository.GetFullPropertyById(
		propertyAfterUpdateModel.GetId(),
	)
	if getPropertyErr != nil {
		return getPropertyErr
	}
	propertyAfterUpdateModelJson, marshalErr := json.Marshal(propertyAfterUpdateModel)
	if marshalErr != nil {
		return marshalErr
	}
	var propertyAfterUpdate dtos.PropertyIdAddressDto
	if unmarshalErr := json.Unmarshal(propertyAfterUpdateModelJson, &propertyAfterUpdate); unmarshalErr != nil {
		return unmarshalErr
	}

	description := "Details updated."

	if propertyBeforeUpdate.Address != propertyAfterUpdate.Address {
		description = fmt.Sprintf(
			"%s Address changed from \"%s\" to \"%s\".",
			description,
			propertyBeforeUpdate.Address,
			propertyAfterUpdate.Address,
		)
	}

	return phs.propertyHistoryRepository.AddPropertyHistoryFromModel(
		propertyAfterUpdateModel,
		description,
		propertyAfterUpdateModel.GetUpdated().Time(),
	)
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
