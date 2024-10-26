package services

import (
	"encoding/json"
	"fmt"
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"time"

	"github.com/pocketbase/pocketbase/models"
)

type IDeviceHistoryServices interface {
	AddNewDeviceHistoryDueToCreateDeviceHook(deviceModel models.Model) error
	AddNewDeviceHistoryDueToUpdateDeviceHook(deviceAfterUpdateModel models.Model) error
}

type DeviceHistoryServices struct {
	deviceRepository        repositories.IDeviceRepository
	deviceHistoryRepository repositories.IDeviceHistoryRepository
}

func (dh DeviceHistoryServices) AddNewDeviceHistoryDueToCreateDeviceHook(
	deviceModel models.Model,
) error {
	return dh.deviceHistoryRepository.AddNewDeviceHistoryFromModel(
		deviceModel,
		"Device created",
		time.Now(),
	)
}

func (dh DeviceHistoryServices) AddNewDeviceHistoryDueToUpdateDeviceHook(
	deviceAfterUpdateModel models.Model,
) error {
	deviceBeforeUpdate, getDeviceErr := dh.deviceRepository.GetDeviceId(deviceAfterUpdateModel.GetId())
	if getDeviceErr != nil {
		return nil
	}
	deviceAfterUpdateModelJson, _ := json.Marshal(deviceAfterUpdateModel)
	var deviceAfterUpdate dtos.DeviceDto
	json.Unmarshal(deviceAfterUpdateModelJson, &deviceAfterUpdate)

	description := "Details updated."

	if deviceBeforeUpdate.Identifier != deviceAfterUpdate.Identifier {
		description = fmt.Sprintf(
			"%s Identifier changed from \"%s\" to \"%s\".",
			description,
			deviceBeforeUpdate.Identifier,
			deviceAfterUpdate.Identifier,
		)
	}

	if deviceBeforeUpdate.Name != deviceAfterUpdate.Name {
		description = fmt.Sprintf(
			"%s Name changed from \"%s\" to \"%s\".",
			description,
			deviceBeforeUpdate.Name,
			deviceAfterUpdate.Name,
		)
	}

	if deviceBeforeUpdate.DefunctReason != deviceAfterUpdate.DefunctReason {
		description = fmt.Sprintf(
			"%s Defunct reason changed from \"%s\" to \"%s\".",
			description,
			deviceBeforeUpdate.DefunctReason,
			deviceAfterUpdate.DefunctReason,
		)
	}

	if deviceBeforeUpdate.Type != deviceAfterUpdate.Type {
		description = fmt.Sprintf(
			"%s Type changed from \"%s\" to \"%s\".",
			description,
			deviceBeforeUpdate.Type,
			deviceAfterUpdate.Type,
		)
	}

	return dh.deviceHistoryRepository.AddNewDeviceHistoryFromModel(
		deviceAfterUpdateModel,
		description,
		deviceAfterUpdateModel.GetUpdated().Time(),
	)
}

func NewDeviceHistoryServices(
	deviceRepository repositories.IDeviceRepository,
	deviceHistoryRepository repositories.IDeviceHistoryRepository,
) IDeviceHistoryServices {
	return &DeviceHistoryServices{
		deviceRepository,
		deviceHistoryRepository,
	}
}
