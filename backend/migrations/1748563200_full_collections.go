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
	m.Register(up1748563200, down1748563200)
}

func up1748563200(db dbx.Builder) error {
	dao := daos.New(db)

	// 1. Modify users: remove default name and avatar fields.
	usersCol, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	filteredSchema := schema.NewSchema()
	for _, f := range usersCol.Schema.Fields() {
		if f.Name != "name" && f.Name != "avatar" {
			filteredSchema.AddField(f)
		}
	}
	usersCol.Schema = filteredSchema
	if err := dao.SaveCollection(usersCol); err != nil {
		return err
	}
	usersColID := usersCol.Id

	// 2. persons
	// Back-relations use _via_fieldName syntax in PocketBase v0.22 filter expressions.
	personsListViewRule := "user.id = @request.auth.id"

	personsCol := &models.Collection{
		Name: "persons",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "pers1", Name: "name", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
			&schema.SchemaField{Id: "pers2", Name: "DOB", Type: schema.FieldTypeDate, Required: true, Options: &schema.DateOptions{}},
			&schema.SchemaField{Id: "pers3", Name: "user", Type: schema.FieldTypeRelation, Required: false,
				Options: &schema.RelationOptions{CollectionId: usersColID, MaxSelect: types.Pointer(1)}},
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
	if err := dao.SaveCollection(personsCol); err != nil {
		return err
	}
	personsColID := personsCol.Id

	// 3. properties
	propListViewRule := "" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || tenants_via_property.person.user.id = @request.auth.id" +
		" || households_via_property.person.user.id = @request.auth.id" +
		" || propertyAgents_via_property.agent.person.user.id = @request.auth.id"

	propUpdateDeleteRule := "" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	propertiesCol := &models.Collection{
		Name: "properties",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "prop1", Name: "address", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
		),
		ListRule:   types.Pointer(propListViewRule),
		ViewRule:   types.Pointer(propListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(propUpdateDeleteRule),
		DeleteRule: types.Pointer(propUpdateDeleteRule),
	}
	if err := dao.SaveCollection(propertiesCol); err != nil {
		return err
	}
	propertiesColID := propertiesCol.Id

	// 4. propertyOwners
	poListViewRule := "" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	poUpdateDeleteRule := "" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	propertyOwnersCol := &models.Collection{
		Name: "propertyOwners",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "pown1", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesColID, MaxSelect: types.Pointer(1)}},
		),
		ListRule:   types.Pointer(poListViewRule),
		ViewRule:   types.Pointer(poListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(poUpdateDeleteRule),
		DeleteRule: types.Pointer(poUpdateDeleteRule),
	}
	if err := dao.SaveCollection(propertyOwnersCol); err != nil {
		return err
	}
	propertyOwnersColID := propertyOwnersCol.Id

	// 5. cobrands
	cobrandListViewRule := "" +
		"cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agents_via_cobrand.person.user.id = @request.auth.id"

	cobrandUpdateDeleteRule := "cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandsCol := &models.Collection{
		Name: "cobrands",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "cbrd1", Name: "name", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
		),
		ListRule:   types.Pointer(cobrandListViewRule),
		ViewRule:   types.Pointer(cobrandListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(cobrandUpdateDeleteRule),
		DeleteRule: types.Pointer(cobrandUpdateDeleteRule),
	}
	if err := dao.SaveCollection(cobrandsCol); err != nil {
		return err
	}
	cobrandsColID := cobrandsCol.Id

	// 6. cobrandAdmins  (unique on user + cobrand)
	caListViewRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.agents_via_cobrand.person.user.id = @request.auth.id"

	caUpdateRule := "cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	caDeleteRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || user.id = @request.auth.id"

	cobrandAdminsCol := &models.Collection{
		Name: "cobrandAdmins",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "cbad1", Name: "user", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: usersColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "cbad2", Name: "cobrand", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: cobrandsColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_cobrandAdmins ON cobrandAdmins (user, cobrand)",
		},
		ListRule:   types.Pointer(caListViewRule),
		ViewRule:   types.Pointer(caListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(caUpdateRule),
		DeleteRule: types.Pointer(caDeleteRule),
	}
	if err := dao.SaveCollection(cobrandAdminsCol); err != nil {
		return err
	}

	// 7. cobrandPropertyManagers  (unique on cobrand + property)
	cpmListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cpmCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cpmDeleteRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandPropMgrsCol := &models.Collection{
		Name: "cobrandPropertyManagers",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "cbpm1", Name: "cobrand", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: cobrandsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "cbpm2", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_cobrandPropertyManagers ON cobrandPropertyManagers (cobrand, property)",
		},
		ListRule:   types.Pointer(cpmListViewRule),
		ViewRule:   types.Pointer(cpmListViewRule),
		CreateRule: types.Pointer(cpmCreateUpdateRule),
		UpdateRule: types.Pointer(cpmCreateUpdateRule),
		DeleteRule: types.Pointer(cpmDeleteRule),
	}
	if err := dao.SaveCollection(cobrandPropMgrsCol); err != nil {
		return err
	}

	// 8. cobrandPropertyOwners  (unique on cobrand + propertyOwner)
	cpoListViewRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cpoCreateUpdateRule := "" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cpoDeleteRule := "" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	cobrandPropOwnersCol := &models.Collection{
		Name: "cobrandPropertyOwners",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "cbpo1", Name: "cobrand", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: cobrandsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "cbpo2", Name: "propertyOwner", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertyOwnersColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_cobrandPropertyOwners ON cobrandPropertyOwners (cobrand, propertyOwner)",
		},
		ListRule:   types.Pointer(cpoListViewRule),
		ViewRule:   types.Pointer(cpoListViewRule),
		CreateRule: types.Pointer(cpoCreateUpdateRule),
		UpdateRule: types.Pointer(cpoCreateUpdateRule),
		DeleteRule: types.Pointer(cpoDeleteRule),
	}
	if err := dao.SaveCollection(cobrandPropOwnersCol); err != nil {
		return err
	}

	// 9. personPropertyOwners  (unique on person + propertyOwner)
	ppoListViewRule := "" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	ppoCUDRule := "" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	personPropOwnersCol := &models.Collection{
		Name: "personPropertyOwners",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "ppow1", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "ppow2", Name: "propertyOwner", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertyOwnersColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_personPropertyOwners ON personPropertyOwners (person, propertyOwner)",
		},
		ListRule:   types.Pointer(ppoListViewRule),
		ViewRule:   types.Pointer(ppoListViewRule),
		CreateRule: types.Pointer(ppoCUDRule),
		UpdateRule: types.Pointer(ppoCUDRule),
		DeleteRule: types.Pointer(ppoCUDRule),
	}
	if err := dao.SaveCollection(personPropOwnersCol); err != nil {
		return err
	}

	// 10. agents  (unique on person + cobrand)
	// PocketBase v0.22 can't handle 3+ back-relations in a chain, so property-based
	// access for agents is omitted here (it would require 3 back-relations).
	agentListViewRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || person.user.id = @request.auth.id"

	agentCreateUpdateRule := "cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	agentDeleteRule := "" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || person.user.id = @request.auth.id"

	agentsCol := &models.Collection{
		Name: "agents",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "agt01", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "agt02", Name: "cobrand", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: cobrandsColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_agents ON agents (person, cobrand)",
		},
		ListRule:   types.Pointer(agentListViewRule),
		ViewRule:   types.Pointer(agentListViewRule),
		CreateRule: types.Pointer(agentCreateUpdateRule),
		UpdateRule: types.Pointer(agentCreateUpdateRule),
		DeleteRule: types.Pointer(agentDeleteRule),
	}
	if err := dao.SaveCollection(agentsCol); err != nil {
		return err
	}
	agentsColID := agentsCol.Id

	// 11. propertyAgents  (unique on agent + property)
	paListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agent.person.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	paCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	paDeleteRule := paCreateUpdateRule + " || agent.person.user.id = @request.auth.id"

	propertyAgentsCol := &models.Collection{
		Name: "propertyAgents",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "pagt1", Name: "agent", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: agentsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "pagt2", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_propertyAgents ON propertyAgents (agent, property)",
		},
		ListRule:   types.Pointer(paListViewRule),
		ViewRule:   types.Pointer(paListViewRule),
		CreateRule: types.Pointer(paCreateUpdateRule),
		UpdateRule: types.Pointer(paCreateUpdateRule),
		DeleteRule: types.Pointer(paDeleteRule),
	}
	if err := dao.SaveCollection(propertyAgentsCol); err != nil {
		return err
	}

	// 12. households  (unique on person + property)
	hhListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id"

	hhCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	hhDeleteRule := hhCreateUpdateRule + " || person.user.id = @request.auth.id"

	householdsCol := &models.Collection{
		Name: "households",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "hhld1", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "hhld2", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_households ON households (person, property)",
		},
		ListRule:   types.Pointer(hhListViewRule),
		ViewRule:   types.Pointer(hhListViewRule),
		CreateRule: types.Pointer(hhCreateUpdateRule),
		UpdateRule: types.Pointer(hhCreateUpdateRule),
		DeleteRule: types.Pointer(hhDeleteRule),
	}
	if err := dao.SaveCollection(householdsCol); err != nil {
		return err
	}

	// 13. tenants  (unique on person + property)
	tenListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id"

	tenCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"

	tenDeleteRule := tenCreateUpdateRule + " || person.user.id = @request.auth.id"

	tenantsCol := &models.Collection{
		Name: "tenants",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "tnt01", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "tnt02", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_tenants ON tenants (person, property)",
		},
		ListRule:   types.Pointer(tenListViewRule),
		ViewRule:   types.Pointer(tenListViewRule),
		CreateRule: types.Pointer(tenCreateUpdateRule),
		UpdateRule: types.Pointer(tenCreateUpdateRule),
		DeleteRule: types.Pointer(tenDeleteRule),
	}
	if err := dao.SaveCollection(tenantsCol); err != nil {
		return err
	}

	// 14. items
	// 3+ back-relation chains cause SQL errors in PocketBase v0.22, so property-owner
	// access is omitted (requires 3 backs via propertyItems→property→propertyOwners→...).
	// Residents (tenant/household/agent) are reachable via 2 back-relations.
	itemListViewRule := "" +
		"personItems_via_item.person.user.id = @request.auth.id" +
		" || propertyItems_via_item.property.tenants_via_property.person.user.id = @request.auth.id" +
		" || propertyItems_via_item.property.households_via_property.person.user.id = @request.auth.id" +
		" || propertyItems_via_item.property.propertyAgents_via_property.agent.person.user.id = @request.auth.id"

	itemUpdateDeleteRule := "personItems_via_item.person.user.id = @request.auth.id"

	itemsCol := &models.Collection{
		Name: "items",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "itm01", Name: "name", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
			&schema.SchemaField{Id: "itm02", Name: "description", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
			&schema.SchemaField{Id: "itm03", Name: "picture", Type: schema.FieldTypeFile, Required: false,
				Options: &schema.FileOptions{MaxSelect: 1, MaxSize: 5242880}},
		),
		ListRule:   types.Pointer(itemListViewRule),
		ViewRule:   types.Pointer(itemListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(itemUpdateDeleteRule),
		DeleteRule: types.Pointer(itemUpdateDeleteRule),
	}
	if err := dao.SaveCollection(itemsCol); err != nil {
		return err
	}
	itemsColID := itemsCol.Id

	// 15. entryDevices  (unique on item)
	// PocketBase v0.22: forward→back→forward→back chains from entryDevices fail.
	// Only the forward→back chain (item → personItems) is supported.
	edListViewRule := "item.personItems_via_item.person.user.id = @request.auth.id"
	edCUDRule := "item.personItems_via_item.person.user.id = @request.auth.id"

	entryDevicesCol := &models.Collection{
		Name: "entryDevices",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "edv01", Name: "item", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: itemsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "edv02", Name: "deviceType", Type: schema.FieldTypeSelect, Required: true,
				Options: &schema.SelectOptions{MaxSelect: 1, Values: []string{"Fob", "Key", "Remote", "RoomKey", "MailboxKey"}}},
			&schema.SchemaField{Id: "edv03", Name: "identifier", Type: schema.FieldTypeText, Required: true, Options: &schema.TextOptions{}},
			&schema.SchemaField{Id: "edv04", Name: "defunctReason", Type: schema.FieldTypeSelect, Required: true,
				Options: &schema.SelectOptions{MaxSelect: 1, Values: []string{"None", "Lost", "Damaged", "Retired", "Stolen"}}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_entryDevices_item ON entryDevices (item)",
		},
		ListRule:   types.Pointer(edListViewRule),
		ViewRule:   types.Pointer(edListViewRule),
		CreateRule: types.Pointer(edCUDRule),
		UpdateRule: types.Pointer(edCUDRule),
		DeleteRule: types.Pointer(edCUDRule),
	}
	if err := dao.SaveCollection(entryDevicesCol); err != nil {
		return err
	}

	// 16. propertyItems  (unique on item + property)
	piListViewRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.propertyAgents_via_property.agent.person.user.id = @request.auth.id"

	piCUDRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id"

	propertyItemsCol := &models.Collection{
		Name: "propertyItems",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "pitm1", Name: "item", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: itemsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "pitm2", Name: "property", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: propertiesColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_propertyItems ON propertyItems (item, property)",
		},
		ListRule:   types.Pointer(piListViewRule),
		ViewRule:   types.Pointer(piListViewRule),
		CreateRule: types.Pointer(piCUDRule),
		UpdateRule: types.Pointer(piCUDRule),
		DeleteRule: types.Pointer(piCUDRule),
	}
	if err := dao.SaveCollection(propertyItemsCol); err != nil {
		return err
	}

	// 17. personItems  (unique on person + item)
	// Property-owner access requires 3+ back-relations; omitted due to PocketBase v0.22 limitation.
	psListViewRule := "person.user.id = @request.auth.id"
	psUpdateDeleteRule := "person.user.id = @request.auth.id"

	personItemsCol := &models.Collection{
		Name: "personItems",
		Type: models.CollectionTypeBase,
		Schema: schema.NewSchema(
			&schema.SchemaField{Id: "psit1", Name: "person", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: personsColID, MaxSelect: types.Pointer(1)}},
			&schema.SchemaField{Id: "psit2", Name: "item", Type: schema.FieldTypeRelation, Required: true,
				Options: &schema.RelationOptions{CollectionId: itemsColID, MaxSelect: types.Pointer(1)}},
		),
		Indexes: types.JsonArray[string]{
			"CREATE UNIQUE INDEX idx_personItems ON personItems (person, item)",
		},
		ListRule:   types.Pointer(psListViewRule),
		ViewRule:   types.Pointer(psListViewRule),
		CreateRule: types.Pointer("@request.auth.id != \"\""),
		UpdateRule: types.Pointer(psUpdateDeleteRule),
		DeleteRule: types.Pointer(psUpdateDeleteRule),
	}
	if err := dao.SaveCollection(personItemsCol); err != nil {
		return err
	}

	return nil
}

