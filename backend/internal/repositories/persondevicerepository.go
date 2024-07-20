package repositories

import (
	"encoding/json"
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type IPersonDeviceRepository interface {
	GetPersonDevicesForProperty(
		personDevicesToAdd []dtos.NewPersonDeviceAndHistoriesDto,
	) (
		[]dtos.PersonDeviceDto,
		error,
	)
	AddNewPersonDevices(
		personDevicesToAdd []dtos.NewPersonDeviceAndHistoriesDto,
	) (
		[]dtos.PersonDeviceDto,
		error,
	)
}

type PersonDeviceRepository struct {
	app *pocketbase.PocketBase
}

func (p PersonDeviceRepository) GetPersonDevicesForProperty(
	personDevicesToAdd []dtos.NewPersonDeviceAndHistoriesDto,
) (
	[]dtos.PersonDeviceDto,
	error,
) {
	var result []dtos.PersonDeviceDto

	var deviceIdsInUseAsInterfaces []interface{}
	for _, deviceDeviceToAdd := range personDevicesToAdd {
		deviceIdsInUseAsInterfaces = append(
			deviceIdsInUseAsInterfaces,
			deviceDeviceToAdd.Device,
		)
	}

	query := p.app.Dao().DB().Select(
	//
	).From(
		"persondevices pd",
	).Where(
		dbx.In("pd.device", deviceIdsInUseAsInterfaces...),
	)
	queryErr := query.All(&result)
	return result, queryErr
}

func (p PersonDeviceRepository) AddNewPersonDevices(
	personDevicesToAdd []dtos.NewPersonDeviceAndHistoriesDto,
) (
	[]dtos.PersonDeviceDto,
	error,
) {
	var newPersonDevices []dtos.PersonDeviceDto

	if transactionErr := p.app.Dao().WithoutHooks().RunInTransaction(func(txDao *daos.Dao) error {
		personDeviceCollection, findCollectionErr := p.app.Dao().FindCollectionByNameOrId("persondevices")
		if findCollectionErr != nil {
			return findCollectionErr
		}
		personDeviceHistoryCollection, findHistoryCollectionErr := p.app.Dao().FindCollectionByNameOrId("persondevicehistory")
		if findHistoryCollectionErr != nil {
			return findHistoryCollectionErr
		}

		for _, personDeviceToAdd := range personDevicesToAdd {
			newPersonDevice := models.NewRecord(personDeviceCollection)
			newPersonDevice.Set("person", personDeviceToAdd.Person)
			newPersonDevice.Set("device", personDeviceToAdd.Device)
			if saveErr := txDao.WithoutHooks().SaveRecord(newPersonDevice); saveErr != nil {
				return saveErr
			}
			newPersonDevices = append(newPersonDevices, dtos.PersonDeviceDto{
				Id:     newPersonDevice.Get("id").(string),
				Person: newPersonDevice.Get("person").(string),
				Device: newPersonDevice.Get("device").(string),
			})
			for _, personDeviceHistory := range personDeviceToAdd.Histories {
				newPersonDeviceHistory := models.NewRecord(personDeviceHistoryCollection)
				newPersonDeviceHistory.Set("description", personDeviceHistory.Description)
				fauxPersonDeviceJson, marshalErr := json.Marshal(dtos.PersonDeviceDto{
					Id:     newPersonDevice.Get("id").(string),
					Person: personDeviceHistory.Person,
					Device: personDeviceHistory.Device,
				})
				if marshalErr != nil {
					return marshalErr
				}
				newPersonDeviceHistory.Set("persondevice", string(fauxPersonDeviceJson))
				newPersonDeviceHistory.Set("stateddatetime", personDeviceHistory.StatedDateTime)
				newPersonDeviceHistory.Set("persondeviceid", newPersonDevice.Get("id"))
				if saveHistoryErr := txDao.WithoutHooks().SaveRecord(newPersonDeviceHistory); saveHistoryErr != nil {
					return saveHistoryErr
				}
			}
		}

		return nil
	}); transactionErr != nil {
		return nil, transactionErr
	}

	return newPersonDevices, nil
}

func NewPersonDeviceRepository(app *pocketbase.PocketBase) IPersonDeviceRepository {
	personDeviceRepository := PersonDeviceRepository{
		app,
	}
	return personDeviceRepository
}
