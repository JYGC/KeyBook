package repositories

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/helpers"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IPropertyRepository interface {
	GetPropertiesManagedByUser(userId string, propertyAddress string) ([]dtos.PropertyIdAddressDto, error)
	GetPropertyIdByName(propertyAddress string) (string, error)
	AddNewProperty(userId string, propertyAddress string) (dtos.PropertyIdAddressDto, error)
}

type PropertyRepository struct {
	app *pocketbase.PocketBase
}

func (p PropertyRepository) GetPropertiesManagedByUser(
	userId string,
	propertyAddress string,
) (
	[]dtos.PropertyIdAddressDto,
	error,
) {
	var result []dtos.PropertyIdAddressDto
	query := p.app.Dao().DB().Select(
		"p.id",
		"p.address",
	).From(
		"properties p, json_each(owners) o",
	).Where(
		dbx.NewExp("o.value = {:userId}", dbx.Params{"userId": userId}),
	).AndWhere(
		dbx.NewExp("p.address = {:address}", dbx.Params{"address": propertyAddress}),
	)
	queryErr := query.All(&result)
	return result, queryErr
}

func (p PropertyRepository) GetPropertyIdByName(propertyAddress string) (string, error) {
	property, err := p.app.Dao().FindFirstRecordByFilter(
		"properties",
		"address = {:address}",
		dbx.Params{"address": propertyAddress},
	)
	if err != nil && !helpers.IsNoRowsResult(err) {
		return "", err
	}
	return property.Get("id").(string), nil
}

func (p PropertyRepository) AddNewProperty(userId string, propertyAddress string) (dtos.PropertyIdAddressDto, error) {
	propertiesCollection, findManagementsCollectionErr := p.app.Dao().FindCollectionByNameOrId("properties")
	if findManagementsCollectionErr != nil {
		return dtos.PropertyIdAddressDto{}, findManagementsCollectionErr
	}

	newProperty := models.NewRecord(propertiesCollection)
	newProperty.Set("address", propertyAddress)
	newProperty.Set("owners", []string{userId})
	if savePropertyErr := p.app.Dao().SaveRecord(newProperty); savePropertyErr != nil {
		return dtos.PropertyIdAddressDto{}, savePropertyErr
	}

	return dtos.PropertyIdAddressDto{
		Id:      newProperty.Get("id").(string),
		Address: newProperty.Get("address").(string),
	}, nil
}

func NewPropertyRepository(app *pocketbase.PocketBase) IPropertyRepository {
	propertyRepository := PropertyRepository{
		app,
	}
	return propertyRepository
}