func down1748563200(db dbx.Builder) error {
	dao := daos.New(db)

	for _, name := range []string{
		"personItems", "propertyItems", "entryDevices", "items",
		"tenants", "households", "propertyAgents", "agents",
		"personPropertyOwners", "cobrandPropertyOwners",
		"cobrandPropertyManagers", "cobrandAdmins", "cobrands",
		"propertyOwners", "properties", "persons",
	} {
		col, err := dao.FindCollectionByNameOrId(name)
		if err != nil {
			continue
		}
		if err := dao.DeleteCollection(col); err != nil {
			return err
		}
	}

	// Restore users name and avatar fields.
	usersCol, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	usersCol.Schema.AddField(&schema.SchemaField{
		Id:      "users_name",
		Type:    schema.FieldTypeText,
		Name:    "name",
		Options: &schema.TextOptions{},
	})
	usersCol.Schema.AddField(&schema.SchemaField{
		Id:   "users_avatar",
		Type: schema.FieldTypeFile,
		Name: "avatar",
		Options: &schema.FileOptions{
			MaxSelect: 1,
			MaxSize:   5242880,
			MimeTypes: []string{"image/jpeg", "image/png", "image/svg+xml", "image/gif", "image/webp"},
		},
	})
	return dao.SaveCollection(usersCol)
}
