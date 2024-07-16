package main

import (
	"keybook/backend/internal/databases"
	"keybook/backend/internal/frontend"
	"keybook/backend/internal/handlers"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"go.uber.org/dig"
)

func startFrontend() {
	mux := http.NewServeMux()
	mux.Handle("/", frontend.SvelteKitHandler("/"))
	log.Fatal(http.ListenAndServe(":5050", mux))
}

func startBackend() {
	container := dig.New()
	container.Provide(pocketbase.New)
	container.Provide(databases.NewHistoryDatabase)

	container.Provide(repositories.NewPersonRepository)
	container.Provide(repositories.NewPropertyRepository)
	container.Provide(repositories.NewOwnershipRepository)
	container.Provide(repositories.NewPersonDeviceRepository)
	container.Provide(repositories.NewHistoryRepository)

	container.Provide(services.NewDataImportServices)
	container.Provide(services.NewPropertyServices)

	container.Provide(handlers.NewDeviceHandlers)
	container.Provide(handlers.NewHistoryHandlers)

	container.Invoke(func(
		app *pocketbase.PocketBase,
		personRepository repositories.IPersonRepository,
		deviceHandlers handlers.IDeviceHandlers,
		historyHandlers handlers.IHistoryHandlers,
	) {
		app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
			handlers.RegisterDeviceHandlersToRouter(e.Router, deviceHandlers)
			handlers.RegisterHistoryHandlersToRouter(e.Router, historyHandlers)
			return nil
		})

		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	})
}

func main() {
	go startFrontend()
	startBackend()
}
