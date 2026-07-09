package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IPersonRepository interface {
	GetPersonById(personId string) (dtos.PersonDto, error)
	CreatePerson(name, DOB, userID string) (dtos.PersonDto, error)
	UpdatePerson(id, name, DOB string) error
	DeletePerson(id string) error
}

type PersonRepository struct {
	app core.App
}

func NewPersonRepository(app core.App) IPersonRepository {
	return &PersonRepository{app}
}

func (r *PersonRepository) GetPersonById(personId string) (dtos.PersonDto, error) {
	query := r.app.Dao().DB().Select(
		"p.id",
		"p.name",
		"p.DOB",
		"p.user",
		"p.profileImage",
	).From(
		"persons p",
	).Where(
		dbx.NewExp("p.id = {:personId}", dbx.Params{"personId": personId}),
	)

	var result dtos.PersonDto
	return result, query.One(&result)
}

func (r *PersonRepository) CreatePerson(name, DOB, userID string) (dtos.PersonDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("persons")
	if err != nil {
		return dtos.PersonDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("name", name)
	record.Set("DOB", DOB)
	if userID != "" {
		record.Set("user", userID)
	}
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.PersonDto{}, err
	}
	return dtos.PersonDto{
		Id:   record.GetId(),
		Name: record.GetString("name"),
		DOB:  record.GetString("DOB"),
		User: record.GetString("user"),
	}, nil
}

func (r *PersonRepository) UpdatePerson(id, name, DOB string) error {
	record, err := r.app.Dao().FindRecordById("persons", id)
	if err != nil {
		return err
	}
	record.Set("name", name)
	record.Set("DOB", DOB)
	return r.app.Dao().SaveRecord(record)
}

func (r *PersonRepository) DeletePerson(id string) error {
	record, err := r.app.Dao().FindRecordById("persons", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}
