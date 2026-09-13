package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(
		addCobrandAdminApprovalAndGateWriteRules,
		dropCobrandAdminApprovalAndUngateWriteRules,
	)
}

// Existing rows are backfilled to approved = true so admins established before
// approval existed as a concept are not retroactively locked out. List/view
// rules and self-resignation clauses stay functionally unchanged (design.md
// section 2.4).
//
// Every comparison against cobrandAdmins_via_cobrand (in both read and write
// rules, whether or not its business logic changes here) is switched from the
// plain "=" operator to PocketBase's "?=" any-match operator. A cobrand can
// now have more than one cobrandAdmins row (an approved admin alongside an
// unapproved one), and PocketBase's plain "=" on a to-many back-relation
// requires ALL related rows to satisfy the condition (a multi-match ambiguity
// guard), not just one — so with two admins, a plain "=" check would reject
// even the admin who really does match. "?=" disables that guard and matches
// on any one related row; when two "?=" conditions share the same
// "..._via_cobrand" prefix (e.g. .user.id and .approved together), PocketBase
// reuses the same joined table alias for both, so they still correlate to the
// same row rather than being satisfiable by two different admins.
func addCobrandAdminApprovalAndGateWriteRules(db dbx.Builder) error {
	dao := daos.New(db)

	// 1. Add the approved field to cobrandAdmins.
	cobrandAdminsCollection, err := dao.FindCollectionByNameOrId("cobrandAdmins")
	if err != nil {
		return err
	}
	cobrandAdminsCollection.Schema.AddField(&schema.SchemaField{
		Id:      "cbad3",
		Name:    "approved",
		Type:    schema.FieldTypeBool,
		Options: &schema.BoolOptions{},
	})

	// 2. List/view unaffected (fix operator only). Lock update to superusers
	// only; gate delete's "remove another admin" clause while leaving
	// self-resignation ("user.id = @request.auth.id") ungated.
	cobrandAdminsListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || cobrand.agents_via_cobrand.person.user.id = @request.auth.id")
	cobrandAdminsCollection.ListRule = cobrandAdminsListViewRule
	cobrandAdminsCollection.ViewRule = cobrandAdminsListViewRule
	cobrandAdminsCollection.UpdateRule = nil
	cobrandAdminsCollection.DeleteRule = types.Pointer("" +
		"(cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || user.id = @request.auth.id")
	if err := dao.SaveCollection(cobrandAdminsCollection); err != nil {
		return err
	}

	// 3. Backfill: existing admins were established before approval existed as
	// a concept, so they are not retroactively locked out. Only cobrandAdmins
	// rows created after this migration default to approved = false.
	if _, err := db.NewQuery("UPDATE {{cobrandAdmins}} SET approved = 1").Execute(); err != nil {
		return err
	}

	// 4. cobrands — list/view fix operator only; update/delete restricted to
	// approved admins.
	cobrandsCollection, err := dao.FindCollectionByNameOrId("cobrands")
	if err != nil {
		return err
	}
	cobrandsListViewRule := types.Pointer("" +
		"cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || agents_via_cobrand.person.user.id = @request.auth.id")
	cobrandsUpdateDeleteRule := types.Pointer(
		"cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrandAdmins_via_cobrand.approved ?= true")
	cobrandsCollection.ListRule = cobrandsListViewRule
	cobrandsCollection.ViewRule = cobrandsListViewRule
	cobrandsCollection.UpdateRule = cobrandsUpdateDeleteRule
	cobrandsCollection.DeleteRule = cobrandsUpdateDeleteRule
	if err := dao.SaveCollection(cobrandsCollection); err != nil {
		return err
	}

	// 5. properties — list/view fix operator only; update/delete: gate the
	// cobrand-owner clause.
	propertiesCollection, err := dao.FindCollectionByNameOrId("properties")
	if err != nil {
		return err
	}
	propertiesListViewRule := types.Pointer("" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || tenants_via_property.person.user.id = @request.auth.id" +
		" || households_via_property.person.user.id = @request.auth.id" +
		" || propertyAgents_via_property.agent.person.user.id = @request.auth.id")
	propertiesUpdateDeleteRule := types.Pointer("" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)")
	propertiesCollection.ListRule = propertiesListViewRule
	propertiesCollection.ViewRule = propertiesListViewRule
	propertiesCollection.UpdateRule = propertiesUpdateDeleteRule
	propertiesCollection.DeleteRule = propertiesUpdateDeleteRule
	if err := dao.SaveCollection(propertiesCollection); err != nil {
		return err
	}

	// 6. propertyOwners — list/view fix operator only; update/delete: gate the
	// cobrand-owner clause.
	propertyOwnersCollection, err := dao.FindCollectionByNameOrId("propertyOwners")
	if err != nil {
		return err
	}
	propertyOwnersListViewRule := types.Pointer("" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	propertyOwnersUpdateDeleteRule := types.Pointer("" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)")
	propertyOwnersCollection.ListRule = propertyOwnersListViewRule
	propertyOwnersCollection.ViewRule = propertyOwnersListViewRule
	propertyOwnersCollection.UpdateRule = propertyOwnersUpdateDeleteRule
	propertyOwnersCollection.DeleteRule = propertyOwnersUpdateDeleteRule
	if err := dao.SaveCollection(propertyOwnersCollection); err != nil {
		return err
	}

	// 7. cobrandPropertyManagers — list/view fix operator only (2 refs);
	// create/update: gate the owning-cobrand clause. delete: same gated
	// clause, plus the managing cobrand's own resignation clause (fix
	// operator, stays ungated).
	cobrandPropertyManagersCollection, err := dao.FindCollectionByNameOrId("cobrandPropertyManagers")
	if err != nil {
		return err
	}
	cobrandPropertyManagersListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	cobrandPropertyManagersCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)"
	cobrandPropertyManagersDeleteRule := cobrandPropertyManagersCreateUpdateRule + " || cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id"
	cobrandPropertyManagersCollection.ListRule = cobrandPropertyManagersListViewRule
	cobrandPropertyManagersCollection.ViewRule = cobrandPropertyManagersListViewRule
	cobrandPropertyManagersCollection.CreateRule = types.Pointer(cobrandPropertyManagersCreateUpdateRule)
	cobrandPropertyManagersCollection.UpdateRule = types.Pointer(cobrandPropertyManagersCreateUpdateRule)
	cobrandPropertyManagersCollection.DeleteRule = types.Pointer(cobrandPropertyManagersDeleteRule)
	if err := dao.SaveCollection(cobrandPropertyManagersCollection); err != nil {
		return err
	}

	// 8. cobrandPropertyOwners — list/view fix operator only (2 refs);
	// create/update: gate the cobrand clause. delete keeps its original
	// business logic (unchanged — the resignation path), but still needs the
	// operator fix since it references the same to-many relation.
	cobrandPropertyOwnersCollection, err := dao.FindCollectionByNameOrId("cobrandPropertyOwners")
	if err != nil {
		return err
	}
	cobrandPropertyOwnersListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	cobrandPropertyOwnersCreateUpdateRule := types.Pointer("" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrand.cobrandAdmins_via_cobrand.approved ?= true)")
	cobrandPropertyOwnersDeleteRule := types.Pointer("" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	cobrandPropertyOwnersCollection.ListRule = cobrandPropertyOwnersListViewRule
	cobrandPropertyOwnersCollection.ViewRule = cobrandPropertyOwnersListViewRule
	cobrandPropertyOwnersCollection.CreateRule = cobrandPropertyOwnersCreateUpdateRule
	cobrandPropertyOwnersCollection.UpdateRule = cobrandPropertyOwnersCreateUpdateRule
	cobrandPropertyOwnersCollection.DeleteRule = cobrandPropertyOwnersDeleteRule
	if err := dao.SaveCollection(cobrandPropertyOwnersCollection); err != nil {
		return err
	}

	// 9. personPropertyOwners — list/view fix operator only (2 refs);
	// create/update/delete share one rule, gate the cobrand-owner clause.
	personPropertyOwnersCollection, err := dao.FindCollectionByNameOrId("personPropertyOwners")
	if err != nil {
		return err
	}
	personPropertyOwnersListViewRule := types.Pointer("" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	personPropertyOwnersCreateUpdateDeleteRule := types.Pointer("" +
		"person.user.id = @request.auth.id" +
		" || (propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)")
	personPropertyOwnersCollection.ListRule = personPropertyOwnersListViewRule
	personPropertyOwnersCollection.ViewRule = personPropertyOwnersListViewRule
	personPropertyOwnersCollection.CreateRule = personPropertyOwnersCreateUpdateDeleteRule
	personPropertyOwnersCollection.UpdateRule = personPropertyOwnersCreateUpdateDeleteRule
	personPropertyOwnersCollection.DeleteRule = personPropertyOwnersCreateUpdateDeleteRule
	if err := dao.SaveCollection(personPropertyOwnersCollection); err != nil {
		return err
	}

	// 10. agents — list/view fix operator only (1 ref; property-based access
	// was never implemented here — PocketBase v0.22 can't handle the
	// resulting 3+ back-relation chain, see 1748563200_full_collections.go).
	// create/update: gate the single cobrand-admin clause. delete: gate the
	// admin-removes-agent clause; the agent's own resignation clause is left
	// ungated.
	agentsCollection, err := dao.FindCollectionByNameOrId("agents")
	if err != nil {
		return err
	}
	agentsListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || person.user.id = @request.auth.id")
	agentsCreateUpdateRule := types.Pointer(
		"cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrand.cobrandAdmins_via_cobrand.approved ?= true")
	agentsDeleteRule := types.Pointer("" +
		"(cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || person.user.id = @request.auth.id")
	agentsCollection.ListRule = agentsListViewRule
	agentsCollection.ViewRule = agentsListViewRule
	agentsCollection.CreateRule = agentsCreateUpdateRule
	agentsCollection.UpdateRule = agentsCreateUpdateRule
	agentsCollection.DeleteRule = agentsDeleteRule
	if err := dao.SaveCollection(agentsCollection); err != nil {
		return err
	}

	// 11. propertyAgents — list/view fix operator only (3 refs); create/update:
	// gate all three cobrand-based clauses. delete: same gated clauses, plus
	// the agent's own resignation ungated.
	propertyAgentsCollection, err := dao.FindCollectionByNameOrId("propertyAgents")
	if err != nil {
		return err
	}
	propertyAgentsListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || agent.person.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	propertyAgentsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (agent.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && agent.cobrand.cobrandAdmins_via_cobrand.approved ?= true)"
	propertyAgentsDeleteRule := propertyAgentsCreateUpdateRule + " || agent.person.user.id = @request.auth.id"
	propertyAgentsCollection.ListRule = propertyAgentsListViewRule
	propertyAgentsCollection.ViewRule = propertyAgentsListViewRule
	propertyAgentsCollection.CreateRule = types.Pointer(propertyAgentsCreateUpdateRule)
	propertyAgentsCollection.UpdateRule = types.Pointer(propertyAgentsCreateUpdateRule)
	propertyAgentsCollection.DeleteRule = types.Pointer(propertyAgentsDeleteRule)
	if err := dao.SaveCollection(propertyAgentsCollection); err != nil {
		return err
	}

	// 12. households — list/view fix operator only (2 refs); create/update:
	// gate the cobrand-owner and cobrand-manager clauses. delete: same, plus
	// the member's own resignation clause left ungated.
	householdsCollection, err := dao.FindCollectionByNameOrId("households")
	if err != nil {
		return err
	}
	householdsListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id")
	householdsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.approved ?= true)"
	householdsDeleteRule := householdsCreateUpdateRule + " || person.user.id = @request.auth.id"
	householdsCollection.ListRule = householdsListViewRule
	householdsCollection.ViewRule = householdsListViewRule
	householdsCollection.CreateRule = types.Pointer(householdsCreateUpdateRule)
	householdsCollection.UpdateRule = types.Pointer(householdsCreateUpdateRule)
	householdsCollection.DeleteRule = types.Pointer(householdsDeleteRule)
	if err := dao.SaveCollection(householdsCollection); err != nil {
		return err
	}

	// 13. tenants — same shape as households.
	tenantsCollection, err := dao.FindCollectionByNameOrId("tenants")
	if err != nil {
		return err
	}
	tenantsListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id")
	tenantsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.approved ?= true)"
	tenantsDeleteRule := tenantsCreateUpdateRule + " || person.user.id = @request.auth.id"
	tenantsCollection.ListRule = tenantsListViewRule
	tenantsCollection.ViewRule = tenantsListViewRule
	tenantsCollection.CreateRule = types.Pointer(tenantsCreateUpdateRule)
	tenantsCollection.UpdateRule = types.Pointer(tenantsCreateUpdateRule)
	tenantsCollection.DeleteRule = types.Pointer(tenantsDeleteRule)
	if err := dao.SaveCollection(tenantsCollection); err != nil {
		return err
	}

	// 14. propertyItems — list/view fix operator only (2 refs);
	// create/update/delete share one rule, gate the cobrand-owner and
	// cobrand-manager clauses. The person-owner clauses (property's and
	// item's) are unaffected.
	propertyItemsCollection, err := dao.FindCollectionByNameOrId("propertyItems")
	if err != nil {
		return err
	}
	propertyItemsListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.propertyAgents_via_property.agent.person.user.id = @request.auth.id")
	propertyItemsCreateUpdateDeleteRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || item.personItems_via_item.person.user.id = @request.auth.id")
	propertyItemsCollection.ListRule = propertyItemsListViewRule
	propertyItemsCollection.ViewRule = propertyItemsListViewRule
	propertyItemsCollection.CreateRule = propertyItemsCreateUpdateDeleteRule
	propertyItemsCollection.UpdateRule = propertyItemsCreateUpdateDeleteRule
	propertyItemsCollection.DeleteRule = propertyItemsCreateUpdateDeleteRule
	return dao.SaveCollection(propertyItemsCollection)
}

