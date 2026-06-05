package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type ICobrandPropertyOwnerRepository interface {
	GetCobrandPropertyOwnerById(id string) (dtos.CobrandPropertyOwnerDto, error)
	GetCobrandPropertyOwnersByPropertyOwnerId(propertyOwnerId string) ([]dtos.CobrandPropertyOwnerDto, error)
	CreateCobrandPropertyOwner(cobrandId, propertyOwnerId string) (dtos.CobrandPropertyOwnerDto, error)
	DeleteCobrandPropertyOwner(id string) error
}

type CobrandPropertyOwnerRepository struct {
	app core.App
}

func NewCobrandPropertyOwnerRepository(app core.App) ICobrandPropertyOwnerRepository {
	return &CobrandPropertyOwnerRepository{app}
}

func (r *CobrandPropertyOwnerRepository) GetCobrandPropertyOwnerById(id string) (dtos.CobrandPropertyOwnerDto, error) {
	record, err := r.app.Dao().FindRecordById("cobrandPropertyOwners", id)
	if err != nil {
		return dtos.CobrandPropertyOwnerDto{}, err
	}
	return cobrandPropertyOwnerRecordToDto(record), nil
}

func (r *CobrandPropertyOwnerRepository) GetCobrandPropertyOwnersByPropertyOwnerId(propertyOwnerId string) ([]dtos.CobrandPropertyOwnerDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("cobrandPropertyOwners",
		dbx.NewExp("propertyOwner = {:propertyOwnerId}", dbx.Params{"propertyOwnerId": propertyOwnerId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.CobrandPropertyOwnerDto, 0, len(records))
	for _, record := range records {
		result = append(result, cobrandPropertyOwnerRecordToDto(record))
	}
	return result, nil
}

func (r *CobrandPropertyOwnerRepository) CreateCobrandPropertyOwner(cobrandId, propertyOwnerId string) (dtos.CobrandPropertyOwnerDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("cobrandPropertyOwners")
	if err != nil {
		return dtos.CobrandPropertyOwnerDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("cobrand", cobrandId)
	record.Set("propertyOwner", propertyOwnerId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.CobrandPropertyOwnerDto{}, err
	}
	return cobrandPropertyOwnerRecordToDto(record), nil
}

func (r *CobrandPropertyOwnerRepository) DeleteCobrandPropertyOwner(id string) error {
	record, err := r.app.Dao().FindRecordById("cobrandPropertyOwners", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func cobrandPropertyOwnerRecordToDto(r *models.Record) dtos.CobrandPropertyOwnerDto {
	return dtos.CobrandPropertyOwnerDto{
		Id:            r.GetId(),
		Cobrand:       r.GetString("cobrand"),
		PropertyOwner: r.GetString("propertyOwner"),
	}
}
