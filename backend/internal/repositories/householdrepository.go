package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IHouseholdRepository interface {
	GetHouseholdById(id string) (dtos.HouseholdDto, error)
	GetHouseholdsByPropertyId(propertyId string) ([]dtos.HouseholdDto, error)
	CreateHousehold(personId, propertyId string) (dtos.HouseholdDto, error)
	DeleteHousehold(id string) error
}

type HouseholdRepository struct {
	app *pocketbase.PocketBase
}

func NewHouseholdRepository(app *pocketbase.PocketBase) IHouseholdRepository {
	return &HouseholdRepository{app}
}

func (r *HouseholdRepository) GetHouseholdById(id string) (dtos.HouseholdDto, error) {
	record, err := r.app.Dao().FindRecordById("households", id)
	if err != nil {
		return dtos.HouseholdDto{}, err
	}
	return householdRecordToDto(record), nil
}

func (r *HouseholdRepository) GetHouseholdsByPropertyId(propertyId string) ([]dtos.HouseholdDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("households",
		dbx.NewExp("property = {:propertyId}", dbx.Params{"propertyId": propertyId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.HouseholdDto, 0, len(records))
	for _, record := range records {
		result = append(result, householdRecordToDto(record))
	}
	return result, nil
}

func (r *HouseholdRepository) CreateHousehold(personId, propertyId string) (dtos.HouseholdDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("households")
	if err != nil {
		return dtos.HouseholdDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("person", personId)
	record.Set("property", propertyId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.HouseholdDto{}, err
	}
	return householdRecordToDto(record), nil
}

func (r *HouseholdRepository) DeleteHousehold(id string) error {
	record, err := r.app.Dao().FindRecordById("households", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func householdRecordToDto(r *models.Record) dtos.HouseholdDto {
	return dtos.HouseholdDto{
		Id:       r.GetId(),
		Person:   r.GetString("person"),
		Property: r.GetString("property"),
	}
}
