package repositories

import (
	"encoding/json"
	"keybook/backend/internal/dtos"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IPersonHistoryRepository interface {
	AddNewPersonHistoryFromModel(
		personModel models.Model,
		description string,
		statedDateTime time.Time,
	) error
}

type PersonHistoryRepository struct {
	app *pocketbase.PocketBase
}

func (ph PersonHistoryRepository) AddNewPersonHistoryFromModel(
	personModel models.Model,
	description string,
	statedDateTime time.Time,
) error {
	personHistoryCollection, findHistoryCollectionErr := ph.app.Dao().FindCollectionByNameOrId("personhistory")
	if findHistoryCollectionErr != nil {
		return findHistoryCollectionErr
	}

	personJson, mashalErr := json.Marshal(personModel)
	if mashalErr != nil {
		return mashalErr
	}

	var propertyIdDto dtos.PropertyIdFromJsonDto
	if unmashalErr := json.Unmarshal(personJson, &propertyIdDto); unmashalErr != nil {
		return unmashalErr
	}

	newPersonHistory := models.NewRecord(personHistoryCollection)
	newPersonHistory.Set("description", description)
	newPersonHistory.Set("person", string(personJson))
	newPersonHistory.Set("statedDateTime", statedDateTime)
	newPersonHistory.Set("personId", personModel.GetId())
	newPersonHistory.Set("property", propertyIdDto.Property)

	if saveHistoryErr := ph.app.Dao().SaveRecord(newPersonHistory); saveHistoryErr != nil {
		return saveHistoryErr
	}
	return nil
}

func NewPersonHistoryRepository(app *pocketbase.PocketBase) IPersonHistoryRepository {
	return PersonHistoryRepository{
		app,
	}
}
