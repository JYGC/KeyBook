package repositories

import (
	"encoding/json"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IPersonDeviceHistoryRepository interface {
	AddNewPersonDeviceHistoryFromModel(
		personDeviceModel models.Model,
		description string,
		statedDateTime time.Time,
	) error
}

type PersonDeviceHistoryRepository struct {
	app *pocketbase.PocketBase
}

func (pdhr PersonDeviceHistoryRepository) AddNewPersonDeviceHistoryFromModel(
	personDeviceModel models.Model,
	description string,
	statedDateTime time.Time,
) error {
	personDeviceHistoryCollection, findHistoryCollectionErr := pdhr.app.Dao().FindCollectionByNameOrId("persondevicehistory")
	if findHistoryCollectionErr != nil {
		return findHistoryCollectionErr
	}
	personDeviceJson, mashalErr := json.Marshal(personDeviceModel)
	if mashalErr != nil {
		return mashalErr
	}
	newPersonDeviceHistory := models.NewRecord(personDeviceHistoryCollection)
	newPersonDeviceHistory.Set("description", description)
	newPersonDeviceHistory.Set("persondevice", string(personDeviceJson))
	newPersonDeviceHistory.Set("stateddatetime", statedDateTime)
	newPersonDeviceHistory.Set("persondeviceid", personDeviceModel.GetId())
	if saveHistoryErr := pdhr.app.Dao().SaveRecord(newPersonDeviceHistory); saveHistoryErr != nil {
		return saveHistoryErr
	}
	return nil
}

func NewPersonDeviceHistoryRepository(app *pocketbase.PocketBase) IPersonDeviceHistoryRepository {
	return PersonDeviceHistoryRepository{
		app,
	}
}
