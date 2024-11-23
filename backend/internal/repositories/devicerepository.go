package repositories

import (
	"encoding/json"
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type IDeviceRepository interface {
	GetDeviceById(deviceId string) (
		dtos.DeviceDto,
		error,
	)
	GetDeviceForPropertyId(
		propertyId string,
		inboundDevices []dtos.NewDeviceDto,
	) (
		[]dtos.DeviceDto,
		error,
	)
	GetDevicesForPropertyId(
		propertyId string,
	) (
		[]dtos.DeviceDto,
		error,
	)
	AddNewDevicesToProperty(
		propertyId string,
		newDevices []dtos.NewDeviceDto,
	) (
		[]dtos.DeviceDto,
		error,
	)
}

type DeviceRepository struct {
	app *pocketbase.PocketBase
}

func (d DeviceRepository) GetDeviceById(deviceId string) (
	dtos.DeviceDto,
	error,
) {
	query := d.app.Dao().DB().Select(
		"d.id",
		"d.name",
		"d.identifier",
		"d.type",
		"d.defunctreason",
		"d.property",
	).From(
		"devices d",
	).Where(
		dbx.NewExp("d.id = {:deviceId}", dbx.Params{"deviceId": deviceId}),
	)

	var result dtos.DeviceDto

	queryErr := query.One(&result)
	return result, queryErr
}

func (d DeviceRepository) GetDeviceForPropertyId(
	propertyId string,
	inboundDevices []dtos.NewDeviceDto,
) (
	[]dtos.DeviceDto,
	error,
) {
	var deviceIdentifiersAsInterfaces []interface{}
	for _, device := range inboundDevices {
		deviceIdentifiersAsInterfaces = append(deviceIdentifiersAsInterfaces, device.Identifier)
	}
	var result []dtos.DeviceDto

	query := d.app.Dao().DB().Select(
		"d.id",
		"d.name",
		"d.identifier",
		"d.type",
		"d.defunctreason",
	).From(
		"devices d",
	).Where(
		dbx.NewExp("d.property = {:propertyId}", dbx.Params{"propertyId": propertyId}),
	).AndWhere(
		dbx.In("d.identifier", deviceIdentifiersAsInterfaces...),
	)
	queryErr := query.All(&result)
	return result, queryErr
}

func (d DeviceRepository) GetDevicesForPropertyId(
	propertyId string,
) (
	[]dtos.DeviceDto,
	error,
) {
	var result []dtos.DeviceDto

	query := d.app.Dao().DB().Select(
		"d.id",
		"d.name",
		"d.identifier",
		"d.type",
		"d.defunctreason",
	).From(
		"devices d",
	).Where(
		dbx.NewExp("d.property = {:propertyId}", dbx.Params{"propertyId": propertyId}),
	)
	queryErr := query.All(&result)
	return result, queryErr
}

func (d DeviceRepository) AddNewDevicesToProperty(
	propertyId string,
	newInboundDevices []dtos.NewDeviceDto,
) (
	[]dtos.DeviceDto,
	error,
) {
	var newDevices []dtos.DeviceDto

	if transactionErr := d.app.Dao().WithoutHooks().RunInTransaction(func(txDao *daos.Dao) error {
		deviceCollection, findCollectionErr := d.app.Dao().FindCollectionByNameOrId("devices")
		if findCollectionErr != nil {
			return findCollectionErr
		}
		deviceHistoryCollection, findHistoryCollectionErr := d.app.Dao().FindCollectionByNameOrId("devicehistory")
		if findHistoryCollectionErr != nil {
			return findHistoryCollectionErr
		}

		for _, newInboundDevice := range newInboundDevices {
			newDevice := models.NewRecord(deviceCollection)
			newDevice.Set("property", propertyId)
			newDevice.Set("name", newInboundDevice.Name)
			newDevice.Set("identifier", newInboundDevice.Identifier)
			newDevice.Set("type", newInboundDevice.Type)
			newDevice.Set("defunctreason", newInboundDevice.DefunctReason)
			if saveErr := txDao.WithoutHooks().SaveRecord(newDevice); saveErr != nil {
				return saveErr
			}
			newDevices = append(newDevices, dtos.DeviceDto{
				Id:            newDevice.Get("id").(string),
				Name:          newDevice.Get("name").(string),
				Identifier:    newDevice.Get("identifier").(string),
				Type:          newDevice.Get("type").(string),
				DefunctReason: newDevice.Get("defunctreason").(string),
			})
			for _, deviceHistory := range newInboundDevice.Histories {
				newDeviceHistory := models.NewRecord(deviceHistoryCollection)
				newDeviceHistory.Set("description", deviceHistory.Description)
				fauxDeviceJson, mashalErr := json.Marshal(dtos.DeviceDto{
					Id:            newDevice.Get("id").(string),
					Name:          newDevice.Get("name").(string),
					Identifier:    newDevice.Get("identifier").(string),
					Type:          newDevice.Get("type").(string),
					DefunctReason: newDevice.Get("defunctreason").(string),
				})
				if mashalErr != nil {
					return mashalErr
				}
				newDeviceHistory.Set("device", string(fauxDeviceJson))
				newDeviceHistory.Set("stateddatetime", deviceHistory.StatedDateTime)
				newDeviceHistory.Set("deviceid", newDevice.Get("id"))
				if saveHistoryErr := txDao.WithoutHooks().SaveRecord(newDeviceHistory); saveHistoryErr != nil {
					return saveHistoryErr
				}
			}
		}

		return nil
	}); transactionErr != nil {
		return nil, transactionErr
	}

	return newDevices, nil
}

func NewDeviceRepository(app *pocketbase.PocketBase) IDeviceRepository {
	return DeviceRepository{
		app,
	}
}
