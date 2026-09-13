package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(createAllCollections, dropAllCollections)
}

func createAllCollections(db dbx.Builder) error {
	dao := daos.New(db)

	usersCollection, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	schemaWithoutNameAndAvatar := schema.NewSchema()
	for _, usersField := range usersCollection.Schema.Fields() {
		if usersField.Name != "name" && usersField.Name != "avatar" {
			schemaWithoutNameAndAvatar.AddField(usersField)
		}
	}
	usersCollection.Schema = schemaWithoutNameAndAvatar
	if err := dao.SaveCollection(usersCollection); err != nil {
		return err
	}
	usersCollectionId := usersCollection.Id

	// Back-relations use _via_fieldName syntax in PocketBase v0.22 filter expressions.
	personsListViewRule := "user.id = @request.auth.id"

	personsCollection := &models.Collection{
		Name: "persons",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "pers1", Name: "name", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
			&schema.SchemaField{Id: "pers2", Name: "DOB", Type: schema.FieldTypeDate, Required: true, Options: &schema.DateOptions{}},
			&schema.SchemaField{Id: "pers3", Name: "user", Type: schema.FieldTypeRelation, Required: false,
				Options: &schema.RelationOptions{CollectionId: usersCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "pers4", Name: "profileImage", Type: schema.FieldTypeFile, Required: false,
				Options: &schema.FileOptions{MaxSelect: 1, MaxSize: 5242880}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_persons_user ON persons (user) WHERE user != ''",
		},
		ListRule:   types.Pointer(personsListViewRule),
		ViewRule:   types.Pointer(personsListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer("user.id = @request.auth.id"),
		DeleteRule: types.Pointer("user.id = @request.auth.id"),
	}
	if err := dao.SaveCollection(personsCollection); err != nil {
		return err
	}
	personsCollectionId := personsCollection.Id

	propertiesListViewRule := "" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || tenants_via_property.person.user.id = @request.auth.id" +
		" || households_via_property.person.user.id = @request.auth.id" +
		" || propertyAgents_via_property.agent.person.user.id = @request.auth.id"

	propertiesUpdateDeleteRule := "" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	propertiesCollection := &models.Collection{
		Name: "properties",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "prop1", Name: "address", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
		),
		ListRule:   types.Pointer(propertiesListViewRule),
		ViewRule:   types.Pointer(propertiesListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(propertiesUpdateDeleteRule),
		DeleteRule: types.Pointer(propertiesUpdateDeleteRule),
	}
	if err := dao.SaveCollection(propertiesCollection); err != nil {
		return err
	}
	propertiesCollectionId := propertiesCollection.Id

	propertyOwnersListViewRule := "" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	propertyOwnersUpdateDeleteRule := "" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	propertyOwnersCollection := &models.Collection{
		Name: "propertyOwners",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "pown1", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesCollectionId, MaxSelect: types.Pointer(1)}},
		),
		ListRule:   types.Pointer(propertyOwnersListViewRule),
		ViewRule:   types.Pointer(propertyOwnersListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(propertyOwnersUpdateDeleteRule),
		DeleteRule: types.Pointer(propertyOwnersUpdateDeleteRule),
	}
	if err := dao.SaveCollection(propertyOwnersCollection); err != nil {
		return err
	}
	propertyOwnersCollectionId := propertyOwnersCollection.Id

	cobrandsListViewRule := "" +
		"cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agents_via_cobrand.person.user.id = @request.auth.id"

	cobrandsUpdateDeleteRule := "cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandsCollection := &models.Collection{
		Name: "cobrands",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "cbrd1", Name: "name", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
		),
		ListRule:   types.Pointer(cobrandsListViewRule),
		ViewRule:   types.Pointer(cobrandsListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(cobrandsUpdateDeleteRule),
		DeleteRule: types.Pointer(cobrandsUpdateDeleteRule),
	}
	if err := dao.SaveCollection(cobrandsCollection); err != nil {
		return err
	}
	cobrandsCollectionId := cobrandsCollection.Id

	cobrandAdminsListViewRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.agents_via_cobrand.person.user.id = @request.auth.id"

	cobrandAdminsUpdateRule := "cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandAdminsDeleteRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || user.id = @request.auth.id"

	cobrandAdminsCollection := &models.Collection{
		Name: "cobrandAdmins",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "cbad1", Name: "user", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: usersCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "cbad2", Name: "cobrand", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: cobrandsCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_cobrandAdmins ON cobrandAdmins (user, cobrand)",
		},
		ListRule:   types.Pointer(cobrandAdminsListViewRule),
		ViewRule:   types.Pointer(cobrandAdminsListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(cobrandAdminsUpdateRule),
		DeleteRule: types.Pointer(cobrandAdminsDeleteRule),
	}
	if err := dao.SaveCollection(cobrandAdminsCollection); err != nil {
		return err
	}

	cobrandPropertyManagersListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandPropertyManagersCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandPropertyManagersDeleteRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandPropertyManagersCollection := &models.Collection{
		Name: "cobrandPropertyManagers",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "cbpm1", Name: "cobrand", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: cobrandsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "cbpm2", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_cobrandPropertyManagers ON cobrandPropertyManagers (cobrand, property)",
		},
		ListRule:   types.Pointer(cobrandPropertyManagersListViewRule),
		ViewRule:   types.Pointer(cobrandPropertyManagersListViewRule),
		CreateRule: types.Pointer(cobrandPropertyManagersCreateUpdateRule),
		UpdateRule: types.Pointer(cobrandPropertyManagersCreateUpdateRule),
		DeleteRule: types.Pointer(cobrandPropertyManagersDeleteRule),
	}
	if err := dao.SaveCollection(cobrandPropertyManagersCollection); err != nil {
		return err
	}

	cobrandPropertyOwnersListViewRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandPropertyOwnersCreateUpdateRule := "" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandPropertyOwnersDeleteRule := "" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandPropertyOwnersCollection := &models.Collection{
		Name: "cobrandPropertyOwners",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "cbpo1", Name: "cobrand", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: cobrandsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "cbpo2", Name: "propertyOwner", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertyOwnersCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_cobrandPropertyOwners ON cobrandPropertyOwners (cobrand, propertyOwner)",
		},
		ListRule:   types.Pointer(cobrandPropertyOwnersListViewRule),
		ViewRule:   types.Pointer(cobrandPropertyOwnersListViewRule),
		CreateRule: types.Pointer(cobrandPropertyOwnersCreateUpdateRule),
		UpdateRule: types.Pointer(cobrandPropertyOwnersCreateUpdateRule),
		DeleteRule: types.Pointer(cobrandPropertyOwnersDeleteRule),
	}
	if err := dao.SaveCollection(cobrandPropertyOwnersCollection); err != nil {
		return err
	}

	personPropertyOwnersListViewRule := "" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	personPropertyOwnersCreateUpdateDeleteRule := "" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	personPropertyOwnersCollection := &models.Collection{
		Name: "personPropertyOwners",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "ppow1", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "ppow2", Name: "propertyOwner", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertyOwnersCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_personPropertyOwners ON personPropertyOwners (person, propertyOwner)",
		},
		ListRule:   types.Pointer(personPropertyOwnersListViewRule),
		ViewRule:   types.Pointer(personPropertyOwnersListViewRule),
		CreateRule: types.Pointer(personPropertyOwnersCreateUpdateDeleteRule),
		UpdateRule: types.Pointer(personPropertyOwnersCreateUpdateDeleteRule),
		DeleteRule: types.Pointer(personPropertyOwnersCreateUpdateDeleteRule),
	}
	if err := dao.SaveCollection(personPropertyOwnersCollection); err != nil {
		return err
	}

	// PocketBase v0.22 can't handle 3+ back-relations in a chain, so property-based
	// access for agents is omitted here (it would require 3 back-relations).
	agentsListViewRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || person.user.id = @request.auth.id"

	agentsCreateUpdateRule := "cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	agentsDeleteRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || person.user.id = @request.auth.id"

	agentsCollection := &models.Collection{
		Name: "agents",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "agt01", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "agt02", Name: "cobrand", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: cobrandsCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_agents ON agents (person, cobrand)",
		},
		ListRule:   types.Pointer(agentsListViewRule),
		ViewRule:   types.Pointer(agentsListViewRule),
		CreateRule: types.Pointer(agentsCreateUpdateRule),
		UpdateRule: types.Pointer(agentsCreateUpdateRule),
		DeleteRule: types.Pointer(agentsDeleteRule),
	}
	if err := dao.SaveCollection(agentsCollection); err != nil {
		return err
	}
	agentsCollectionId := agentsCollection.Id

	propertyAgentsListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agent.person.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	propertyAgentsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	propertyAgentsDeleteRule := propertyAgentsCreateUpdateRule + " || agent.person.user.id = @request.auth.id"

	propertyAgentsCollection := &models.Collection{
		Name: "propertyAgents",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "pagt1", Name: "agent", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: agentsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "pagt2", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_propertyAgents ON propertyAgents (agent, property)",
		},
		ListRule:   types.Pointer(propertyAgentsListViewRule),
		ViewRule:   types.Pointer(propertyAgentsListViewRule),
		CreateRule: types.Pointer(propertyAgentsCreateUpdateRule),
		UpdateRule: types.Pointer(propertyAgentsCreateUpdateRule),
		DeleteRule: types.Pointer(propertyAgentsDeleteRule),
	}
	if err := dao.SaveCollection(propertyAgentsCollection); err != nil {
		return err
	}

	householdsListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id"

	householdsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	householdsDeleteRule := householdsCreateUpdateRule + " || person.user.id = @request.auth.id"

	householdsCollection := &models.Collection{
		Name: "households",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "hhld1", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "hhld2", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_households ON households (person, property)",
		},
		ListRule:   types.Pointer(householdsListViewRule),
		ViewRule:   types.Pointer(householdsListViewRule),
		CreateRule: types.Pointer(householdsCreateUpdateRule),
		UpdateRule: types.Pointer(householdsCreateUpdateRule),
		DeleteRule: types.Pointer(householdsDeleteRule),
	}
	if err := dao.SaveCollection(householdsCollection); err != nil {
		return err
	}

	tenantsListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id"

	tenantsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	tenantsDeleteRule := tenantsCreateUpdateRule + " || person.user.id = @request.auth.id"

	tenantsCollection := &models.Collection{
		Name: "tenants",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "tnt01", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "tnt02", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_tenants ON tenants (person, property)",
		},
		ListRule:   types.Pointer(tenantsListViewRule),
		ViewRule:   types.Pointer(tenantsListViewRule),
		CreateRule: types.Pointer(tenantsCreateUpdateRule),
		UpdateRule: types.Pointer(tenantsCreateUpdateRule),
		DeleteRule: types.Pointer(tenantsDeleteRule),
	}
	if err := dao.SaveCollection(tenantsCollection); err != nil {
		return err
	}

	// 3+ back-relation chains cause SQL errors in PocketBase v0.22, so property-owner
	// access is omitted (requires 3 backs via propertyItems→property→propertyOwners→...).
	// Residents (tenant/household/agent) are reachable via 2 back-relations.
	itemsListViewRule := "" +
		"personItems_via_item.person.user.id = @request.auth.id" +
		" || propertyItems_via_item.property.tenants_via_property.person.user.id = @request.auth.id" +
		" || propertyItems_via_item.property.households_via_property.person.user.id = @request.auth.id" +
		" || propertyItems_via_item.property.propertyAgents_via_property.agent.person.user.id = @request.auth.id"

	itemsUpdateDeleteRule := "personItems_via_item.person.user.id = @request.auth.id"

	itemsCollection := &models.Collection{
		Name: "items",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "itm01", Name: "name", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
			&schema.SchemaField{Id: "itm02", Name: "description", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
			&schema.SchemaField{Id: "itm03", Name: "picture", Type: schema.FieldTypeFile, Required: false,
				Options: &schema.FileOptions{MaxSelect: 1, MaxSize: 5242880}},
		),
		ListRule:   types.Pointer(itemsListViewRule),
		ViewRule:   types.Pointer(itemsListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(itemsUpdateDeleteRule),
		DeleteRule: types.Pointer(itemsUpdateDeleteRule),
	}
	if err := dao.SaveCollection(itemsCollection); err != nil {
		return err
	}
	itemsCollectionId := itemsCollection.Id

	// PocketBase v0.22: forward→back→forward→back chains from entryDevices fail.
	// Only the forward→back chain (item → personItems) is supported.
	entryDevicesListViewRule := "item.personItems_via_item.person.user.id = @request.auth.id"
	entryDevicesCreateUpdateDeleteRule := "item.personItems_via_item.person.user.id = @request.auth.id"

	entryDevicesCollection := &models.Collection{
		Name: "entryDevices",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "edv01", Name: "item", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: itemsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "edv02", Name: "deviceType", Type: schema.FieldTypeSelect, Required: true,
				Options: &schema.SelectOptions{MaxSelect: 1, Values: []string{"Fob", "Key", "Remote", "RoomKey", "MailboxKey"}}},
			&schema.SchemaField{Id: "edv03", Name: "identifier", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
			&schema.SchemaField{Id: "edv04", Name: "defunctReason", Type: schema.FieldTypeSelect, Required: true,
				Options: &schema.SelectOptions{MaxSelect: 1, Values: []string{"None", "Lost", "Damaged", "Retired", "Stolen"}}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_entryDevices_item ON entryDevices (item)",
		},
		ListRule:   types.Pointer(entryDevicesListViewRule),
		ViewRule:   types.Pointer(entryDevicesListViewRule),
		CreateRule: types.Pointer(entryDevicesCreateUpdateDeleteRule),
		UpdateRule: types.Pointer(entryDevicesCreateUpdateDeleteRule),
		DeleteRule: types.Pointer(entryDevicesCreateUpdateDeleteRule),
	}
	if err := dao.SaveCollection(entryDevicesCollection); err != nil {
		return err
	}

	propertyItemsListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.propertyAgents_via_property.agent.person.user.id = @request.auth.id"

	propertyItemsCreateUpdateDeleteRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id"

	propertyItemsCollection := &models.Collection{
		Name: "propertyItems",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "pitm1", Name: "item", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: itemsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "pitm2", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_propertyItems ON propertyItems (item, property)",
		},
		ListRule:   types.Pointer(propertyItemsListViewRule),
		ViewRule:   types.Pointer(propertyItemsListViewRule),
		CreateRule: types.Pointer(propertyItemsCreateUpdateDeleteRule),
		UpdateRule: types.Pointer(propertyItemsCreateUpdateDeleteRule),
		DeleteRule: types.Pointer(propertyItemsCreateUpdateDeleteRule),
	}
	if err := dao.SaveCollection(propertyItemsCollection); err != nil {
		return err
	}

	// Property-owner access requires 3+ back-relations; omitted due to PocketBase v0.22 limitation.
	personItemsListViewRule := "person.user.id = @request.auth.id"
	personItemsUpdateDeleteRule := "person.user.id = @request.auth.id"

	personItemsCollection := &models.Collection{
		Name: "personItems",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "psit1", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsCollectionId, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "psit2", Name: "item", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: itemsCollectionId, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_personItems ON personItems (person, item)",
		},
		ListRule:   types.Pointer(personItemsListViewRule),
		ViewRule:   types.Pointer(personItemsListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(personItemsUpdateDeleteRule),
		DeleteRule: types.Pointer(personItemsUpdateDeleteRule),
	}
	if err := dao.SaveCollection(personItemsCollection); err != nil {
		return err
	}

	return nil
}

