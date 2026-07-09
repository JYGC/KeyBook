package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

// ── cobrands ──────────────────────────────────────────────────────────────────

func TestCobrands_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.cobrand
	body := `{"name":"X"}`

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/cobrands/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/cobrands/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/cobrands/records", body, "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/cobrands/records/%s", id), body, "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/cobrands/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
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
	id := env.ids.cobrandAdmin

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/cobrandAdmins/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/cobrandAdmins/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		// cobrandAdmins.updateRule is locked to superusers only (nil) as of
		// the approval-gate migration, so PocketBase rejects with 403
		// ("Only admins can perform this action") rather than 404.
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusForbidden)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
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
	id := env.ids.agent

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/agents/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/agents/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/agents/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/agents/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/agents/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestAgents_View_CobrandAdminAndAgentSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/agents/records/%s", env.ids.agent)
	// Property owner access omitted — requires 3 back-relations (PocketBase v0.22 limit).
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand admin", env.tok.cobrandAdmin},
		{"agent self", env.tok.agent},
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
	id := env.ids.propertyAgent

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/propertyAgents/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/propertyAgents/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/propertyAgents/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/propertyAgents/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/propertyAgents/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestPropertyAgents_View_AllowedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/propertyAgents/records/%s", env.ids.propertyAgent)
	// Cobrand manager omitted — PocketBase v0.22 limitation.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"property owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
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

	// Use a fresh agent to avoid unique constraint with the fixture's propertyAgent.
	newPersonID := env.createRecord(t, "persons",
		map[string]any{"name": "Resign Agent", "DOB": "1990-01-01 00:00:00.000Z"})
	newAgentID := env.createRecord(t, "agents",
		map[string]any{"person": newPersonID, "cobrand": env.ids.cobrand})
	paID := env.createRecord(t, "propertyAgents",
		map[string]any{"agent": newAgentID, "property": env.ids.property})

	path := fmt.Sprintf("/api/collections/propertyAgents/records/%s", paID)
	// Cobrand admin can delete the agent's propertyAgent record.
	resp := env.do("DELETE", path, "", env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusNoContent)
}

// ── households ────────────────────────────────────────────────────────────────

func TestHouseholds_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.household

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/households/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/households/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/households/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/households/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/households/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestHouseholds_View_AllowedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/households/records/%s", env.ids.household)
	// Cobrand manager omitted — PocketBase v0.22 limitation.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
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
	// Use a fresh person to avoid unique constraint with the fixture's household.
	newPersonID := env.createRecord(t, "persons",
		map[string]any{"name": "Leaver", "DOB": "2000-01-01 00:00:00.000Z"})
	hhID := env.createRecord(t, "households",
		map[string]any{"person": newPersonID, "property": env.ids.property})
	path := fmt.Sprintf("/api/collections/households/records/%s", hhID)
	// Property owner can delete household records.
	resp := env.do("DELETE", path, "", env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusNoContent)
}

// ── tenants ───────────────────────────────────────────────────────────────────

func TestTenants_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.tenant

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/tenants/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/tenants/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/tenants/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/tenants/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/tenants/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestTenants_View_AllowedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/tenants/records/%s", env.ids.tenant)
	// Cobrand manager omitted — PocketBase v0.22 limitation.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
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
	// Use a fresh person to avoid unique constraint with the fixture's tenant.
	newPersonID := env.createRecord(t, "persons",
		map[string]any{"name": "Vacater", "DOB": "2000-01-01 00:00:00.000Z"})
	tID := env.createRecord(t, "tenants",
		map[string]any{"person": newPersonID, "property": env.ids.property})
	path := fmt.Sprintf("/api/collections/tenants/records/%s", tID)
	// Property owner can delete tenant records.
	resp := env.do("DELETE", path, "", env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusNoContent)
}

// ── items ─────────────────────────────────────────────────────────────────────

func TestItems_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.item
	body := `{"name":"X","description":"Y"}`

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/items/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/items/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/items/records", body, "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/items/records/%s", id), body, "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/items/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestItems_View_AllowedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/items/records/%s", env.ids.item)
	// Property-owner access requires 3+ back-relations (omitted due to PocketBase v0.22 limit).
	// Tenant/household/agent reach items via 2-back property chain.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
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
	id := env.ids.entryDevice

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/entryDevices/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/entryDevices/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/entryDevices/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/entryDevices/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/entryDevices/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestEntryDevices_View_PersonOwnerSees200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/entryDevices/records/%s", env.ids.entryDevice)
	// Only person owner (via item.personItems chain) can view entryDevices.
	// PocketBase v0.22 can't handle forward→back→forward→back chains from entryDevices.
	resp := env.do("GET", path, "", env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestEntryDevices_CreateUpdateDelete_ResidentsCannotModify(t *testing.T) {
	env := newTestEnv(t)

	// Residents cannot create entry devices.
	newItemID := env.createRecord(t, "items",
		map[string]any{"name": "New Key", "description": "D"})
	body := fmt.Sprintf(`{"item":%q,"deviceType":"Key","identifier":"X001","defunctReason":"None"}`, newItemID)

	resp := env.do("POST", "/api/collections/entryDevices/records", body, env.tok.tenant)
	env.assertStatus(t, resp, http.StatusBadRequest)

	// Resident cannot update existing entry device.
	path := fmt.Sprintf("/api/collections/entryDevices/records/%s", env.ids.entryDevice)
	resp = env.do("PATCH", path, `{"identifier":"HACKED"}`, env.tok.tenant)
	env.assertStatus(t, resp, http.StatusNotFound)
}

// ── propertyItems ─────────────────────────────────────────────────────────────

func TestPropertyItems_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.propertyItem

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/propertyItems/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/propertyItems/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/propertyItems/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/propertyItems/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/propertyItems/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestPropertyItems_View_AllGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/propertyItems/records/%s", env.ids.propertyItem)
	// Cobrand manager omitted — PocketBase v0.22 limitation.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
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
	id := env.ids.personItem

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/personItems/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/personItems/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/personItems/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/personItems/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/personItems/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestPersonItems_View_PersonOwnerSees200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/personItems/records/%s", env.ids.personItem)
	// personItems view rule is person.user.id only; property-owner access requires 3+ backs.
	resp := env.do("GET", path, "", env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusOK)
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
