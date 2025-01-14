package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type IPersonRepository interface {
	GetPersonById(personId string) (
		dtos.PersonDto,
		error,
	)
}

type PersonRepository struct {
	app *pocketbase.PocketBase
}

func (p PersonRepository) GetPersonById(personId string) (
	dtos.PersonDto,
	error,
) {
	query := p.app.Dao().DB().Select(
		"p.id",
		"p.name",
		"p.type",
		"p.property",
	).From(
		"persons p",
	).Where(
		dbx.NewExp("p.id = {:personId}", dbx.Params{"personId": personId}),
	)

	var result dtos.PersonDto

	queryErr := query.One(&result)
	return result, queryErr
}

func NewPersonRepository(app *pocketbase.PocketBase) IPersonRepository {
	personRepository := PersonRepository{
		app,
	}
	return personRepository
}
