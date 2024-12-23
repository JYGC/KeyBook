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
}

func (pdhs PersonDeviceHistoryServices) AddNewPersonDeviceHistoryDueToCreatePersonDeviceHook(
	personDeviceModel models.Model,
) error {
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
	personDeviceAfterUpdateModelJson, _ := json.Marshal(personDeviceAfterUpdateModel)
	var personDeviceAfterUpdate dtos.PersonDeviceDto
	json.Unmarshal(personDeviceAfterUpdateModelJson, &personDeviceAfterUpdate)

	return pdhs.personDeviceHistoryRepository.AddNewPersonDeviceHistoryFromModel(
		personDeviceAfterUpdateModel,
		"Assignment changed",
		personDeviceAfterUpdateModel.GetUpdated().Time(),
	)
}

func NewPersonDeviceHistoryServices(
	personDeviceHistoryRepository repositories.IPersonDeviceHistoryRepository,
) IPersonDeviceHistoryServices {
	return PersonDeviceHistoryServices{
		personDeviceHistoryRepository,
	}
}
