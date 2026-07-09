package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type ICobrandPropertyManagerRepository interface {
	GetCobrandPropertyManagerById(id string) (dtos.CobrandPropertyManagerDto, error)
	GetCobrandPropertyManagersByPropertyId(propertyId string) ([]dtos.CobrandPropertyManagerDto, error)
	CreateCobrandPropertyManager(cobrandId, propertyId string) (dtos.CobrandPropertyManagerDto, error)
	DeleteCobrandPropertyManager(id string) error
}

type CobrandPropertyManagerRepository struct {
	app core.App
}

func NewCobrandPropertyManagerRepository(app core.App) ICobrandPropertyManagerRepository {
	return &CobrandPropertyManagerRepository{app}
}

func (r *CobrandPropertyManagerRepository) GetCobrandPropertyManagerById(id string) (dtos.CobrandPropertyManagerDto, error) {
	record, err := r.app.Dao().FindRecordById("cobrandPropertyManagers", id)
	if err != nil {
		return dtos.CobrandPropertyManagerDto{}, err
	}
	return cobrandPropertyManagerRecordToDto(record), nil
}

func (r *CobrandPropertyManagerRepository) GetCobrandPropertyManagersByPropertyId(propertyId string) ([]dtos.CobrandPropertyManagerDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("cobrandPropertyManagers",
		dbx.NewExp("property = {:propertyId}", dbx.Params{"propertyId": propertyId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.CobrandPropertyManagerDto, 0, len(records))
	for _, record := range records {
		result = append(result, cobrandPropertyManagerRecordToDto(record))
	}
	return result, nil
}

func (r *CobrandPropertyManagerRepository) CreateCobrandPropertyManager(cobrandId, propertyId string) (dtos.CobrandPropertyManagerDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("cobrandPropertyManagers")
	if err != nil {
		return dtos.CobrandPropertyManagerDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("cobrand", cobrandId)
	record.Set("property", propertyId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.CobrandPropertyManagerDto{}, err
	}
	return cobrandPropertyManagerRecordToDto(record), nil
}

func (r *CobrandPropertyManagerRepository) DeleteCobrandPropertyManager(id string) error {
	record, err := r.app.Dao().FindRecordById("cobrandPropertyManagers", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func cobrandPropertyManagerRecordToDto(r *models.Record) dtos.CobrandPropertyManagerDto {
	return dtos.CobrandPropertyManagerDto{
		Id:       r.GetId(),
		Cobrand:  r.GetString("cobrand"),
		Property: r.GetString("property"),
	}
}
