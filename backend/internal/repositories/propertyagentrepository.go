package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IPropertyAgentRepository interface {
	GetPropertyAgentById(id string) (dtos.PropertyAgentDto, error)
	GetPropertyAgentsByPropertyId(propertyId string) ([]dtos.PropertyAgentDto, error)
	CreatePropertyAgent(agentId, propertyId string) (dtos.PropertyAgentDto, error)
	DeletePropertyAgent(id string) error
}

type PropertyAgentRepository struct {
	app core.App
}

func NewPropertyAgentRepository(app core.App) IPropertyAgentRepository {
	return &PropertyAgentRepository{app}
}

func (r *PropertyAgentRepository) GetPropertyAgentById(id string) (dtos.PropertyAgentDto, error) {
	record, err := r.app.Dao().FindRecordById("propertyAgents", id)
	if err != nil {
		return dtos.PropertyAgentDto{}, err
	}
	return propertyAgentRecordToDto(record), nil
}

func (r *PropertyAgentRepository) GetPropertyAgentsByPropertyId(propertyId string) ([]dtos.PropertyAgentDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("propertyAgents",
		dbx.NewExp("property = {:propertyId}", dbx.Params{"propertyId": propertyId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.PropertyAgentDto, 0, len(records))
	for _, record := range records {
		result = append(result, propertyAgentRecordToDto(record))
	}
	return result, nil
}

func (r *PropertyAgentRepository) CreatePropertyAgent(agentId, propertyId string) (dtos.PropertyAgentDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("propertyAgents")
	if err != nil {
		return dtos.PropertyAgentDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("agent", agentId)
	record.Set("property", propertyId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.PropertyAgentDto{}, err
	}
	return propertyAgentRecordToDto(record), nil
}

func (r *PropertyAgentRepository) DeletePropertyAgent(id string) error {
	record, err := r.app.Dao().FindRecordById("propertyAgents", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func propertyAgentRecordToDto(r *models.Record) dtos.PropertyAgentDto {
	return dtos.PropertyAgentDto{
		Id:       r.GetId(),
		Agent:    r.GetString("agent"),
		Property: r.GetString("property"),
	}
}
