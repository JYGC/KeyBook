package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IPropertyRepository interface {
	GetPropertiesForUserByPropertyName(userId string, propertyAddress string) ([]dtos.PropertyIdAddressDto, error)
	IsPropertyIdBelongToUser(userId string, propertyId string) (bool, error)
	AddNewProperty(userId string, propertyAddress string) (dtos.PropertyIdAddressDto, error)
}

type PropertyRepository struct {
	app *pocketbase.PocketBase
}

func (p PropertyRepository) GetPropertiesForUserByPropertyName(
	userId string,
	propertyAddress string,
) (
	[]dtos.PropertyIdAddressDto,
	error,
) {
	var result []dtos.PropertyIdAddressDto
	query := p.app.Dao().WithoutHooks().DB().Select(
		"pt.id",
		"pt.address",
	).From(
		"properties pt, json_each(owners) o",
	).Where(
		dbx.NewExp("o.value = {:userId}", dbx.Params{"userId": userId}),
	).AndWhere(
		dbx.NewExp("pt.address = {:address}", dbx.Params{"address": propertyAddress}),
	)
	queryErr := query.All(&result)
	return result, queryErr
}

func (p PropertyRepository) IsPropertyIdBelongToUser(
	userId string,
	propertyId string,
) (
	bool,
	error,
) {
	var result []dtos.PropertyIdAddressDto
	query := p.app.Dao().WithoutHooks().DB().Select(
		"pt.id",
		"pt.address",
	).From(
		"properties pt, json_each(owners) o",
	).Where(
		dbx.NewExp("o.value = {:userId}", dbx.Params{"userId": userId}),
	).AndWhere(
		dbx.NewExp("pt.id = {:propertyId}", dbx.Params{"propertyId": propertyId}),
	).Limit(1)

	queryErr := query.All(&result)
	return len(result) == 1, queryErr
}

func (p PropertyRepository) AddNewProperty(userId string, propertyAddress string) (dtos.PropertyIdAddressDto, error) {
	propertiesCollection, findManagementsCollectionErr := p.app.Dao().FindCollectionByNameOrId("properties")
	if findManagementsCollectionErr != nil {
		return dtos.PropertyIdAddressDto{}, findManagementsCollectionErr
	}

	newProperty := models.NewRecord(propertiesCollection)
	newProperty.Set("address", propertyAddress)
	newProperty.Set("owners", []string{userId})
	if savePropertyErr := p.app.Dao().WithoutHooks().SaveRecord(newProperty); savePropertyErr != nil {
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
