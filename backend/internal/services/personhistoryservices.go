package services

import (
	"encoding/json"
	"fmt"
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"time"

	"github.com/pocketbase/pocketbase/models"
)

type IPersonHistoryServices interface {
	AddNewPersonHistoryDueToCreatePersonHook(
		personModel models.Model,
	) error
	AddNewPersonHistoryDueToUpdatePersonHook(
		personAfterUpdateModel models.Model,
	) error
}

type PersonHistoryServices struct {
	personRepository        repositories.IPersonRepository
	personHistoryRepository repositories.IPersonHistoryRepository
}

func (ph PersonHistoryServices) AddNewPersonHistoryDueToCreatePersonHook(
	personModel models.Model,
) error {
	return ph.personHistoryRepository.AddNewPersonHistoryFromModel(
		personModel,
		"Person created",
		time.Now(),
	)
}

func (ph PersonHistoryServices) AddNewPersonHistoryDueToUpdatePersonHook(
	personAfterUpdateModel models.Model,
) error {
	personBeforeUpdate, getPersonErr := ph.personRepository.GetPersonById(
		personAfterUpdateModel.GetId(),
	)
	if getPersonErr != nil {
		return getPersonErr
	}
	personAfterUpdateModelJson, _ := json.Marshal(personAfterUpdateModel)
	var personAfterUpdate dtos.PersonDto
	json.Unmarshal(personAfterUpdateModelJson, &personAfterUpdate)

	description := ""

	if personBeforeUpdate.Name != personAfterUpdate.Name {
		description = fmt.Sprintf(
			"%s Name changed from \"%s\" to \"%s\".",
			description,
			personBeforeUpdate.Name,
			personAfterUpdate.Name,
		)
	}

	if personBeforeUpdate.Type != personAfterUpdate.Type {
		description = fmt.Sprintf(
			"%s Type changed from \"%s\" to \"%s\".",
			description,
			personBeforeUpdate.Type,
			personAfterUpdate.Type,
		)
	}

	if description == "" {
		return nil
	}

	return ph.personHistoryRepository.AddNewPersonHistoryFromModel(
		personAfterUpdateModel,
		description,
		personAfterUpdateModel.GetUpdated().Time(),
	)
}

func NewPersonHistoryServices(
	personRepository repositories.IPersonRepository,
	personHistoryRepository repositories.IPersonHistoryRepository,
) IPersonHistoryServices {
	return &PersonHistoryServices{
		personRepository,
		personHistoryRepository,
	}
}
