package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type ICobrandRepository interface {
	GetCobrandById(id string) (dtos.CobrandDto, error)
	CreateCobrand(name string) (dtos.CobrandDto, error)
	UpdateCobrand(id, name string) error
	DeleteCobrand(id string) error
}

type CobrandRepository struct {
	app core.App
}

func NewCobrandRepository(app core.App) ICobrandRepository {
	return &CobrandRepository{app}
}

func (r *CobrandRepository) GetCobrandById(id string) (dtos.CobrandDto, error) {
	record, err := r.app.Dao().FindRecordById("cobrands", id)
	if err != nil {
		return dtos.CobrandDto{}, err
	}
	return cobrandRecordToDto(record), nil
}

func (r *CobrandRepository) CreateCobrand(name string) (dtos.CobrandDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("cobrands")
	if err != nil {
		return dtos.CobrandDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("name", name)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.CobrandDto{}, err
	}
	return cobrandRecordToDto(record), nil
}

func (r *CobrandRepository) UpdateCobrand(id, name string) error {
	record, err := r.app.Dao().FindRecordById("cobrands", id)
	if err != nil {
		return err
	}
	record.Set("name", name)
	return r.app.Dao().SaveRecord(record)
}

func (r *CobrandRepository) DeleteCobrand(id string) error {
	record, err := r.app.Dao().FindRecordById("cobrands", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func cobrandRecordToDto(r *models.Record) dtos.CobrandDto {
	return dtos.CobrandDto{
		Id:   r.GetId(),
		Name: r.GetString("name"),
	}
}
