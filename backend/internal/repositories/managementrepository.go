package repositories

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type IManagementRepository interface{}

type ManagementRepository struct {
	app *pocketbase.PocketBase
}

func (o ManagementRepository) AddManagement(userId string, propertyId string) (string, error) {
	managementCollection, findCollectionErr := o.app.Dao().FindCollectionByNameOrId("managements")
	if findCollectionErr != nil {
		return "", findCollectionErr
	}

	newManagement := models.NewRecord(managementCollection)
	newManagement.Set("user", userId)
	newManagement.Set("property", propertyId)

	if saveRecordErr := o.app.Dao().SaveRecord(newManagement); saveRecordErr != nil {
		return "", saveRecordErr
	}

	return newManagement.Get("id").(string), nil
}

func NewOwnershipRepository(app *pocketbase.PocketBase) IManagementRepository {
	return ManagementRepository{
		app,
	}
}
