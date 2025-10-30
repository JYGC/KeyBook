package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type IDeviceRepository interface {
	GetDeviceById(deviceId string) (
		dtos.DeviceDto,
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
		"d.defunctReason",
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

func NewDeviceRepository(app *pocketbase.PocketBase) IDeviceRepository {
	return DeviceRepository{
		app,
	}
}
