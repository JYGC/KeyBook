package main

import (
	"errors"
	_ "keybook/backend/migrations"
	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"go.uber.org/dig"
)

func startBackend() {
	container := dig.New()
	container.Provide(pocketbase.New)
	container.Provide(func(app *pocketbase.PocketBase) core.App { return app })

	// Repositories
	container.Provide(repositories.NewPersonRepository)
	container.Provide(repositories.NewPropertyRepository)
	container.Provide(repositories.NewItemRepository)
	container.Provide(repositories.NewPropertyItemRepository)
	container.Provide(repositories.NewPersonItemRepository)
	container.Provide(repositories.NewEntryDeviceRepository)
	container.Provide(repositories.NewCobrandRepository)
	container.Provide(repositories.NewCobrandAdminRepository)
	container.Provide(repositories.NewCobrandPropertyManagerRepository)
	container.Provide(repositories.NewCobrandPropertyOwnerRepository)
	container.Provide(repositories.NewPropertyOwnerRepository)
	container.Provide(repositories.NewPersonPropertyOwnerRepository)
	container.Provide(repositories.NewAgentRepository)
	container.Provide(repositories.NewPropertyAgentRepository)
	container.Provide(repositories.NewHouseholdRepository)
	container.Provide(repositories.NewTenantRepository)

	// Services
	container.Provide(services.NewPersonService)
	container.Provide(services.NewPropertyService)
	container.Provide(services.NewItemService)
	container.Provide(services.NewCobrandService)
	container.Provide(services.NewAgentService)
	container.Provide(services.NewPropertyOwnerService)

	// Application services
	container.Provide(application.NewPersonApplicationService)
	container.Provide(application.NewPropertyApplicationService)
	container.Provide(application.NewItemApplicationService)
	container.Provide(application.NewCobrandApplicationService)
	container.Provide(application.NewAgentApplicationService)
	container.Provide(application.NewPropertyOwnerApplicationService)

	invokeErr := container.Invoke(func(
		app *pocketbase.PocketBase,
		personSvc services.IPersonService,
		propertySvc services.IPropertyService,
		itemSvc services.IItemService,
		cobrandSvc services.ICobrandService,
		agentSvc services.IAgentService,
		propertyOwnerSvc services.IPropertyOwnerService,
		entryDeviceRepo repositories.IEntryDeviceRepository,
		cobrandAdminRepo repositories.ICobrandAdminRepository,
		personPropertyOwnerRepo repositories.IPersonPropertyOwnerRepository,
		cobrandPropertyOwnerRepo repositories.ICobrandPropertyOwnerRepository,
		propertyAgentRepo repositories.IPropertyAgentRepository,
		propertyOwnerRepo repositories.IPropertyOwnerRepository,
	) {
		app.OnRecordBeforeCreateRequest("persons").Add(func(e *core.RecordCreateEvent) error {
			return personSvc.ValidatePerson(e.Record.GetString("name"), e.Record.GetString("DOB"))
		})

		app.OnRecordBeforeCreateRequest("properties").Add(func(e *core.RecordCreateEvent) error {
			return propertySvc.ValidateProperty(e.Record.GetString("address"))
		})

		app.OnRecordBeforeCreateRequest("items").Add(func(e *core.RecordCreateEvent) error {
			return itemSvc.ValidateItem(e.Record.GetString("name"))
		})

		app.OnRecordBeforeUpdateRequest("entryDevices").Add(func(e *core.RecordUpdateEvent) error {
			current, err := entryDeviceRepo.GetEntryDeviceById(e.Record.GetId())
			if err != nil {
				return err
			}
			return itemSvc.ValidateEntryDeviceTransition(current.DefunctReason, e.Record.GetString("defunctReason"))
		})

		app.OnRecordBeforeCreateRequest("cobrands").Add(func(e *core.RecordCreateEvent) error {
			return cobrandSvc.ValidateCobrand(e.Record.GetString("name"))
		})

		app.OnRecordBeforeCreateRequest("cobrandAdmins").Add(func(e *core.RecordCreateEvent) error {
			existing, err := cobrandAdminRepo.GetCobrandAdminsByCobrandId(e.Record.GetString("cobrand"))
			if err != nil {
				return err
			}
			return cobrandSvc.EnsureAdminIsUnique(existing, e.Record.GetString("user"))
		})

		app.OnRecordBeforeCreateRequest("personPropertyOwners").Add(func(e *core.RecordCreateEvent) error {
			existing, err := personPropertyOwnerRepo.GetPersonPropertyOwnersByPropertyOwnerId(e.Record.GetString("propertyOwner"))
			if err != nil {
				return err
			}
			return propertyOwnerSvc.EnsureNoDuplicatePersonOwner(existing, e.Record.GetString("person"))
		})

		app.OnRecordBeforeCreateRequest("cobrandPropertyOwners").Add(func(e *core.RecordCreateEvent) error {
			existing, err := cobrandPropertyOwnerRepo.GetCobrandPropertyOwnersByPropertyOwnerId(e.Record.GetString("propertyOwner"))
			if err != nil {
				return err
			}
			return propertyOwnerSvc.EnsureNoDuplicateCobrandOwner(existing, e.Record.GetString("cobrand"))
		})

		app.OnRecordBeforeCreateRequest("propertyAgents").Add(func(e *core.RecordCreateEvent) error {
			existing, err := propertyAgentRepo.GetPropertyAgentsByPropertyId(e.Record.GetString("property"))
			if err != nil {
				return err
			}
			return agentSvc.EnsureNoDuplicatePropertyAgent(existing, e.Record.GetString("agent"))
		})

		app.OnRecordBeforeDeleteRequest("propertyOwners").Add(func(e *core.RecordDeleteEvent) error {
			owners, err := propertyOwnerRepo.GetPropertyOwnersByPropertyId(e.Record.GetString("property"))
			if err != nil {
				return err
			}
			if len(owners) <= 1 {
				return errors.New("cannot delete the last property owner")
			}
			return nil
		})

		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	})
	if invokeErr != nil {
		log.Fatal(invokeErr)
	}
}

func main() {
	startBackend()
}