// Every restored rule string — including list/view, whose business logic never
// changed — is copied verbatim from 1748563200_full_collections.go.
func dropCobrandAdminApprovalAndUngateWriteRules(db dbx.Builder) error {
	dao := daos.New(db)

	propertyItemsCollection, err := dao.FindCollectionByNameOrId("propertyItems")
	if err != nil {
		return err
	}
	propertyItemsListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.propertyAgents_via_property.agent.person.user.id = @request.auth.id")
	propertyItemsCreateUpdateDeleteRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id")
	propertyItemsCollection.ListRule = propertyItemsListViewRule
	propertyItemsCollection.ViewRule = propertyItemsListViewRule
	propertyItemsCollection.CreateRule = propertyItemsCreateUpdateDeleteRule
	propertyItemsCollection.UpdateRule = propertyItemsCreateUpdateDeleteRule
	propertyItemsCollection.DeleteRule = propertyItemsCreateUpdateDeleteRule
	if err := dao.SaveCollection(propertyItemsCollection); err != nil {
		return err
	}

	tenantsCollection, err := dao.FindCollectionByNameOrId("tenants")
	if err != nil {
		return err
	}
	tenantsListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id")
	tenantsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"
	tenantsCollection.ListRule = tenantsListViewRule
	tenantsCollection.ViewRule = tenantsListViewRule
	tenantsCollection.CreateRule = types.Pointer(tenantsCreateUpdateRule)
	tenantsCollection.UpdateRule = types.Pointer(tenantsCreateUpdateRule)
	tenantsCollection.DeleteRule = types.Pointer(tenantsCreateUpdateRule + " || person.user.id = @request.auth.id")
	if err := dao.SaveCollection(tenantsCollection); err != nil {
		return err
	}

	householdsCollection, err := dao.FindCollectionByNameOrId("households")
	if err != nil {
		return err
	}
	householdsListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id")
	householdsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"
	householdsCollection.ListRule = householdsListViewRule
	householdsCollection.ViewRule = householdsListViewRule
	householdsCollection.CreateRule = types.Pointer(householdsCreateUpdateRule)
	householdsCollection.UpdateRule = types.Pointer(householdsCreateUpdateRule)
	householdsCollection.DeleteRule = types.Pointer(householdsCreateUpdateRule + " || person.user.id = @request.auth.id")
	if err := dao.SaveCollection(householdsCollection); err != nil {
		return err
	}

	propertyAgentsCollection, err := dao.FindCollectionByNameOrId("propertyAgents")
	if err != nil {
		return err
	}
	propertyAgentsListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agent.person.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	propertyAgentsCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"
	propertyAgentsCollection.ListRule = propertyAgentsListViewRule
	propertyAgentsCollection.ViewRule = propertyAgentsListViewRule
	propertyAgentsCollection.CreateRule = types.Pointer(propertyAgentsCreateUpdateRule)
	propertyAgentsCollection.UpdateRule = types.Pointer(propertyAgentsCreateUpdateRule)
	propertyAgentsCollection.DeleteRule = types.Pointer(propertyAgentsCreateUpdateRule + " || agent.person.user.id = @request.auth.id")
	if err := dao.SaveCollection(propertyAgentsCollection); err != nil {
		return err
	}

	agentsCollection, err := dao.FindCollectionByNameOrId("agents")
	if err != nil {
		return err
	}
	agentsListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || person.user.id = @request.auth.id" +
		" || propertyAgents.property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyAgents.property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyAgents.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	agentsCreateUpdateRule := types.Pointer("cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	agentsCollection.ListRule = agentsListViewRule
	agentsCollection.ViewRule = agentsListViewRule
	agentsCollection.CreateRule = agentsCreateUpdateRule
	agentsCollection.UpdateRule = agentsCreateUpdateRule
	agentsCollection.DeleteRule = types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || person.user.id = @request.auth.id")
	if err := dao.SaveCollection(agentsCollection); err != nil {
		return err
	}

	personPropertyOwnersCollection, err := dao.FindCollectionByNameOrId("personPropertyOwners")
	if err != nil {
		return err
	}
	personPropertyOwnersListViewRule := types.Pointer("" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	personPropertyOwnersCreateUpdateDeleteRule := types.Pointer("" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	personPropertyOwnersCollection.ListRule = personPropertyOwnersListViewRule
	personPropertyOwnersCollection.ViewRule = personPropertyOwnersListViewRule
	personPropertyOwnersCollection.CreateRule = personPropertyOwnersCreateUpdateDeleteRule
	personPropertyOwnersCollection.UpdateRule = personPropertyOwnersCreateUpdateDeleteRule
	personPropertyOwnersCollection.DeleteRule = personPropertyOwnersCreateUpdateDeleteRule
	if err := dao.SaveCollection(personPropertyOwnersCollection); err != nil {
		return err
	}

	cobrandPropertyOwnersCollection, err := dao.FindCollectionByNameOrId("cobrandPropertyOwners")
	if err != nil {
		return err
	}
	cobrandPropertyOwnersListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cobrandPropertyOwnersCreateUpdateRule := types.Pointer("" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cobrandPropertyOwnersCollection.ListRule = cobrandPropertyOwnersListViewRule
	cobrandPropertyOwnersCollection.ViewRule = cobrandPropertyOwnersListViewRule
	cobrandPropertyOwnersCollection.CreateRule = cobrandPropertyOwnersCreateUpdateRule
	cobrandPropertyOwnersCollection.UpdateRule = cobrandPropertyOwnersCreateUpdateRule
	cobrandPropertyOwnersCollection.DeleteRule = cobrandPropertyOwnersCreateUpdateRule
	if err := dao.SaveCollection(cobrandPropertyOwnersCollection); err != nil {
		return err
	}

	cobrandPropertyManagersCollection, err := dao.FindCollectionByNameOrId("cobrandPropertyManagers")
	if err != nil {
		return err
	}
	cobrandPropertyManagersListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cobrandPropertyManagersCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"
	cobrandPropertyManagersCollection.ListRule = cobrandPropertyManagersListViewRule
	cobrandPropertyManagersCollection.ViewRule = cobrandPropertyManagersListViewRule
	cobrandPropertyManagersCollection.CreateRule = types.Pointer(cobrandPropertyManagersCreateUpdateRule)
	cobrandPropertyManagersCollection.UpdateRule = types.Pointer(cobrandPropertyManagersCreateUpdateRule)
	cobrandPropertyManagersCollection.DeleteRule = types.Pointer(
		cobrandPropertyManagersCreateUpdateRule + " || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	if err := dao.SaveCollection(cobrandPropertyManagersCollection); err != nil {
		return err
	}

	propertyOwnersCollection, err := dao.FindCollectionByNameOrId("propertyOwners")
	if err != nil {
		return err
	}
	propertyOwnersListViewRule := types.Pointer("" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	propertyOwnersUpdateDeleteRule := types.Pointer("" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	propertyOwnersCollection.ListRule = propertyOwnersListViewRule
	propertyOwnersCollection.ViewRule = propertyOwnersListViewRule
	propertyOwnersCollection.UpdateRule = propertyOwnersUpdateDeleteRule
	propertyOwnersCollection.DeleteRule = propertyOwnersUpdateDeleteRule
	if err := dao.SaveCollection(propertyOwnersCollection); err != nil {
		return err
	}

	propertiesCollection, err := dao.FindCollectionByNameOrId("properties")
	if err != nil {
		return err
	}
	propertiesListViewRule := types.Pointer("" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || tenants_via_property.person.user.id = @request.auth.id" +
		" || households_via_property.person.user.id = @request.auth.id" +
		" || propertyAgents_via_property.agent.person.user.id = @request.auth.id")
	propertiesUpdateDeleteRule := types.Pointer("" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	propertiesCollection.ListRule = propertiesListViewRule
	propertiesCollection.ViewRule = propertiesListViewRule
	propertiesCollection.UpdateRule = propertiesUpdateDeleteRule
	propertiesCollection.DeleteRule = propertiesUpdateDeleteRule
	if err := dao.SaveCollection(propertiesCollection); err != nil {
		return err
	}

	cobrandsCollection, err := dao.FindCollectionByNameOrId("cobrands")
	if err != nil {
		return err
	}
	cobrandsListViewRule := types.Pointer("" +
		"cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agents_via_cobrand.person.user.id = @request.auth.id")
	cobrandsUpdateDeleteRule := types.Pointer("cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cobrandsCollection.ListRule = cobrandsListViewRule
	cobrandsCollection.ViewRule = cobrandsListViewRule
	cobrandsCollection.UpdateRule = cobrandsUpdateDeleteRule
	cobrandsCollection.DeleteRule = cobrandsUpdateDeleteRule
	if err := dao.SaveCollection(cobrandsCollection); err != nil {
		return err
	}

	cobrandAdminsCollection, err := dao.FindCollectionByNameOrId("cobrandAdmins")
	if err != nil {
		return err
	}
	cobrandAdminsListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.agents_via_cobrand.person.user.id = @request.auth.id")
	cobrandAdminsCollection.ListRule = cobrandAdminsListViewRule
	cobrandAdminsCollection.ViewRule = cobrandAdminsListViewRule
	cobrandAdminsCollection.UpdateRule = types.Pointer("cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cobrandAdminsCollection.DeleteRule = types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || user.id = @request.auth.id")
	if field := cobrandAdminsCollection.Schema.GetFieldByName("approved"); field != nil {
		cobrandAdminsCollection.Schema.RemoveField(field.Id)
	}
	return dao.SaveCollection(cobrandAdminsCollection)
}
