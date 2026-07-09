package repositories

import (
	"keybook/backend/internal/dtos"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type IPropertyItemRepository interface {
	GetPropertyItemById(id string) (dtos.PropertyItemDto, error)
	GetPropertyItemsByPropertyId(propertyId string) ([]dtos.PropertyItemDto, error)
	GetPropertyItemsByItemId(itemId string) ([]dtos.PropertyItemDto, error)
	CreatePropertyItem(itemId, propertyId string) (dtos.PropertyItemDto, error)
	DeletePropertyItem(id string) error
}

type PropertyItemRepository struct {
	app core.App
}

func NewPropertyItemRepository(app core.App) IPropertyItemRepository {
	return &PropertyItemRepository{app}
}

func (r *PropertyItemRepository) GetPropertyItemById(id string) (dtos.PropertyItemDto, error) {
	record, err := r.app.Dao().FindRecordById("propertyItems", id)
	if err != nil {
		return dtos.PropertyItemDto{}, err
	}
	return propertyItemRecordToDto(record), nil
}

func (r *PropertyItemRepository) GetPropertyItemsByPropertyId(propertyId string) ([]dtos.PropertyItemDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("propertyItems",
		dbx.NewExp("property = {:propertyId}", dbx.Params{"propertyId": propertyId}))
	if err != nil {
		return nil, err
	}
	return propertyItemRecordsToDtos(records), nil
}

func (r *PropertyItemRepository) GetPropertyItemsByItemId(itemId string) ([]dtos.PropertyItemDto, error) {
	records, err := r.app.Dao().FindRecordsByExpr("propertyItems",
		dbx.NewExp("item = {:itemId}", dbx.Params{"itemId": itemId}))
	if err != nil {
		return nil, err
	}
	return propertyItemRecordsToDtos(records), nil
}

func (r *PropertyItemRepository) CreatePropertyItem(itemId, propertyId string) (dtos.PropertyItemDto, error) {
	col, err := r.app.Dao().FindCollectionByNameOrId("propertyItems")
	if err != nil {
		return dtos.PropertyItemDto{}, err
	}
	record := models.NewRecord(col)
	record.Set("item", itemId)
	record.Set("property", propertyId)
	if err := r.app.Dao().SaveRecord(record); err != nil {
		return dtos.PropertyItemDto{}, err
	}
	return propertyItemRecordToDto(record), nil
}

func (r *PropertyItemRepository) DeletePropertyItem(id string) error {
	record, err := r.app.Dao().FindRecordById("propertyItems", id)
	if err != nil {
		return err
	}
	return r.app.Dao().DeleteRecord(record)
}

func propertyItemRecordToDto(r *models.Record) dtos.PropertyItemDto {
	return dtos.PropertyItemDto{
		Id:       r.GetId(),
		Item:     r.GetString("item"),
		Property: r.GetString("property"),
	}
}

func propertyItemRecordsToDtos(records []*models.Record) []dtos.PropertyItemDto {
	result := make([]dtos.PropertyItemDto, 0, len(records))
	for _, r := range records {
		result = append(result, propertyItemRecordToDto(r))
	}
	return result
}