func dropAllCollections(db dbx.Builder) error {
	dao := daos.New(db)

	for _, name := range []string{
		"personItems", "propertyItems", "entryDevices", "items",
		"tenants", "households", "propertyAgents", "agents",
		"personPropertyOwners", "cobrandPropertyOwners",
		"cobrandPropertyManagers", "cobrandAdmins", "cobrands",
		"propertyOwners", "properties", "persons",
	} {
		collection, err := dao.FindCollectionByNameOrId(name)
		if err != nil {
			continue
		}
		if err := dao.DeleteCollection(collection); err != nil {
			return err
		}
	}

	// Restore users name and avatar fields.
	usersCollection, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	usersCollection.Schema.AddField(&schema.SchemaField{
		Id:      "users_name",
		Type:    schema.FieldTypeText,
		Name:    "name",
		Options: &schema.TextOptions{},
	})
	usersCollection.Schema.AddField(&schema.SchemaField{
		Id:   "users_avatar",
		Type: schema.FieldTypeFile,
		Name: "avatar",
		Options: &schema.FileOptions{
			MaxSelect: 1,
			MaxSize:   5242880,
			MimeTypes: []string{"image/jpeg", "image/png", "image/svg+xml", "image/gif", "image/webp"},
		},
	})
	return dao.SaveCollection(usersCollection)
}
