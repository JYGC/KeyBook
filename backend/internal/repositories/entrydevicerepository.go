package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IEntryDeviceRepository interface {
	GetEntryDeviceById(id string) (dtos.EntryDeviceDto, error)
	CreateEntryDevice(itemId, deviceType, identifier, defunctReason string) (dtos.EntryDeviceDto, error)
	UpdateEntryDevice(id, deviceType, identifier, defunctReason string) error
	DeleteEntryDevice(id string) error
}

type EntryDeviceRepository struct {
	app core.App
}

func NewEntryDeviceRepository(app core.App) IEntryDeviceRepository {
	return &EntryDeviceRepository{app}
}

func (r *EntryDeviceRepository) GetEntryDeviceById(id string) (dtos.EntryDeviceDto, error) {
	record, err := r.app.Dao().FindRecordById("entryDevices", id)
	if err != nil {
		return dtos.EntryDeviceDto{}, err
	}
	return entryDeviceRecordToDto(record), nil
}

func (r *EntryDeviceRepository) CreateEntryDevice(itemId, deviceType, identifier, defunctReason string) (dtos.EntryDeviceDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("entryDevices")
	if err != nil {
		return dtos.EntryDeviceDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("item", itemId)
	record.Set("deviceType", deviceType)
	record.Set("identifier", identifier)
	record.Set("defunctReason", defunctReason)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.EntryDeviceDto{}, err
	}
	return entryDeviceRecordToDto(record), nil
}

func (r *EntryDeviceRepository) UpdateEntryDevice(id, deviceType, identifier, defunctReason string) error {
	record, err := r.app.Dao().FindRecordById("entryDevices", id)
	if err != nil {
		return err
	}
	record.Set("deviceType", deviceType)
	record.Set("identifier", identifier)
	record.Set("defunctReason", defunctReason)
	return r.app.Dao().SaveRecord(record)
}

func (r *EntryDeviceRepository) DeleteEntryDevice(id string) error {
	record, err := r.app.Dao().FindRecordById("entryDevices", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func entryDeviceRecordToDto(r *models.Record) dtos.EntryDeviceDto {
	return dtos.EntryDeviceDto{
		Id:            r.GetId(),
		Item:          r.GetString("item"),
		DeviceType:    r.GetString("deviceType"),
		Identifier:    r.GetString("identifier"),
		DefunctReason: r.GetString("defunctReason"),
	}
}
