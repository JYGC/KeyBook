package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(up1783150478, down1783150478)
}

// up1783150478 adds cobrandAdmins.approved (default false), backfills existing
// rows to approved = true so already-established admins are not retroactively
// locked out, locks the cobrandAdmins update rule to superusers only, and
// rewrites every access rule that authorizes a create/update/delete via
// "cobrandAdmins_via_cobrand.user.id = @request.auth.id" to also require
// approved = true on that same relation path. List/view rules and
// self-resignation clauses are left functionally unchanged (see design.md
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
func up1783150478(db dbx.Builder) error {
	dao := daos.New(db)

	// 1. Add the approved field to cobrandAdmins.
	cobrandAdminsCol, err := dao.FindCollectionByNameOrId("cobrandAdmins")
	if err != nil {
		return err
	}
	cobrandAdminsCol.Schema.AddField(&schema.SchemaField{
		Id:      "cbad3",
		Name:    "approved",
		Type:    schema.FieldTypeBool,
		Options: &schema.BoolOptions{},
	})

	// 2. List/view unaffected (fix operator only). Lock update to superusers
	// only; gate delete's "remove another admin" clause while leaving
	// self-resignation ("user.id = @request.auth.id") ungated.
	caListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || cobrand.agents_via_cobrand.person.user.id = @request.auth.id")
	cobrandAdminsCol.ListRule = caListViewRule
	cobrandAdminsCol.ViewRule = caListViewRule
	cobrandAdminsCol.UpdateRule = nil
	cobrandAdminsCol.DeleteRule = types.Pointer("" +
		"(cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || user.id = @request.auth.id")
	if err := dao.SaveCollection(cobrandAdminsCol); err != nil {
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
	cobrandsCol, err := dao.FindCollectionByNameOrId("cobrands")
	if err != nil {
		return err
	}
	cobrandListViewRule := types.Pointer("" +
		"cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || agents_via_cobrand.person.user.id = @request.auth.id")
	cobrandUpdateDeleteRule := types.Pointer(
		"cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrandAdmins_via_cobrand.approved ?= true")
	cobrandsCol.ListRule = cobrandListViewRule
	cobrandsCol.ViewRule = cobrandListViewRule
	cobrandsCol.UpdateRule = cobrandUpdateDeleteRule
	cobrandsCol.DeleteRule = cobrandUpdateDeleteRule
	if err := dao.SaveCollection(cobrandsCol); err != nil {
		return err
	}

	// 5. properties — list/view fix operator only; update/delete: gate the
	// cobrand-owner clause.
	propertiesCol, err := dao.FindCollectionByNameOrId("properties")
	if err != nil {
		return err
	}
	propListViewRule := types.Pointer("" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || tenants_via_property.person.user.id = @request.auth.id" +
		" || households_via_property.person.user.id = @request.auth.id" +
		" || propertyAgents_via_property.agent.person.user.id = @request.auth.id")
	propUpdateDeleteRule := types.Pointer("" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)")
	propertiesCol.ListRule = propListViewRule
	propertiesCol.ViewRule = propListViewRule
	propertiesCol.UpdateRule = propUpdateDeleteRule
	propertiesCol.DeleteRule = propUpdateDeleteRule
	if err := dao.SaveCollection(propertiesCol); err != nil {
		return err
	}

	// 6. propertyOwners — list/view fix operator only; update/delete: gate the
	// cobrand-owner clause.
	propertyOwnersCol, err := dao.FindCollectionByNameOrId("propertyOwners")
	if err != nil {
		return err
	}
	poListViewRule := types.Pointer("" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	poUpdateDeleteRule := types.Pointer("" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)")
	propertyOwnersCol.ListRule = poListViewRule
	propertyOwnersCol.ViewRule = poListViewRule
	propertyOwnersCol.UpdateRule = poUpdateDeleteRule
	propertyOwnersCol.DeleteRule = poUpdateDeleteRule
	if err := dao.SaveCollection(propertyOwnersCol); err != nil {
		return err
	}

	// 7. cobrandPropertyManagers — list/view fix operator only (2 refs);
	// create/update: gate the owning-cobrand clause. delete: same gated
	// clause, plus the managing cobrand's own resignation clause (fix
	// operator, stays ungated).
	cpmCol, err := dao.FindCollectionByNameOrId("cobrandPropertyManagers")
	if err != nil {
		return err
	}
	cpmListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	cpmCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)"
	cpmDeleteRule := cpmCreateUpdateRule + " || cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id"
	cpmCol.ListRule = cpmListViewRule
	cpmCol.ViewRule = cpmListViewRule
	cpmCol.CreateRule = types.Pointer(cpmCreateUpdateRule)
	cpmCol.UpdateRule = types.Pointer(cpmCreateUpdateRule)
	cpmCol.DeleteRule = types.Pointer(cpmDeleteRule)
	if err := dao.SaveCollection(cpmCol); err != nil {
		return err
	}

	// 8. cobrandPropertyOwners — list/view fix operator only (2 refs);
	// create/update: gate the cobrand clause. delete keeps its original
	// business logic (unchanged — the resignation path), but still needs the
	// operator fix since it references the same to-many relation.
	cpoCol, err := dao.FindCollectionByNameOrId("cobrandPropertyOwners")
	if err != nil {
		return err
	}
	cpoListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	cpoCreateUpdateRule := types.Pointer("" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrand.cobrandAdmins_via_cobrand.approved ?= true)")
	cpoDeleteRule := types.Pointer("" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	cpoCol.ListRule = cpoListViewRule
	cpoCol.ViewRule = cpoListViewRule
	cpoCol.CreateRule = cpoCreateUpdateRule
	cpoCol.UpdateRule = cpoCreateUpdateRule
	cpoCol.DeleteRule = cpoDeleteRule
	if err := dao.SaveCollection(cpoCol); err != nil {
		return err
	}

	// 9. personPropertyOwners — list/view fix operator only (2 refs);
	// create/update/delete share one rule, gate the cobrand-owner clause.
	ppoCol, err := dao.FindCollectionByNameOrId("personPropertyOwners")
	if err != nil {
		return err
	}
	ppoListViewRule := types.Pointer("" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	ppoCUDRule := types.Pointer("" +
		"person.user.id = @request.auth.id" +
		" || (propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)")
	ppoCol.ListRule = ppoListViewRule
	ppoCol.ViewRule = ppoListViewRule
	ppoCol.CreateRule = ppoCUDRule
	ppoCol.UpdateRule = ppoCUDRule
	ppoCol.DeleteRule = ppoCUDRule
	if err := dao.SaveCollection(ppoCol); err != nil {
		return err
	}

	// 10. agents — list/view fix operator only (1 ref; property-based access
	// was never implemented here — PocketBase v0.22 can't handle the
	// resulting 3+ back-relation chain, see 1748563200_full_collections.go).
	// create/update: gate the single cobrand-admin clause. delete: gate the
	// admin-removes-agent clause; the agent's own resignation clause is left
	// ungated.
	agentsCol, err := dao.FindCollectionByNameOrId("agents")
	if err != nil {
		return err
	}
	agentListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || person.user.id = @request.auth.id")
	agentCreateUpdateRule := types.Pointer(
		"cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrand.cobrandAdmins_via_cobrand.approved ?= true")
	agentDeleteRule := types.Pointer("" +
		"(cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || person.user.id = @request.auth.id")
	agentsCol.ListRule = agentListViewRule
	agentsCol.ViewRule = agentListViewRule
	agentsCol.CreateRule = agentCreateUpdateRule
	agentsCol.UpdateRule = agentCreateUpdateRule
	agentsCol.DeleteRule = agentDeleteRule
	if err := dao.SaveCollection(agentsCol); err != nil {
		return err
	}

	// 11. propertyAgents — list/view fix operator only (3 refs); create/update:
	// gate all three cobrand-based clauses. delete: same gated clauses, plus
	// the agent's own resignation ungated.
	paCol, err := dao.FindCollectionByNameOrId("propertyAgents")
	if err != nil {
		return err
	}
	paListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || agent.person.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id")
	paCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (agent.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id && agent.cobrand.cobrandAdmins_via_cobrand.approved ?= true)"
	paDeleteRule := paCreateUpdateRule + " || agent.person.user.id = @request.auth.id"
	paCol.ListRule = paListViewRule
	paCol.ViewRule = paListViewRule
	paCol.CreateRule = types.Pointer(paCreateUpdateRule)
	paCol.UpdateRule = types.Pointer(paCreateUpdateRule)
	paCol.DeleteRule = types.Pointer(paDeleteRule)
	if err := dao.SaveCollection(paCol); err != nil {
		return err
	}

	// 12. households — list/view fix operator only (2 refs); create/update:
	// gate the cobrand-owner and cobrand-manager clauses. delete: same, plus
	// the member's own resignation clause left ungated.
	hhCol, err := dao.FindCollectionByNameOrId("households")
	if err != nil {
		return err
	}
	hhListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id")
	hhCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.approved ?= true)"
	hhDeleteRule := hhCreateUpdateRule + " || person.user.id = @request.auth.id"
	hhCol.ListRule = hhListViewRule
	hhCol.ViewRule = hhListViewRule
	hhCol.CreateRule = types.Pointer(hhCreateUpdateRule)
	hhCol.UpdateRule = types.Pointer(hhCreateUpdateRule)
	hhCol.DeleteRule = types.Pointer(hhDeleteRule)
	if err := dao.SaveCollection(hhCol); err != nil {
		return err
	}

	// 13. tenants — same shape as households.
	tenCol, err := dao.FindCollectionByNameOrId("tenants")
	if err != nil {
		return err
	}
	tenListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id")
	tenCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.approved ?= true)"
	tenDeleteRule := tenCreateUpdateRule + " || person.user.id = @request.auth.id"
	tenCol.ListRule = tenListViewRule
	tenCol.ViewRule = tenListViewRule
	tenCol.CreateRule = types.Pointer(tenCreateUpdateRule)
	tenCol.UpdateRule = types.Pointer(tenCreateUpdateRule)
	tenCol.DeleteRule = types.Pointer(tenDeleteRule)
	if err := dao.SaveCollection(tenCol); err != nil {
		return err
	}

	// 14. propertyItems — list/view fix operator only (2 refs);
	// create/update/delete share one rule, gate the cobrand-owner and
	// cobrand-manager clauses. The person-owner clauses (property's and
	// item's) are unaffected.
	piCol, err := dao.FindCollectionByNameOrId("propertyItems")
	if err != nil {
		return err
	}
	piListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.propertyAgents_via_property.agent.person.user.id = @request.auth.id")
	piCUDRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || (property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || (property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id ?= @request.auth.id" +
		" && property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.approved ?= true)" +
		" || item.personItems_via_item.person.user.id = @request.auth.id")
	piCol.ListRule = piListViewRule
	piCol.ViewRule = piListViewRule
	piCol.CreateRule = piCUDRule
	piCol.UpdateRule = piCUDRule
	piCol.DeleteRule = piCUDRule
	return dao.SaveCollection(piCol)
}

// down1783150478 restores every rule modified by up1783150478 (including
// list/view) to its exact original string from 1748563200_full_collections.go,
// and drops the approved field.
func down1783150478(db dbx.Builder) error {
	dao := daos.New(db)

	propertyItemsCol, err := dao.FindCollectionByNameOrId("propertyItems")
	if err != nil {
		return err
	}
	piListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.propertyAgents_via_property.agent.person.user.id = @request.auth.id")
	piCUDRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || item.personItems_via_item.person.user.id = @request.auth.id")
	propertyItemsCol.ListRule = piListViewRule
	propertyItemsCol.ViewRule = piListViewRule
	propertyItemsCol.CreateRule = piCUDRule
	propertyItemsCol.UpdateRule = piCUDRule
	propertyItemsCol.DeleteRule = piCUDRule
	if err := dao.SaveCollection(propertyItemsCol); err != nil {
		return err
	}

	tenantsCol, err := dao.FindCollectionByNameOrId("tenants")
	if err != nil {
		return err
	}
	tenListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id")
	tenCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"
	tenantsCol.ListRule = tenListViewRule
	tenantsCol.ViewRule = tenListViewRule
	tenantsCol.CreateRule = types.Pointer(tenCreateUpdateRule)
	tenantsCol.UpdateRule = types.Pointer(tenCreateUpdateRule)
	tenantsCol.DeleteRule = types.Pointer(tenCreateUpdateRule + " || person.user.id = @request.auth.id")
	if err := dao.SaveCollection(tenantsCol); err != nil {
		return err
	}

	householdsCol, err := dao.FindCollectionByNameOrId("households")
	if err != nil {
		return err
	}
	hhListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.households_via_property.person.user.id = @request.auth.id" +
		" || property.tenants_via_property.person.user.id = @request.auth.id")
	hhCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"
	householdsCol.ListRule = hhListViewRule
	householdsCol.ViewRule = hhListViewRule
	householdsCol.CreateRule = types.Pointer(hhCreateUpdateRule)
	householdsCol.UpdateRule = types.Pointer(hhCreateUpdateRule)
	householdsCol.DeleteRule = types.Pointer(hhCreateUpdateRule + " || person.user.id = @request.auth.id")
	if err := dao.SaveCollection(householdsCol); err != nil {
		return err
	}

	propertyAgentsCol, err := dao.FindCollectionByNameOrId("propertyAgents")
	if err != nil {
		return err
	}
	paListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agent.person.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	paCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agent.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"
	propertyAgentsCol.ListRule = paListViewRule
	propertyAgentsCol.ViewRule = paListViewRule
	propertyAgentsCol.CreateRule = types.Pointer(paCreateUpdateRule)
	propertyAgentsCol.UpdateRule = types.Pointer(paCreateUpdateRule)
	propertyAgentsCol.DeleteRule = types.Pointer(paCreateUpdateRule + " || agent.person.user.id = @request.auth.id")
	if err := dao.SaveCollection(propertyAgentsCol); err != nil {
		return err
	}

	agentsCol, err := dao.FindCollectionByNameOrId("agents")
	if err != nil {
		return err
	}
	agentListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || person.user.id = @request.auth.id" +
		" || propertyAgents.property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyAgents.property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyAgents.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	agentCreateUpdateRule := types.Pointer("cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	agentsCol.ListRule = agentListViewRule
	agentsCol.ViewRule = agentListViewRule
	agentsCol.CreateRule = agentCreateUpdateRule
	agentsCol.UpdateRule = agentCreateUpdateRule
	agentsCol.DeleteRule = types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || person.user.id = @request.auth.id")
	if err := dao.SaveCollection(agentsCol); err != nil {
		return err
	}

	personPropertyOwnersCol, err := dao.FindCollectionByNameOrId("personPropertyOwners")
	if err != nil {
		return err
	}
	ppoListViewRule := types.Pointer("" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	ppoCUDRule := types.Pointer("" +
		"person.user.id = @request.auth.id" +
		" || propertyOwner.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	personPropertyOwnersCol.ListRule = ppoListViewRule
	personPropertyOwnersCol.ViewRule = ppoListViewRule
	personPropertyOwnersCol.CreateRule = ppoCUDRule
	personPropertyOwnersCol.UpdateRule = ppoCUDRule
	personPropertyOwnersCol.DeleteRule = ppoCUDRule
	if err := dao.SaveCollection(personPropertyOwnersCol); err != nil {
		return err
	}

	cobrandPropertyOwnersCol, err := dao.FindCollectionByNameOrId("cobrandPropertyOwners")
	if err != nil {
		return err
	}
	cpoListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwner.property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cpoCreateUpdateRule := types.Pointer("" +
		"propertyOwner.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cobrandPropertyOwnersCol.ListRule = cpoListViewRule
	cobrandPropertyOwnersCol.ViewRule = cpoListViewRule
	cobrandPropertyOwnersCol.CreateRule = cpoCreateUpdateRule
	cobrandPropertyOwnersCol.UpdateRule = cpoCreateUpdateRule
	cobrandPropertyOwnersCol.DeleteRule = cpoCreateUpdateRule
	if err := dao.SaveCollection(cobrandPropertyOwnersCol); err != nil {
		return err
	}

	cobrandPropertyManagersCol, err := dao.FindCollectionByNameOrId("cobrandPropertyManagers")
	if err != nil {
		return err
	}
	cpmListViewRule := types.Pointer("" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cpmCreateUpdateRule := "" +
		"property.propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || property.propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id"
	cobrandPropertyManagersCol.ListRule = cpmListViewRule
	cobrandPropertyManagersCol.ViewRule = cpmListViewRule
	cobrandPropertyManagersCol.CreateRule = types.Pointer(cpmCreateUpdateRule)
	cobrandPropertyManagersCol.UpdateRule = types.Pointer(cpmCreateUpdateRule)
	cobrandPropertyManagersCol.DeleteRule = types.Pointer(
		cpmCreateUpdateRule + " || cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	if err := dao.SaveCollection(cobrandPropertyManagersCol); err != nil {
		return err
	}

	propertyOwnersCol, err := dao.FindCollectionByNameOrId("propertyOwners")
	if err != nil {
		return err
	}
	poListViewRule := types.Pointer("" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || property.cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	poUpdateDeleteRule := types.Pointer("" +
		"personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	propertyOwnersCol.ListRule = poListViewRule
	propertyOwnersCol.ViewRule = poListViewRule
	propertyOwnersCol.UpdateRule = poUpdateDeleteRule
	propertyOwnersCol.DeleteRule = poUpdateDeleteRule
	if err := dao.SaveCollection(propertyOwnersCol); err != nil {
		return err
	}

	propertiesCol, err := dao.FindCollectionByNameOrId("properties")
	if err != nil {
		return err
	}
	propListViewRule := types.Pointer("" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || tenants_via_property.person.user.id = @request.auth.id" +
		" || households_via_property.person.user.id = @request.auth.id" +
		" || propertyAgents_via_property.agent.person.user.id = @request.auth.id")
	propUpdateDeleteRule := types.Pointer("" +
		"propertyOwners_via_property.personPropertyOwners_via_propertyOwner.person.user.id = @request.auth.id" +
		" || propertyOwners_via_property.cobrandPropertyOwners_via_propertyOwner.cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	propertiesCol.ListRule = propListViewRule
	propertiesCol.ViewRule = propListViewRule
	propertiesCol.UpdateRule = propUpdateDeleteRule
	propertiesCol.DeleteRule = propUpdateDeleteRule
	if err := dao.SaveCollection(propertiesCol); err != nil {
		return err
	}

	cobrandsCol, err := dao.FindCollectionByNameOrId("cobrands")
	if err != nil {
		return err
	}
	cobrandListViewRule := types.Pointer("" +
		"cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || agents_via_cobrand.person.user.id = @request.auth.id")
	cobrandUpdateDeleteRule := types.Pointer("cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cobrandsCol.ListRule = cobrandListViewRule
	cobrandsCol.ViewRule = cobrandListViewRule
	cobrandsCol.UpdateRule = cobrandUpdateDeleteRule
	cobrandsCol.DeleteRule = cobrandUpdateDeleteRule
	if err := dao.SaveCollection(cobrandsCol); err != nil {
		return err
	}

	cobrandAdminsCol, err := dao.FindCollectionByNameOrId("cobrandAdmins")
	if err != nil {
		return err
	}
	caListViewRule := types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || cobrand.agents_via_cobrand.person.user.id = @request.auth.id")
	cobrandAdminsCol.ListRule = caListViewRule
	cobrandAdminsCol.ViewRule = caListViewRule
	cobrandAdminsCol.UpdateRule = types.Pointer("cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id")
	cobrandAdminsCol.DeleteRule = types.Pointer("" +
		"cobrand.cobrandAdmins_via_cobrand.user.id = @request.auth.id" +
		" || user.id = @request.auth.id")
	if field := cobrandAdminsCol.Schema.GetFieldByName("approved"); field != nil {
		cobrandAdminsCol.Schema.RemoveField(field.Id)
	}
	return dao.SaveCollection(cobrandAdminsCol)
}
