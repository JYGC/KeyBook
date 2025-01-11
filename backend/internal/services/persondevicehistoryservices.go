package services

import (
	"encoding/json"
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"time"

	"github.com/pocketbase/pocketbase/models"
)

type IPersonDeviceHistoryServices interface {
	AddNewPersonDeviceHistoryDueToCreatePersonDeviceHook(
		devicePersonModel models.Model,
	) error
	AddNewPersonDeviceHistoryDueToUpdatePersonDeviceHook(
		devicePersonAfterUpdateModel models.Model,
	) error
}

type PersonDeviceHistoryServices struct {
	//personDeviceRepository        repositories.IPersonDeviceRepository
	personDeviceHistoryRepository repositories.IPersonDeviceHistoryRepository
	deviceHistory                 repositories.IDeviceRepository
}

func (pdhs PersonDeviceHistoryServices) AddNewPersonDeviceHistoryDueToCreatePersonDeviceHook(
	personDeviceModel models.Model,
) error {
	personDeviceModelJson, mashalErr := json.Marshal(personDeviceModel)
	if mashalErr != nil {
		return mashalErr
	}

	var personDevice dtos.PersonDeviceDto
	if unmashalErr := json.Unmarshal(personDeviceModelJson, &personDevice); unmashalErr != nil {
		return unmashalErr
	}

	return pdhs.personDeviceHistoryRepository.AddNewPersonDeviceHistoryFromModel(
		personDeviceModel,
		"Person device assignment created",
		time.Now(),
	)
}

func (pdhs PersonDeviceHistoryServices) AddNewPersonDeviceHistoryDueToUpdatePersonDeviceHook(
	personDeviceAfterUpdateModel models.Model,
) error {
	// personDeviceBeforeUpdate, getPersonDeviceErr := pdhs.personDeviceRepository.GetPersonDeviceById(
	// 	personDeviceAfterUpdateModel.GetId(),
	// )
	// if getPersonDeviceErr != nil {
	// 	return getPersonDeviceErr
	// }

	personDeviceAfterUpdateModelJson, mashalErr := json.Marshal(personDeviceAfterUpdateModel)
	if mashalErr != nil {
		return mashalErr
	}

	var personDeviceAfterUpdate dtos.PersonDeviceDto
	if unmashalErr := json.Unmarshal(personDeviceAfterUpdateModelJson, &personDeviceAfterUpdate); unmashalErr != nil {
		return unmashalErr
	}

	return pdhs.personDeviceHistoryRepository.AddNewPersonDeviceHistoryFromModel(
		personDeviceAfterUpdateModel,
		"Assignment changed",
		personDeviceAfterUpdateModel.GetUpdated().Time(),
	)
}

func NewPersonDeviceHistoryServices(
	personDeviceHistoryRepository repositories.IPersonDeviceHistoryRepository,
	deviceHistory repositories.IDeviceRepository,
) IPersonDeviceHistoryServices {
	return PersonDeviceHistoryServices{
		personDeviceHistoryRepository,
		deviceHistory,
	}
}
