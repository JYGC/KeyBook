package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IPropertyRepository interface {
	GetPropertyById(propertyId string) (dtos.PropertyDto, error)
	CreateProperty(address string) (dtos.PropertyDto, error)
	UpdateProperty(id, address string) error
	DeleteProperty(id string) error
}

type PropertyRepository struct {
	app core.App
}

func NewPropertyRepository(app core.App) IPropertyRepository {
	return &PropertyRepository{app}
}

func (r *PropertyRepository) GetPropertyById(propertyId string) (dtos.PropertyDto, error) {
	query := r.app.Dao().DB().Select(
		"p.id",
		"p.address",
	).From(
		"properties p",
	).Where(
		dbx.NewExp("p.id = {:propertyId}", dbx.Params{"propertyId": propertyId}),
	)

	var result dtos.PropertyDto
	return result, query.One(&result)
}

func (r *PropertyRepository) CreateProperty(address string) (dtos.PropertyDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("properties")
	if err != nil {
		return dtos.PropertyDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("address", address)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.PropertyDto{}, err
	}
	return dtos.PropertyDto{
		Id:      record.GetId(),
		Address: record.GetString("address"),
	}, nil
}

func (r *PropertyRepository) UpdateProperty(id, address string) error {
	record, err := r.app.Dao().FindRecordById("properties", id)
	if err != nil {
		return err
	}
	record.Set("address", address)
	return r.app.Dao().SaveRecord(record)
}

func (r *PropertyRepository) DeleteProperty(id string) error {
	record, err := r.app.Dao().FindRecordById("properties", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}
