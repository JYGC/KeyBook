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
	queryErr := p.app.Dao().DB().Select(
		"id",
		"address",
	).From(
		"properties",
	).InnerJoin(
		"managements",
		dbx.NewExp("properties.id = managements.property"),
	).Where(
		dbx.NewExp("managements.user = {:user}", dbx.Params{"user": userId}),
	).AndWhere(
		dbx.NewExp("properties.address = {:address}", dbx.Params{"address": propertyAddress}),
	).All(&result)
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

	managementsCollection, findManagementsCollectionErr := p.app.Dao().FindCollectionByNameOrId("managements")
	if findManagementsCollectionErr != nil {
		return dtos.PropertyIdAddressDto{}, findManagementsCollectionErr
	}

	// p.app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
	// })

	newProperty := models.NewRecord(propertiesCollection)
	newProperty.Set("address", propertyAddress)
	if savePropertyErr := p.app.Dao().SaveRecord(newProperty); savePropertyErr != nil {
		return dtos.PropertyIdAddressDto{}, savePropertyErr
	}

	newManagement := models.NewRecord(managementsCollection)
	newManagement.Set("user", userId)
	newManagement.Set("property", newProperty.Get("id"))
	if saveManagementErr := p.app.Dao().SaveRecord(newManagement); saveManagementErr != nil {
		return dtos.PropertyIdAddressDto{}, saveManagementErr
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
