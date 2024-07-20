package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type IPersonRepository interface {
	GetPersonsForPropertyIdByPersonNames(
		propertyId string,
		personNames []string,
	) (
		[]dtos.PersonIdNameDto,
		error,
	)
	GetPersonsForPropertyId(
		propertyId string,
	) (
		[]dtos.PersonIdNameDto,
		error,
	)
	AddNewPersonsToProperty(
		propertyId string,
		personNames []string,
	) (
		[]dtos.PersonIdNameDto,
		error,
	)
}

type PersonRepository struct {
	app *pocketbase.PocketBase
}

func (p PersonRepository) GetPersonsForPropertyIdByPersonNames(
	propertyId string,
	personNames []string,
) (
	[]dtos.PersonIdNameDto,
	error,
) {
	var result []dtos.PersonIdNameDto

	var propertyNamesAsInterfaces []interface{}
	for _, personName := range personNames {
		propertyNamesAsInterfaces = append(propertyNamesAsInterfaces, personName)
	}

	query := p.app.Dao().DB().Select(
		"pr.id",
		"pr.name",
	).From(
		"persons pr",
	).Where(
		dbx.NewExp("pr.property = {:propertyId}", dbx.Params{"propertyId": propertyId}),
	).AndWhere(
		dbx.In("pr.name", propertyNamesAsInterfaces...),
	)
	queryErr := query.All(&result)
	return result, queryErr
}

func (p PersonRepository) GetPersonsForPropertyId(
	propertyId string,
) (
	[]dtos.PersonIdNameDto,
	error,
) {
	var result []dtos.PersonIdNameDto

	query := p.app.Dao().DB().Select(
		"pr.id",
		"pr.name",
	).From(
		"persons pr",
	).Where(
		dbx.NewExp("pr.property = {:propertyId}", dbx.Params{"propertyId": propertyId}),
	)
	queryErr := query.All(&result)
	return result, queryErr
}

func (p PersonRepository) AddNewPersonsToProperty(
	propertyId string,
	personNames []string,
) (
	[]dtos.PersonIdNameDto,
	error,
) {
	var newPersons []dtos.PersonIdNameDto

	if transactionErr := p.app.Dao().WithoutHooks().RunInTransaction(func(txDao *daos.Dao) error {
		personsCollection, findCollectionErr := p.app.Dao().FindCollectionByNameOrId("persons")
		if findCollectionErr != nil {
			return findCollectionErr
		}

		for _, personName := range personNames {
			newPerson := models.NewRecord(personsCollection)
			newPerson.Set("name", personName)
			newPerson.Set("type", "Tenant")
			newPerson.Set("property", propertyId)
			if saveErr := txDao.WithoutHooks().SaveRecord(newPerson); saveErr != nil {
				return saveErr
			}
			newPersons = append(newPersons, dtos.PersonIdNameDto{
				Id:   newPerson.Get("id").(string),
				Name: newPerson.Get("name").(string),
			})
		}

		return nil
	}); transactionErr != nil {
		return nil, transactionErr
	}

	return newPersons, nil
}

func NewPersonRepository(app *pocketbase.PocketBase) IPersonRepository {
	personRepository := PersonRepository{
		app,
	}
	return personRepository
}
