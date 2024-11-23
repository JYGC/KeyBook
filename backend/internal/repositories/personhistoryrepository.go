package repositories

import (
	"encoding/json"
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
	newPersonHistory := models.NewRecord(personHistoryCollection)
	newPersonHistory.Set("description", description)
	newPersonHistory.Set("person", string(personJson))
	newPersonHistory.Set("stateddatetime", statedDateTime)
	newPersonHistory.Set("deviceid", personModel.GetId())
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
