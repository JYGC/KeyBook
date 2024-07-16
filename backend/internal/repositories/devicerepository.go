package repositories

import (
	"github.com/pocketbase/pocketbase"
)

type IDeviceRepository interface {
}

type DeviceRepository struct {
	app *pocketbase.PocketBase
}

func NewDeviceRepository(app *pocketbase.PocketBase) IPersonDeviceRepository {
	personDeviceRepository := PersonDeviceRepository{
		app,
	}
	return personDeviceRepository
}
