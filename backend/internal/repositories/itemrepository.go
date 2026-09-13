package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IItemRepository interface {
	GetItemById(id string) (dtos.ItemDto, error)
	CreateItem(name, description string) (dtos.ItemDto, error)
	UpdateItem(id, name, description string) error
	DeleteItem(id string) error
}

type ItemRepository struct {
	app core.App
}

func NewItemRepository(app core.App) IItemRepository {
	return &ItemRepository{app}
}

func (r *ItemRepository) GetItemById(id string) (dtos.ItemDto, error) {
	record, err := r.app.Dao().FindRecordById("items", id)
	if err != nil {
		return dtos.ItemDto{}, err
	}
	return itemRecordToDto(record), nil
}

func (r *ItemRepository) CreateItem(name, description string) (dtos.ItemDto, error) {
	itemsCollection, err := r.app.Dao().FindCollectionByNameOrId("items")
	if err != nil {
		return dtos.ItemDto{}, err
	}
	record := models.NewRecord(itemsCollection)
	record.Set("name", name)
	record.Set("description", description)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.ItemDto{}, err
	}
	return itemRecordToDto(record), nil
}

func (r *ItemRepository) UpdateItem(id, name, description string) error {
	record, err := r.app.Dao().FindRecordById("items", id)
	if err != nil {
		return err
	}
	record.Set("name", name)
	record.Set("description", description)
	return r.app.Dao().SaveRecord(record)
}

func (r *ItemRepository) DeleteItem(id string) error {
	record, err := r.app.Dao().FindRecordById("items", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func itemRecordToDto(record *models.Record) dtos.ItemDto {
	return dtos.ItemDto{
		Id:          record.GetId(),
		Name:        record.GetString("name"),
		Description: record.GetString("description"),
		Picture:     record.GetString("picture"),
	}
}
