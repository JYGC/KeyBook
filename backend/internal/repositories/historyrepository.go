package repositories

import (
	"encoding/json"

	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

type IHistoryRepository interface {
	AddHistories(collectionName string, models []interface{}) error
	GetHistories(collectionName string) ([]interface{}, error)
}

type HistoryRepository struct {
	historyDb *clover.DB
}

func (p HistoryRepository) AddHistories(collectionName string, models []interface{}) error {
	newDocuments := []*document.Document{}
	for _, model := range models {
		modelJson, marshalErr := json.Marshal(model)
		if marshalErr != nil {
			return marshalErr
		}

		var modelMap map[string]interface{}
		json.Unmarshal([]byte(modelJson), &modelMap)

		newDocument := document.NewDocument()
		newDocument.SetAll(modelMap)
		newDocuments = append(newDocuments, newDocument)
	}
	insertErr := p.historyDb.Insert(collectionName, newDocuments...)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (p HistoryRepository) GetHistories(collectionName string) ([]interface{}, error) {
	documents, findAllErr := p.historyDb.FindAll(query.NewQuery(collectionName))
	if findAllErr != nil {
		return nil, findAllErr
	}

	models := []interface{}{}
	for _, document := range documents {
		model := document.ToMap()
		models = append(models, model)
	}
	return models, nil
}

func NewHistoryRepository(historyDb *clover.DB) IHistoryRepository {
	return HistoryRepository{
		historyDb,
	}
}
