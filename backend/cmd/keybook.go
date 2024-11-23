package main

import (
	"fmt"
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

	constructors = append(constructors, repositories.NewDeviceRepository)
	constructors = append(constructors, repositories.NewPersonRepository)
	constructors = append(constructors, repositories.NewPropertyRepository)
	constructors = append(constructors, repositories.NewPersonDeviceRepository)
	constructors = append(constructors, repositories.NewDeviceHistoryRepository)
	constructors = append(constructors, repositories.NewPersonHistoryRepository)

	constructors = append(constructors, services.NewDataImportServices)
	constructors = append(constructors, services.NewDeviceServices)
	constructors = append(constructors, services.NewPersonDeviceServices)
	constructors = append(constructors, services.NewPersonServices)
	constructors = append(constructors, services.NewPropertyServices)
	constructors = append(constructors, services.NewDeviceHistoryServices)
	constructors = append(constructors, services.NewPersonHistoryServices)

	constructors = append(constructors, handlers.NewDeviceHandlers)
	constructors = append(constructors, handlers.NewDataImportHandlers)

	for _, constructor := range constructors {
		provideConstructorErr := container.Provide(constructor)
		if provideConstructorErr != nil {
			fmt.Printf("provideConstructorErr: %v\n", provideConstructorErr)
			return
		}
	}

	invokeErr := container.Invoke(func(
		app *pocketbase.PocketBase,
		deviceHistoryServices services.IDeviceHistoryServices,
		personHistoryServices services.IPersonHistoryServices,
		deviceHandlers handlers.IDeviceHoldingHandlers,
		dataImportHandlers handlers.IDataImportHandlers,
	) {
		app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
			handlers.RegisterDeviceHandlersToRouter(e.Router, deviceHandlers)
			handlers.RegisterDataImportHandlersToRouter(e.Router, dataImportHandlers)
			return nil
		})

		app.OnModelAfterCreate("devices").Add(func(e *core.ModelEvent) error {
			if updateErr := deviceHistoryServices.AddNewDeviceHistoryDueToCreateDeviceHook(e.Model); updateErr != nil {
				return updateErr
			}
			return nil
		})

		app.OnModelBeforeUpdate("devices").Add(func(e *core.ModelEvent) error {
			if updateErr := deviceHistoryServices.AddNewDeviceHistoryDueToUpdateDeviceHook(e.Model); updateErr != nil {
				return updateErr
			}
			return nil
		})

		app.OnModelAfterCreate("persons").Add(func(e *core.ModelEvent) error {
			if updateErr := personHistoryServices.AddNewPersonHistoryDueToCreatePersonHook(e.Model); updateErr != nil {
				return updateErr
			}
			return nil
		})

		app.OnModelBeforeUpdate("persons").Add(func(e *core.ModelEvent) error {
			if updateErr := personHistoryServices.AddNewPersonHistoryDueToUpdatePersonHook(e.Model); updateErr != nil {
				return updateErr
			}
			return nil
		})

		app.OnModelAfterCreate("persondevices").Add(func(e *core.ModelEvent) error {
			return nil
		})

		app.OnModelBeforeUpdate("persondevices").Add(func(e *core.ModelEvent) error {
			return nil
		})

		app.OnModelAfterCreate("properties").Add(func(e *core.ModelEvent) error {
			return nil
		})

		app.OnModelBeforeUpdate("properties").Add(func(e *core.ModelEvent) error {
			return nil
		})

		app.OnModelAfterCreate("users").Add(func(e *core.ModelEvent) error {
			return nil
		})

		app.OnModelBeforeUpdate("users").Add(func(e *core.ModelEvent) error {
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
