package repositories

import (
	"github.com/pocketbase/pocketbase"
)

type IPersonDeviceRepository interface {
}

type PersonDeviceRepository struct {
	app *pocketbase.PocketBase
}

func NewPersonDeviceRepository(app *pocketbase.PocketBase) IPersonDeviceRepository {
	personDeviceRepository := PersonDeviceRepository{
		app,
	}
	return personDeviceRepository
}
