package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type ITenantRepository interface {
	GetTenantById(id string) (dtos.TenantDto, error)
	GetTenantsByPropertyId(propertyId string) ([]dtos.TenantDto, error)
	CreateTenant(personId, propertyId string) (dtos.TenantDto, error)
	DeleteTenant(id string) error
}

type TenantRepository struct {
	app core.App
}

func NewTenantRepository(app core.App) ITenantRepository {
	return &TenantRepository{app}
}

func (r *TenantRepository) GetTenantById(id string) (dtos.TenantDto, error) {
	record, err := r.app.Dao().FindRecordById("tenants", id)
	if err != nil {
		return dtos.TenantDto{}, err
	}
	return tenantRecordToDto(record), nil
}

func (r *TenantRepository) GetTenantsByPropertyId(propertyId string) ([]dtos.TenantDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("tenants",
		dbx.NewExp("property = {:propertyId}", dbx.Params{"propertyId": propertyId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.TenantDto, 0, len(records))
	for _, record := range records {
		result = append(result, tenantRecordToDto(record))
	}
	return result, nil
}

func (r *TenantRepository) CreateTenant(personId, propertyId string) (dtos.TenantDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("tenants")
	if err != nil {
		return dtos.TenantDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("person", personId)
	record.Set("property", propertyId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.TenantDto{}, err
	}
	return tenantRecordToDto(record), nil
}

func (r *TenantRepository) DeleteTenant(id string) error {
	record, err := r.app.Dao().FindRecordById("tenants", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func tenantRecordToDto(r *models.Record) dtos.TenantDto {
	return dtos.TenantDto{
		Id:       r.GetId(),
		Person:   r.GetString("person"),
		Property: r.GetString("property"),
	}
}
