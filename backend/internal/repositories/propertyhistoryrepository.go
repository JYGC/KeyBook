package repositories

import (
	"encoding/json"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IPropertyHistoryRepository interface {
	AddPropertyHistoryFromModel(
		propertyModel models.Model,
		description string,
		statedDateTime time.Time,
	) error
}

type PropertyHistoryRepository struct {
	app *pocketbase.PocketBase
}

func (phr PropertyHistoryRepository) AddPropertyHistoryFromModel(
	propertyModel models.Model,
	description string,
	statedDateTime time.Time,
) error {
	propertyHistoryCollection, findHistoryCollectionErr := phr.app.Dao().FindCollectionByNameOrId("propertyhistory")
	if findHistoryCollectionErr != nil {
		return findHistoryCollectionErr
	}
	propertyJson, mashalErr := json.Marshal(propertyModel)
	if mashalErr != nil {
		return mashalErr
	}
	newPropertyHistory := models.NewRecord(propertyHistoryCollection)
	newPropertyHistory.Set("description", description)
	newPropertyHistory.Set("property", string(propertyJson))
	newPropertyHistory.Set("stateddatetime", statedDateTime)
	newPropertyHistory.Set("propertyid", propertyModel.GetId())
	if saveHistoryErr := phr.app.Dao().SaveRecord(newPropertyHistory); saveHistoryErr != nil {
		return saveHistoryErr
	}
	return nil
}

func NewPropertyHistoryRepository(app *pocketbase.PocketBase) IPropertyHistoryRepository {
	return &PropertyHistoryRepository{
		app,
	}
}
