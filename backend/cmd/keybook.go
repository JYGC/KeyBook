package main

import (
	"fmt"
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

	var constructors []interface{}

	constructors = append(constructors, pocketbase.New)
	constructors = append(constructors, databases.NewHistoryDatabase)

	constructors = append(constructors, repositories.NewDeviceRepository)
	constructors = append(constructors, repositories.NewPersonRepository)
	constructors = append(constructors, repositories.NewPropertyRepository)
	constructors = append(constructors, repositories.NewPersonDeviceRepository)
	constructors = append(constructors, repositories.NewHistoryRepository)

	constructors = append(constructors, services.NewDataImportServices)
	constructors = append(constructors, services.NewDeviceServices)
	constructors = append(constructors, services.NewPersonDeviceServices)
	constructors = append(constructors, services.NewPersonServices)
	constructors = append(constructors, services.NewPropertyServices)

	constructors = append(constructors, handlers.NewDeviceHandlers)
	constructors = append(constructors, handlers.NewHistoryHandlers)

	for _, constructor := range constructors {
		provideConstructorErr := container.Provide(constructor)
		if provideConstructorErr != nil {
			fmt.Printf("provideConstructorErr: %v\n", provideConstructorErr)
			return
		}
	}

	invokeErr := container.Invoke(func(
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
	if invokeErr != nil {
		fmt.Printf("invokeErr: %v\n", invokeErr)
		return
	}
}

func main() {
	go startFrontend()
	startBackend()
}
