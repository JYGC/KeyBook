package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(up1749240000, down1749240000)
}

// setCascadeDelete sets CascadeDelete on the named relation field in each collection.
func setCascadeDelete(db dbx.Builder, value bool) error {
	dao := daos.New(db)

	targets := []struct {
		collection string
		fieldId    string
	}{
		{"entryDevices", "edv01"},
		{"personItems", "psit2"},
		{"propertyItems", "pitm1"},
	}

	for _, t := range targets {
		col, err := dao.FindCollectionByNameOrId(t.collection)
		if err != nil {
			return err
		}
		field := col.Schema.GetFieldById(t.fieldId)
		if field == nil {
			continue
		}
		opts, ok := field.Options.(*schema.RelationOptions)
		if !ok {
			continue
		}
		opts.CascadeDelete = value
		if err := dao.SaveCollection(col); err != nil {
			return err
		}
	}
	return nil
}

// up1749240000 enables cascade delete on item-relation fields in entryDevices,
// personItems, and propertyItems. Deleting an item now automatically removes all
// dependent records. The item's deleteRule is evaluated first (while dependents
// still exist), then PocketBase cascades the deletions.
func up1749240000(db dbx.Builder) error {
	return setCascadeDelete(db, true)
}

func down1749240000(db dbx.Builder) error {
	return setCascadeDelete(db, false)
}
