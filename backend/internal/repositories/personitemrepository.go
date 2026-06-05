package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IPersonItemRepository interface {
	GetPersonItemById(id string) (dtos.PersonItemDto, error)
	GetPersonItemsByPersonId(personId string) ([]dtos.PersonItemDto, error)
	GetPersonItemsByItemId(itemId string) ([]dtos.PersonItemDto, error)
	CreatePersonItem(personId, itemId string) (dtos.PersonItemDto, error)
	DeletePersonItem(id string) error
}

type PersonItemRepository struct {
	app core.App
}

func NewPersonItemRepository(app core.App) IPersonItemRepository {
	return &PersonItemRepository{app}
}

func (r *PersonItemRepository) GetPersonItemById(id string) (dtos.PersonItemDto, error) {
	record, err := r.app.Dao().FindRecordById("personItems", id)
	if err != nil {
		return dtos.PersonItemDto{}, err
	}
	return personItemRecordToDto(record), nil
}

func (r *PersonItemRepository) GetPersonItemsByPersonId(personId string) ([]dtos.PersonItemDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("personItems",
		dbx.NewExp("person = {:personId}", dbx.Params{"personId": personId}))
	if err != nil {
		return nil, err
	}
	return personItemRecordsToDtos(records), nil
}

func (r *PersonItemRepository) GetPersonItemsByItemId(itemId string) ([]dtos.PersonItemDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("personItems",
		dbx.NewExp("item = {:itemId}", dbx.Params{"itemId": itemId}))
	if err != nil {
		return nil, err
	}
	return personItemRecordsToDtos(records), nil
}

func (r *PersonItemRepository) CreatePersonItem(personId, itemId string) (dtos.PersonItemDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("personItems")
	if err != nil {
		return dtos.PersonItemDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("person", personId)
	record.Set("item", itemId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.PersonItemDto{}, err
	}
	return personItemRecordToDto(record), nil
}

func (r *PersonItemRepository) DeletePersonItem(id string) error {
	record, err := r.app.Dao().FindRecordById("personItems", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func personItemRecordToDto(r *models.Record) dtos.PersonItemDto {
	return dtos.PersonItemDto{
		Id:     r.GetId(),
		Person: r.GetString("person"),
		Item:   r.GetString("item"),
	}
}

func personItemRecordsToDtos(records []*models.Record) []dtos.PersonItemDto {
	result := make([]dtos.PersonItemDto, 0, len(records))
	for _, r := range records {
		result = append(result, personItemRecordToDto(r))
	}
	return result
}
