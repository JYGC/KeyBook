package repositories

import (
	"encoding/json"
	"keybook/backend/internal/dtos"
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

	var propertyIdDto dtos.PersonDeviceDto
	if unmashalErr := json.Unmarshal(personDeviceJson, &propertyIdDto); unmashalErr != nil {
		return unmashalErr
	}

	newPersonDeviceHistory := models.NewRecord(personDeviceHistoryCollection)
	newPersonDeviceHistory.Set("description", description)
	newPersonDeviceHistory.Set("persondevice", string(personDeviceJson))
	newPersonDeviceHistory.Set("stateddatetime", statedDateTime)
	newPersonDeviceHistory.Set("persondeviceid", personDeviceModel.GetId())
	newPersonDeviceHistory.Set("property", propertyIdDto.Property)
	newPersonDeviceHistory.Set("deviceid", propertyIdDto.Device)
	newPersonDeviceHistory.Set("personid", propertyIdDto.Person)

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
