package repositories_test

import (
	"testing"

	_ "keybook/backend/migrations"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

// newApp creates a fresh PocketBase test instance with migrations applied.
func newApp(t *testing.T) *tests.TestApp {
	t.Helper()
	app, err := tests.NewTestApp(t.TempDir())
	if err != nil {
		t.Fatalf("new test app: %v", err)
	}
	t.Cleanup(app.Cleanup)
	return app
}

// createRecord inserts a record directly via the DAO (bypasses access rules).
func createRecord(t *testing.T, app *tests.TestApp, collection string, fields map[string]any) *models.Record {
	t.Helper()
	col, err := app.Dao().FindCollectionByNameOrId(collection)
	if err != nil {
		t.Fatalf("find collection %q: %v", collection, err)
	}
	r := models.NewRecord(col)
	for k, v := range fields {
		r.Set(k, v)
	}
	if err := app.Dao().SaveRecord(r); err != nil {
		t.Fatalf("save %q record: %v", collection, err)
	}
	return r
}
