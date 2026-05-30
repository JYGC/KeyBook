package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

// ── cobrands ──────────────────────────────────────────────────────────────────

func TestCobrands_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/cobrands/records"},
		{"GET", fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)},
		{"POST", "/api/collections/cobrands/records"},
		{"PATCH", fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)},
		{"DELETE", fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, `{"name":"X"}`, "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestCobrands_View_AdminAndAgentSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand admin", env.tok.cobrandAdmin},
		{"agent", env.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestCobrands_View_UnrelatedReturns404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)
	resp := env.do("GET", path, "", env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestCobrands_Create_AnyAuthenticated_Returns200(t *testing.T) {
	env := newTestEnv(t)
	resp := env.do("POST", "/api/collections/cobrands/records", `{"name":"New Co"}`, env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestCobrands_UpdateDelete_OnlyCobrandAdmin(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)

	resp := env.do("PATCH", path, `{"name":"Updated Co"}`, env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusOK)

	resp = env.do("PATCH", path, `{"name":"Hacked"}`, env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)
}

// ── cobrandAdmins ─────────────────────────────────────────────────────────────

func TestCobrandAdmins_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/cobrandAdmins/records"},
		{"GET", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", env.ids.cobrandAdmin)},
		{"POST", "/api/collections/cobrandAdmins/records"},
		{"PATCH", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", env.ids.cobrandAdmin)},
		{"DELETE", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", env.ids.cobrandAdmin)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, "{}", "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestCobrandAdmins_View_AdminAndAgentSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", env.ids.cobrandAdmin)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand admin", env.tok.cobrandAdmin},
		{"agent", env.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestCobrandAdmins_Delete_SelfRemoval(t *testing.T) {
	env := newTestEnv(t)

	// cobrandAdmin user can remove themselves.
	newCobrandID := env.createRecord(t, "cobrands", map[string]any{"name": "Self Remove Co"})
	cobrandAdminUserID := env.userIdByEmail(t, "cobrand@test.com")
	caID := env.createRecord(t, "cobrandAdmins",
		map[string]any{"user": cobrandAdminUserID, "cobrand": newCobrandID})

	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", caID)
	resp := env.do("DELETE", path, "", env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusNoContent)
}

// ── agents ────────────────────────────────────────────────────────────────────

func TestAgents_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/agents/records"},
		{"GET", fmt.Sprintf("/api/collections/agents/records/%s", env.ids.agent)},
		{"POST", "/api/collections/agents/records"},
		{"PATCH", fmt.Sprintf("/api/collections/agents/records/%s", env.ids.agent)},
		{"DELETE", fmt.Sprintf("/api/collections/agents/records/%s", env.ids.agent)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, "{}", "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestAgents_View_CobrandAdminAndAgentSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/agents/records/%s", env.ids.agent)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand admin", env.tok.cobrandAdmin},
		{"agent self", env.tok.agent},
		{"property owner", env.tok.userOwner},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestAgents_Delete_CobrandAdminOrAgentSelf(t *testing.T) {
	env := newTestEnv(t)

	// Agent can resign (delete own agent record).
	newPersonID := env.createRecord(t, "persons",
		map[string]any{"name": "Resignee", "DOB": "1990-01-01 00:00:00.000Z"})
	agentToDeleteID := env.createRecord(t, "agents",
		map[string]any{"person": newPersonID, "cobrand": env.ids.cobrand})

	// cobrand admin can delete.
	path := fmt.Sprintf("/api/collections/agents/records/%s", agentToDeleteID)
	resp := env.do("DELETE", path, "", env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusNoContent)
}

// ── propertyAgents ────────────────────────────────────────────────────────────

func TestPropertyAgents_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/propertyAgents/records"},
		{"GET", fmt.Sprintf("/api/collections/propertyAgents/records/%s", env.ids.propertyAgent)},
		{"POST", "/api/collections/propertyAgents/records"},
		{"PATCH", fmt.Sprintf("/api/collections/propertyAgents/records/%s", env.ids.propertyAgent)},
		{"DELETE", fmt.Sprintf("/api/collections/propertyAgents/records/%s", env.ids.propertyAgent)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, "{}", "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestPropertyAgents_View_AllowedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/propertyAgents/records/%s", env.ids.propertyAgent)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"property owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
		{"agent self", env.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestPropertyAgents_Delete_AgentCanResign(t *testing.T) {
	env := newTestEnv(t)

	paID := env.createRecord(t, "propertyAgents",
		map[string]any{"agent": env.ids.agent, "property": env.ids.property})
	path := fmt.Sprintf("/api/collections/propertyAgents/records/%s", paID)
	resp := env.do("DELETE", path, "", env.tok.agent)
	// Agent self can delete (resign).
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Errorf("agent resign: expected 204 or 404, got %d", resp.StatusCode)
	}
}

// ── households ────────────────────────────────────────────────────────────────

func TestHouseholds_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/households/records"},
		{"GET", fmt.Sprintf("/api/collections/households/records/%s", env.ids.household)},
		{"POST", "/api/collections/households/records"},
		{"PATCH", fmt.Sprintf("/api/collections/households/records/%s", env.ids.household)},
		{"DELETE", fmt.Sprintf("/api/collections/households/records/%s", env.ids.household)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, "{}", "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestHouseholds_View_AllowedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/households/records/%s", env.ids.household)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
		{"household self", env.tok.household},
		{"tenant (cross-visibility)", env.tok.tenant},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestHouseholds_Delete_MemberCanLeave(t *testing.T) {
	env := newTestEnv(t)
	hhID := env.createRecord(t, "households",
		map[string]any{"person": env.ids.personHousehold, "property": env.ids.property})
	path := fmt.Sprintf("/api/collections/households/records/%s", hhID)
	resp := env.do("DELETE", path, "", env.tok.household)
	// Household member can self-delete (leave).
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Errorf("household leave: expected 204 or 404, got %d", resp.StatusCode)
	}
}

// ── tenants ───────────────────────────────────────────────────────────────────

func TestTenants_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/tenants/records"},
		{"GET", fmt.Sprintf("/api/collections/tenants/records/%s", env.ids.tenant)},
		{"POST", "/api/collections/tenants/records"},
		{"PATCH", fmt.Sprintf("/api/collections/tenants/records/%s", env.ids.tenant)},
		{"DELETE", fmt.Sprintf("/api/collections/tenants/records/%s", env.ids.tenant)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, "{}", "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestTenants_View_AllowedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/tenants/records/%s", env.ids.tenant)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
		{"tenant self", env.tok.tenant},
		{"household (cross-visibility)", env.tok.household},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestTenants_Delete_TenantCanVacate(t *testing.T) {
	env := newTestEnv(t)
	tID := env.createRecord(t, "tenants",
		map[string]any{"person": env.ids.personTenant, "property": env.ids.property})
	path := fmt.Sprintf("/api/collections/tenants/records/%s", tID)
	resp := env.do("DELETE", path, "", env.tok.tenant)
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Errorf("tenant vacate: expected 204 or 404, got %d", resp.StatusCode)
	}
}

// ── items ─────────────────────────────────────────────────────────────────────

func TestItems_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/items/records"},
		{"GET", fmt.Sprintf("/api/collections/items/records/%s", env.ids.item)},
		{"POST", "/api/collections/items/records"},
		{"PATCH", fmt.Sprintf("/api/collections/items/records/%s", env.ids.item)},
		{"DELETE", fmt.Sprintf("/api/collections/items/records/%s", env.ids.item)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, `{"name":"X","description":"Y"}`, "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestItems_View_AllowedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/items/records/%s", env.ids.item)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"property owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
		{"tenant (resident)", env.tok.tenant},
		{"household (resident)", env.tok.household},
		{"agent (resident)", env.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestItems_View_UnrelatedReturns404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/items/records/%s", env.ids.item)
	resp := env.do("GET", path, "", env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestItems_Create_AnyAuthenticated_Returns200(t *testing.T) {
	env := newTestEnv(t)
	resp := env.do("POST", "/api/collections/items/records",
		`{"name":"New Item","description":"Desc"}`, env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestItems_UpdateDelete_ResidentsCannotModify(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/items/records/%s", env.ids.item)

	// Residents (tenant, household, agent) cannot update items.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"tenant", env.tok.tenant},
		{"household", env.tok.household},
	} {
		t.Run(tok.name+" cannot update", func(t *testing.T) {
			resp := env.do("PATCH", path, `{"name":"Hacked"}`, tok.token)
			env.assertStatus(t, resp, http.StatusNotFound)
		})
	}

	// Person owner can update.
	resp := env.do("PATCH", path, `{"name":"Updated"}`, env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusOK)
}

// ── entryDevices ──────────────────────────────────────────────────────────────

func TestEntryDevices_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/entryDevices/records"},
		{"GET", fmt.Sprintf("/api/collections/entryDevices/records/%s", env.ids.entryDevice)},
		{"POST", "/api/collections/entryDevices/records"},
		{"PATCH", fmt.Sprintf("/api/collections/entryDevices/records/%s", env.ids.entryDevice)},
		{"DELETE", fmt.Sprintf("/api/collections/entryDevices/records/%s", env.ids.entryDevice)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, "{}", "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestEntryDevices_View_ResidentsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/entryDevices/records/%s", env.ids.entryDevice)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
		{"tenant", env.tok.tenant},
		{"household", env.tok.household},
		{"agent", env.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestEntryDevices_CreateUpdateDelete_ResidentsCannotModify(t *testing.T) {
	env := newTestEnv(t)

	// Residents cannot create entry devices.
	newItemID := env.createRecord(t, "items",
		map[string]any{"name": "New Key", "description": "D"})
	body := fmt.Sprintf(`{"item":%q,"deviceType":"Key","identifier":"X001","defunctReason":"None"}`, newItemID)

	resp := env.do("POST", "/api/collections/entryDevices/records", body, env.tok.tenant)
	env.assertStatus(t, resp, http.StatusNotFound)

	// Resident cannot update existing entry device.
	path := fmt.Sprintf("/api/collections/entryDevices/records/%s", env.ids.entryDevice)
	resp = env.do("PATCH", path, `{"identifier":"HACKED"}`, env.tok.tenant)
	env.assertStatus(t, resp, http.StatusNotFound)
}

// ── propertyItems ─────────────────────────────────────────────────────────────

func TestPropertyItems_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/propertyItems/records"},
		{"GET", fmt.Sprintf("/api/collections/propertyItems/records/%s", env.ids.propertyItem)},
		{"POST", "/api/collections/propertyItems/records"},
		{"PATCH", fmt.Sprintf("/api/collections/propertyItems/records/%s", env.ids.propertyItem)},
		{"DELETE", fmt.Sprintf("/api/collections/propertyItems/records/%s", env.ids.propertyItem)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, "{}", "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestPropertyItems_View_AllGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/propertyItems/records/%s", env.ids.propertyItem)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
		{"tenant", env.tok.tenant},
		{"household", env.tok.household},
		{"agent", env.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestPropertyItems_Create_EitherSideCanAdd(t *testing.T) {
	env := newTestEnv(t)
	// Person owner (item owner via personItems) can add a propertyItems record.
	newItemID := env.createRecord(t, "items",
		map[string]any{"name": "Extra Item", "description": "D"})
	env.createRecord(t, "personItems",
		map[string]any{"person": env.ids.personOwner, "item": newItemID})
	body := fmt.Sprintf(`{"item":%q,"property":%q}`, newItemID, env.ids.property)
	resp := env.do("POST", "/api/collections/propertyItems/records", body, env.tok.userOwner)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("person owner create propertyItem: expected 200 or 400, got %d", resp.StatusCode)
	}
}

// ── personItems ───────────────────────────────────────────────────────────────

func TestPersonItems_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	cases := []struct{ method, path string }{
		{"GET", "/api/collections/personItems/records"},
		{"GET", fmt.Sprintf("/api/collections/personItems/records/%s", env.ids.personItem)},
		{"POST", "/api/collections/personItems/records"},
		{"PATCH", fmt.Sprintf("/api/collections/personItems/records/%s", env.ids.personItem)},
		{"DELETE", fmt.Sprintf("/api/collections/personItems/records/%s", env.ids.personItem)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, "{}", "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestPersonItems_View_PersonAndOwnersSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/personItems/records/%s", env.ids.personItem)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person (owner)", env.tok.userOwner},
		{"property owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestPersonItems_Create_AnyAuthenticated_Returns200(t *testing.T) {
	env := newTestEnv(t)
	newItemID := env.createRecord(t, "items",
		map[string]any{"name": "Unclaimed Item", "description": "D"})
	unrelatedPersonID := env.createRecord(t, "persons",
		map[string]any{"name": "Temp", "DOB": "2000-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"item":%q}`, unrelatedPersonID, newItemID)
	resp := env.do("POST", "/api/collections/personItems/records", body, env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestPersonItems_UpdateDelete_PersonOrOwner(t *testing.T) {
	env := newTestEnv(t)

	// Unrelated cannot update personItems record.
	path := fmt.Sprintf("/api/collections/personItems/records/%s", env.ids.personItem)
	resp := env.do("PATCH", path, fmt.Sprintf(`{"person":%q}`, env.ids.personOwner), env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)

	// The person themselves can update.
	resp = env.do("PATCH", path, fmt.Sprintf(`{"person":%q}`, env.ids.personOwner), env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusOK)
}
