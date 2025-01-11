package services

import (
	"encoding/json"
	"fmt"
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
	personRepository              repositories.IPersonRepository
	personDeviceHistoryRepository repositories.IPersonDeviceHistoryRepository
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

	newDeviceHolder, getNewDeviceHolderErr := pdhs.personRepository.GetPersonById(personDevice.Person)
	if getNewDeviceHolderErr != nil {
		return getNewDeviceHolderErr
	}

	return pdhs.personDeviceHistoryRepository.AddNewPersonDeviceHistoryFromModel(
		personDeviceModel,
		fmt.Sprintf("Device given to %s", newDeviceHolder.Name),
		time.Now(),
	)
}

func (pdhs PersonDeviceHistoryServices) AddNewPersonDeviceHistoryDueToUpdatePersonDeviceHook(
	personDeviceAfterUpdateModel models.Model,
) error {
	personDeviceAfterUpdateModelJson, mashalErr := json.Marshal(personDeviceAfterUpdateModel)
	if mashalErr != nil {
		return mashalErr
	}

	var personDeviceAfterUpdate dtos.PersonDeviceDto
	if unmashalErr := json.Unmarshal(personDeviceAfterUpdateModelJson, &personDeviceAfterUpdate); unmashalErr != nil {
		return unmashalErr
	}

	newDeviceHolder, getNewDeviceHolderErr := pdhs.personRepository.GetPersonById(personDeviceAfterUpdate.Person)
	if getNewDeviceHolderErr != nil {
		return getNewDeviceHolderErr
	}

	return pdhs.personDeviceHistoryRepository.AddNewPersonDeviceHistoryFromModel(
		personDeviceAfterUpdateModel,
		fmt.Sprintf("Device given to %s", newDeviceHolder.Name),
		personDeviceAfterUpdateModel.GetUpdated().Time(),
	)
}

func NewPersonDeviceHistoryServices(
	personRepository repositories.IPersonRepository,
	personDeviceHistoryRepository repositories.IPersonDeviceHistoryRepository,

) IPersonDeviceHistoryServices {
	return PersonDeviceHistoryServices{
		personRepository,
		personDeviceHistoryRepository,
	}
}
