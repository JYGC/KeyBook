package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IAgentRepository interface {
	GetAgentById(id string) (dtos.AgentDto, error)
	GetAgentsByCobrandId(cobrandId string) ([]dtos.AgentDto, error)
	CreateAgent(personId, cobrandId string) (dtos.AgentDto, error)
	DeleteAgent(id string) error
}

type AgentRepository struct {
	app core.App
}

func NewAgentRepository(app core.App) IAgentRepository {
	return &AgentRepository{app}
}

func (r *AgentRepository) GetAgentById(id string) (dtos.AgentDto, error) {
	record, err := r.app.Dao().FindRecordById("agents", id)
	if err != nil {
		return dtos.AgentDto{}, err
	}
	return agentRecordToDto(record), nil
}

func (r *AgentRepository) GetAgentsByCobrandId(cobrandId string) ([]dtos.AgentDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("agents",
		dbx.NewExp("cobrand = {:cobrandId}", dbx.Params{"cobrandId": cobrandId}))
	if err != nil {
		return nil, err
	}
	result := make([]dtos.AgentDto, 0, len(records))
	for _, record := range records {
		result = append(result, agentRecordToDto(record))
	}
	return result, nil
}

func (r *AgentRepository) CreateAgent(personId, cobrandId string) (dtos.AgentDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("agents")
	if err != nil {
		return dtos.AgentDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("person", personId)
	record.Set("cobrand", cobrandId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.AgentDto{}, err
	}
	return agentRecordToDto(record), nil
}

func (r *AgentRepository) DeleteAgent(id string) error {
	record, err := r.app.Dao().FindRecordById("agents", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func agentRecordToDto(r *models.Record) dtos.AgentDto {
	return dtos.AgentDto{
		Id:      r.GetId(),
		Person:  r.GetString("person"),
		Cobrand: r.GetString("cobrand"),
	}
}
