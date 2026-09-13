package main

import (
	"errors"
	"keybook/backend/internal/application"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
	_ "keybook/backend/migrations"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"go.uber.org/dig"
)

func startBackend() {
	container := dig.New()
	container.Provide(pocketbase.New)
	container.Provide(func(app *pocketbase.PocketBase) core.App { return app })

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

	container.Provide(services.NewPersonService)
	container.Provide(services.NewPropertyService)
	container.Provide(services.NewItemService)
	container.Provide(services.NewCobrandService)
	container.Provide(services.NewAgentService)
	container.Provide(services.NewPropertyOwnerService)

	container.Provide(application.NewPersonApplicationService)
	container.Provide(application.NewPropertyApplicationService)
	container.Provide(application.NewItemApplicationService)
	container.Provide(application.NewCobrandApplicationService)
	container.Provide(application.NewAgentApplicationService)
	container.Provide(application.NewPropertyOwnerApplicationService)

	invokeErr := container.Invoke(func(
		app *pocketbase.PocketBase,
		personService services.IPersonService,
		propertyService services.IPropertyService,
		itemService services.IItemService,
		cobrandService services.ICobrandService,
		agentService services.IAgentService,
		propertyOwnerService services.IPropertyOwnerService,
		entryDeviceRepository repositories.IEntryDeviceRepository,
		cobrandAdminRepository repositories.ICobrandAdminRepository,
		personPropertyOwnerRepository repositories.IPersonPropertyOwnerRepository,
		cobrandPropertyOwnerRepository repositories.ICobrandPropertyOwnerRepository,
		propertyAgentRepository repositories.IPropertyAgentRepository,
		propertyOwnerRepository repositories.IPropertyOwnerRepository,
	) {
		app.OnRecordBeforeCreateRequest("persons").Add(func(personCreateEvent *core.RecordCreateEvent) error {
			return personService.ValidatePerson(
				personCreateEvent.Record.GetString("name"),
				personCreateEvent.Record.GetString("DOB"),
			)
		})

		app.OnRecordBeforeCreateRequest("properties").Add(func(propertyCreateEvent *core.RecordCreateEvent) error {
			return propertyService.ValidateProperty(propertyCreateEvent.Record.GetString("address"))
		})

		app.OnRecordBeforeCreateRequest("items").Add(func(itemCreateEvent *core.RecordCreateEvent) error {
			return itemService.ValidateItem(itemCreateEvent.Record.GetString("name"))
		})

		app.OnRecordBeforeUpdateRequest("entryDevices").Add(func(entryDeviceUpdateEvent *core.RecordUpdateEvent) error {
			currentEntryDevice, err := entryDeviceRepository.GetEntryDeviceById(entryDeviceUpdateEvent.Record.GetId())
			if err != nil {
				return err
			}
			return itemService.ValidateEntryDeviceTransition(
				currentEntryDevice.DefunctReason,
				entryDeviceUpdateEvent.Record.GetString("defunctReason"),
			)
		})

		app.OnRecordBeforeCreateRequest("cobrands").Add(func(cobrandCreateEvent *core.RecordCreateEvent) error {
			return cobrandService.ValidateCobrand(cobrandCreateEvent.Record.GetString("name"))
		})

		app.OnRecordBeforeCreateRequest("cobrandAdmins").Add(func(cobrandAdminCreateEvent *core.RecordCreateEvent) error {
			existingAdmins, err := cobrandAdminRepository.GetCobrandAdminsByCobrandId(
				cobrandAdminCreateEvent.Record.GetString("cobrand"),
			)
			if err != nil {
				return err
			}
			if err := cobrandService.EnsureAdminIsUnique(
				existingAdmins,
				cobrandAdminCreateEvent.Record.GetString("user"),
			); err != nil {
				return err
			}

			// PocketBase superusers bypass the inviter-approval check below —
			// they have no cobrandAdmins record of their own to check.
			requestInfo := apis.RequestInfo(cobrandAdminCreateEvent.HttpContext)
			if requestInfo.Admin != nil {
				return nil
			}
			inviterUserId := ""
			if requestInfo.AuthRecord != nil {
				inviterUserId = requestInfo.AuthRecord.Id
			}
			return cobrandService.EnsureInviterIsApprovedAdmin(existingAdmins, inviterUserId)
		})

		app.OnRecordBeforeCreateRequest("personPropertyOwners").Add(func(personPropertyOwnerCreateEvent *core.RecordCreateEvent) error {
			existingPersonOwners, err := personPropertyOwnerRepository.GetPersonPropertyOwnersByPropertyOwnerId(
				personPropertyOwnerCreateEvent.Record.GetString("propertyOwner"),
			)
			if err != nil {
				return err
			}
			return propertyOwnerService.EnsureNoDuplicatePersonOwner(
				existingPersonOwners,
				personPropertyOwnerCreateEvent.Record.GetString("person"),
			)
		})

		app.OnRecordBeforeCreateRequest("cobrandPropertyOwners").Add(func(cobrandPropertyOwnerCreateEvent *core.RecordCreateEvent) error {
			existingCobrandOwners, err := cobrandPropertyOwnerRepository.GetCobrandPropertyOwnersByPropertyOwnerId(
				cobrandPropertyOwnerCreateEvent.Record.GetString("propertyOwner"),
			)
			if err != nil {
				return err
			}
			return propertyOwnerService.EnsureNoDuplicateCobrandOwner(
				existingCobrandOwners,
				cobrandPropertyOwnerCreateEvent.Record.GetString("cobrand"),
			)
		})

		app.OnRecordBeforeCreateRequest("propertyAgents").Add(func(propertyAgentCreateEvent *core.RecordCreateEvent) error {
			existingPropertyAgents, err := propertyAgentRepository.GetPropertyAgentsByPropertyId(
				propertyAgentCreateEvent.Record.GetString("property"),
			)
			if err != nil {
				return err
			}
			return agentService.EnsureNoDuplicatePropertyAgent(
				existingPropertyAgents,
				propertyAgentCreateEvent.Record.GetString("agent"),
			)
		})

		app.OnRecordBeforeDeleteRequest("propertyOwners").Add(func(propertyOwnerDeleteEvent *core.RecordDeleteEvent) error {
			ownersOfSameProperty, err := propertyOwnerRepository.GetPropertyOwnersByPropertyId(
				propertyOwnerDeleteEvent.Record.GetString("property"),
			)
			if err != nil {
				return err
			}
			if len(ownersOfSameProperty) <= 1 {
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
