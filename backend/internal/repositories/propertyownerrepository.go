package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IPropertyOwnerRepository interface {
	GetPropertyOwnerById(id string) (dtos.PropertyOwnerDto, error)
	GetPropertyOwnersByPropertyId(propertyId string) ([]dtos.PropertyOwnerDto, error)
	CreatePropertyOwner(propertyId string) (dtos.PropertyOwnerDto, error)
	DeletePropertyOwner(id string) error
}

type PropertyOwnerRepository struct {
	app core.App
}

func NewPropertyOwnerRepository(app core.App) IPropertyOwnerRepository {
	return &PropertyOwnerRepository{app}
}

func (r *PropertyOwnerRepository) GetPropertyOwnerById(id string) (dtos.PropertyOwnerDto, error) {
	record, err := r.app.Dao().FindRecordById("propertyOwners", id)
	if err != nil {
		return dtos.PropertyOwnerDto{}, err
	}
	return propertyOwnerRecordToDto(record), nil
}

func (r *PropertyOwnerRepository) GetPropertyOwnersByPropertyId(propertyId string) ([]dtos.PropertyOwnerDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("propertyOwners",
		dbx.NewExp("property = {:propertyId}", dbx.Params{"propertyId": propertyId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.PropertyOwnerDto, 0, len(records))
	for _, record := range records {
		result = append(result, propertyOwnerRecordToDto(record))
	}
	return result, nil
}

func (r *PropertyOwnerRepository) CreatePropertyOwner(propertyId string) (dtos.PropertyOwnerDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("propertyOwners")
	if err != nil {
		return dtos.PropertyOwnerDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("property", propertyId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.PropertyOwnerDto{}, err
	}
	return propertyOwnerRecordToDto(record), nil
}

func (r *PropertyOwnerRepository) DeletePropertyOwner(id string) error {
	record, err := r.app.Dao().FindRecordById("propertyOwners", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func propertyOwnerRecordToDto(r *models.Record) dtos.PropertyOwnerDto {
	return dtos.PropertyOwnerDto{
		Id:       r.GetId(),
		Property: r.GetString("property"),
	}
}
