package main

import (
	"fmt"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"go.uber.org/dig"
)

func startBackend() {
	container := dig.New()

	var constructors []interface{}

	constructors = append(constructors, pocketbase.New)

	constructors = append(constructors, repositories.NewDeviceRepository)
	constructors = append(constructors, repositories.NewPersonRepository)
	constructors = append(constructors, repositories.NewPropertyRepository)
	constructors = append(constructors, repositories.NewPropertyHistoryRepository)
	constructors = append(constructors, repositories.NewDeviceHistoryRepository)
	constructors = append(constructors, repositories.NewPersonHistoryRepository)
	constructors = append(constructors, repositories.NewPersonDeviceHistoryRepository)

	constructors = append(constructors, services.NewPropertyHistoryServices)
	constructors = append(constructors, services.NewDeviceHistoryServices)
	constructors = append(constructors, services.NewPersonHistoryServices)
	constructors = append(constructors, services.NewPersonDeviceHistoryServices)

	for _, constructor := range constructors {
		provideConstructorErr := container.Provide(constructor)
		if provideConstructorErr != nil {
			fmt.Printf("provideConstructorErr: %v\n", provideConstructorErr)
			return
		}
	}

	invokeErr := container.Invoke(func(
		app *pocketbase.PocketBase,
		propertyHistoryServices services.IPropertyHistoryServices,
		deviceHistoryServices services.IDeviceHistoryServices,
		personHistoryServices services.IPersonHistoryServices,
		personDeviceHistoryServices services.IPersonDeviceHistoryServices,
	) {
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
			if updateErr := personDeviceHistoryServices.AddNewPersonDeviceHistoryDueToCreatePersonDeviceHook(e.Model); updateErr != nil {
				return updateErr
			}
			return nil
		})

		app.OnModelBeforeUpdate("persondevices").Add(func(e *core.ModelEvent) error {
			if updateErr := personDeviceHistoryServices.AddNewPersonDeviceHistoryDueToUpdatePersonDeviceHook(e.Model); updateErr != nil {
				return updateErr
			}
			return nil
		})

		app.OnModelAfterCreate("properties").Add(func(e *core.ModelEvent) error {
			return propertyHistoryServices.AddPropertyHistoryDueToCreatePropertyHook(e.Model)
		})

		app.OnModelBeforeUpdate("properties").Add(func(e *core.ModelEvent) error {
			return propertyHistoryServices.AddPropertyHistoryDueToUpdatePropertyHook(e.Model)
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
	startBackend()
}
