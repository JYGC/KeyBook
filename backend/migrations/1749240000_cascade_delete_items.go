package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(enableCascadeDeleteOnItemRelations, disableCascadeDeleteOnItemRelations)
}

func setCascadeDeleteOnItemRelationFields(db dbx.Builder, cascadeDelete bool) error {
	dao := daos.New(db)

	itemRelationFields := []struct {
		collectionName    string
		itemRelationField string
	}{
		{"entryDevices", "edv01"},
		{"personItems", "psit2"},
		{"propertyItems", "pitm1"},
	}

	for _, itemRelation := range itemRelationFields {
		collection, err := dao.FindCollectionByNameOrId(itemRelation.collectionName)
		if err != nil {
			return err
		}
		field := collection.Schema.GetFieldById(itemRelation.itemRelationField)
		if field == nil {
			continue
		}
		relationOptions, fieldIsRelation := field.Options.(*schema.RelationOptions)
		if !fieldIsRelation {
			continue
		}
		relationOptions.CascadeDelete = cascadeDelete
		if err := dao.SaveCollection(collection); err != nil {
			return err
		}
	}
	return nil
}

// The item's own deleteRule is evaluated first, while its dependents still
// exist; only then does PocketBase cascade the deletions.
func enableCascadeDeleteOnItemRelations(db dbx.Builder) error {
	return setCascadeDeleteOnItemRelationFields(db, true)
}

func disableCascadeDeleteOnItemRelations(db dbx.Builder) error {
	return setCascadeDeleteOnItemRelationFields(db, false)
}
