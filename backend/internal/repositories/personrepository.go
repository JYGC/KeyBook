package repositories

import (
	"keybook/backend/internal/helpers"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IPersonRepository interface {
	GetPersonForUser(loggedInUser *models.Record) (*models.Record, error)
	CreatePersonForUser(userId string) (string, error)
}

type PersonRepository struct {
	app *pocketbase.PocketBase
}

func (p PersonRepository) GetPersonForUser(loggedInUser *models.Record) (*models.Record, error) {
	return p.app.Dao().FindFirstRecordByFilter(
		"persons",
		"user = {:user}",
		dbx.Params{"user": loggedInUser.Get("id")},
	)
}

func (p PersonRepository) CreatePersonForUser(userId string) (string, error) {
	user, findUserErr := p.app.Dao().FindRecordById("users", userId)
	if findUserErr != nil {
		return "", findUserErr
	}

	personsCollection, findCollectionErr := p.app.Dao().FindCollectionByNameOrId("persons")
	if findCollectionErr != nil && !helpers.IsNoRowsResult(findCollectionErr) {
		return "", findCollectionErr
	}
	newPerson := models.NewRecord(personsCollection)
	newPerson.Set("name", user.Get("name"))
	newPerson.Set("user", user.Get("id"))
	if addPersonError := p.app.Dao().SaveRecord(newPerson); addPersonError != nil {
		return "", addPersonError
	}

	return newPerson.Get("id").(string), nil
}

func NewPersonRepository(app *pocketbase.PocketBase) IPersonRepository {
	personRepository := PersonRepository{
		app,
	}
	return personRepository
}
