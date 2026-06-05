package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type ICobrandAdminRepository interface {
	GetCobrandAdminById(id string) (dtos.CobrandAdminDto, error)
	GetCobrandAdminsByCobrandId(cobrandId string) ([]dtos.CobrandAdminDto, error)
	CreateCobrandAdmin(userId, cobrandId string) (dtos.CobrandAdminDto, error)
	DeleteCobrandAdmin(id string) error
}

type CobrandAdminRepository struct {
	app core.App
}

func NewCobrandAdminRepository(app core.App) ICobrandAdminRepository {
	return &CobrandAdminRepository{app}
}

func (r *CobrandAdminRepository) GetCobrandAdminById(id string) (dtos.CobrandAdminDto, error) {
	record, err := r.app.Dao().FindRecordById("cobrandAdmins", id)
	if err != nil {
		return dtos.CobrandAdminDto{}, err
	}
	return cobrandAdminRecordToDto(record), nil
}

func (r *CobrandAdminRepository) GetCobrandAdminsByCobrandId(cobrandId string) ([]dtos.CobrandAdminDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("cobrandAdmins",
		dbx.NewExp("cobrand = {:cobrandId}", dbx.Params{"cobrandId": cobrandId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.CobrandAdminDto, 0, len(records))
	for _, record := range records {
		result = append(result, cobrandAdminRecordToDto(record))
	}
	return result, nil
}

func (r *CobrandAdminRepository) CreateCobrandAdmin(userId, cobrandId string) (dtos.CobrandAdminDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("cobrandAdmins")
	if err != nil {
		return dtos.CobrandAdminDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("user", userId)
	record.Set("cobrand", cobrandId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.CobrandAdminDto{}, err
	}
	return cobrandAdminRecordToDto(record), nil
}

func (r *CobrandAdminRepository) DeleteCobrandAdmin(id string) error {
	record, err := r.app.Dao().FindRecordById("cobrandAdmins", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func cobrandAdminRecordToDto(r *models.Record) dtos.CobrandAdminDto {
	return dtos.CobrandAdminDto{
		Id:      r.GetId(),
		User:    r.GetString("user"),
		Cobrand: r.GetString("cobrand"),
	}
}
