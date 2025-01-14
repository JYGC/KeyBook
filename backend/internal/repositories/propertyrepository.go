package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type IPropertyRepository interface {
	GetFullPropertyById(
		propertyId string,
	) (
		dtos.PropertyIdAddressDto,
		error,
	)
}

type PropertyRepository struct {
	app *pocketbase.PocketBase
}

func (p PropertyRepository) GetFullPropertyById(
	propertyId string,
) (
	dtos.PropertyIdAddressDto,
	error,
) {
	query := p.app.Dao().DB().Select(
		"p.id",
		"p.address",
		"p.owners",
		"p.managers",
	).From(
		"properties p",
	).Where(
		dbx.NewExp("p.id = {:propertyId}", dbx.Params{"propertyId": propertyId}),
	)

	var result dtos.PropertyIdAddressDto

	queryErr := query.One(&result)
	return result, queryErr
}

func NewPropertyRepository(app *pocketbase.PocketBase) IPropertyRepository {
	propertyRepository := PropertyRepository{
		app,
	}
	return propertyRepository
}
