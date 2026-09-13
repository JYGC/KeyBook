package application_test

import (
	"testing"

	_ "keybook/backend/migrations"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func newTestAppWithMigrations(t *testing.T) *tests.TestApp {
	t.Helper()
	app, err := tests.NewTestApp(t.TempDir())
	if err != nil {
		t.Fatalf("new test app: %v", err)
	}
	t.Cleanup(app.Cleanup)
	return app
}

func createRecordBypassingAccessRules(t *testing.T, app *tests.TestApp, collection string, fields map[string]any) *models.Record {
	t.Helper()
	targetCollection, err := app.Dao().FindCollectionByNameOrId(collection)
	if err != nil {
		t.Fatalf("find collection %q: %v", collection, err)
	}
	record := models.NewRecord(targetCollection)
	for fieldName, fieldValue := range fields {
		record.Set(fieldName, fieldValue)
	}
	if err := app.Dao().SaveRecord(record); err != nil {
		t.Fatalf("save %q record: %v", collection, err)
	}
	return record
}
