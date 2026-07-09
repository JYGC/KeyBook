package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IPersonPropertyOwnerRepository interface {
	GetPersonPropertyOwnerById(id string) (dtos.PersonPropertyOwnerDto, error)
	GetPersonPropertyOwnersByPropertyOwnerId(propertyOwnerId string) ([]dtos.PersonPropertyOwnerDto, error)
	CreatePersonPropertyOwner(personId, propertyOwnerId string) (dtos.PersonPropertyOwnerDto, error)
	DeletePersonPropertyOwner(id string) error
}

type PersonPropertyOwnerRepository struct {
	app core.App
}

func NewPersonPropertyOwnerRepository(app core.App) IPersonPropertyOwnerRepository {
	return &PersonPropertyOwnerRepository{app}
}

func (r *PersonPropertyOwnerRepository) GetPersonPropertyOwnerById(id string) (dtos.PersonPropertyOwnerDto, error) {
	record, err := r.app.Dao().FindRecordById("personPropertyOwners", id)
	if err != nil {
		return dtos.PersonPropertyOwnerDto{}, err
	}
	return personPropertyOwnerRecordToDto(record), nil
}

func (r *PersonPropertyOwnerRepository) GetPersonPropertyOwnersByPropertyOwnerId(propertyOwnerId string) ([]dtos.PersonPropertyOwnerDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("personPropertyOwners",
		dbx.NewExp("propertyOwner = {:propertyOwnerId}", dbx.Params{"propertyOwnerId": propertyOwnerId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.PersonPropertyOwnerDto, 0, len(records))
	for _, record := range records {
		result = append(result, personPropertyOwnerRecordToDto(record))
	}
	return result, nil
}

func (r *PersonPropertyOwnerRepository) CreatePersonPropertyOwner(personId, propertyOwnerId string) (dtos.PersonPropertyOwnerDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("personPropertyOwners")
	if err != nil {
		return dtos.PersonPropertyOwnerDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("person", personId)
	record.Set("propertyOwner", propertyOwnerId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.PersonPropertyOwnerDto{}, err
	}
	return personPropertyOwnerRecordToDto(record), nil
}

func (r *PersonPropertyOwnerRepository) DeletePersonPropertyOwner(id string) error {
	record, err := r.app.Dao().FindRecordById("personPropertyOwners", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func personPropertyOwnerRecordToDto(r *models.Record) dtos.PersonPropertyOwnerDto {
	return dtos.PersonPropertyOwnerDto{
		Id:            r.GetId(),
		Person:        r.GetString("person"),
		PropertyOwner: r.GetString("propertyOwner"),
	}
}
