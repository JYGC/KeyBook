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
	GetHoldingPersonNameByDeviceId(
		personDeviceId string,
	) (
		dtos.PersonIdNameDto,
		error,
	)
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

func (pdr PersonDeviceRepository) GetHoldingPersonNameByDeviceId(
	personDeviceId string,
) (
	dtos.PersonIdNameDto,
	error,
) {
	query := pdr.app.Dao().DB().Select(
		"p.id",
		"p.name",
	).From(
		"persondevices pd",
	).InnerJoin(
		"persons p",
		dbx.NewExp("pd.person = p.id"),
	).Where(
		dbx.NewExp("d.id = {:personDeviceId}", dbx.Params{"personDeviceId": personDeviceId}),
	)

	var result dtos.PersonIdNameDto

	queryErr := query.One(&result)
	return result, queryErr
}

func (pdr PersonDeviceRepository) GetPersonDevicesForProperty(
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

	query := pdr.app.Dao().DB().Select(
	//
	).From(
		"persondevices pd",
	).Where(
		dbx.In("pd.device", deviceIdsInUseAsInterfaces...),
	)
	queryErr := query.All(&result)
	return result, queryErr
}

func (pdr PersonDeviceRepository) AddNewPersonDevices(
	personDevicesToAdd []dtos.NewPersonDeviceAndHistoriesDto,
) (
	[]dtos.PersonDeviceDto,
	error,
) {
	var newPersonDevices []dtos.PersonDeviceDto

	if transactionErr := pdr.app.Dao().WithoutHooks().RunInTransaction(func(txDao *daos.Dao) error {
		personDeviceCollection, findCollectionErr := pdr.app.Dao().FindCollectionByNameOrId("persondevices")
		if findCollectionErr != nil {
			return findCollectionErr
		}
		personDeviceHistoryCollection, findHistoryCollectionErr := pdr.app.Dao().FindCollectionByNameOrId("persondevicehistory")
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
