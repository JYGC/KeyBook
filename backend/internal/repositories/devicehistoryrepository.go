package repositories

import (
	"encoding/json"
	"keybook/backend/internal/dtos"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IDeviceHistoryRepository interface {
	AddNewDeviceHistoryFromModel(
		deviceModel models.Model,
		description string,
		statedDateTime time.Time,
	) error
}

type DeviceHistoryRepository struct {
	app *pocketbase.PocketBase
}

func (dh DeviceHistoryRepository) AddNewDeviceHistoryFromModel(
	deviceModel models.Model,
	description string,
	statedDateTime time.Time,
) error {
	deviceHistoryCollection, findHistoryCollectionErr := dh.app.Dao().FindCollectionByNameOrId("deviceHistory")
	if findHistoryCollectionErr != nil {
		return findHistoryCollectionErr
	}

	deviceJson, mashalErr := json.Marshal(deviceModel)
	if mashalErr != nil {
		return mashalErr
	}

	var propertyIdDto dtos.PropertyIdFromJsonDto
	if unmashalErr := json.Unmarshal(deviceJson, &propertyIdDto); unmashalErr != nil {
		return unmashalErr
	}

	newDeviceHistory := models.NewRecord(deviceHistoryCollection)
	newDeviceHistory.Set("description", description)
	newDeviceHistory.Set("device", string(deviceJson))
	newDeviceHistory.Set("statedDateTime", statedDateTime)
	newDeviceHistory.Set("deviceId", deviceModel.GetId())
	newDeviceHistory.Set("property", propertyIdDto.Property)

	if saveHistoryErr := dh.app.Dao().SaveRecord(newDeviceHistory); saveHistoryErr != nil {
		return saveHistoryErr
	}
	return nil
}

func NewDeviceHistoryRepository(app *pocketbase.PocketBase) IDeviceHistoryRepository {
	return DeviceHistoryRepository{
		app,
	}
}
