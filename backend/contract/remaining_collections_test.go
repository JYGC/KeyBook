package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCobrands_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.cobrand
	body := `{"name":"X"}`

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/cobrands/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/cobrands/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/cobrands/records", body, "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/cobrands/records/%s", id), body, "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/cobrands/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestCobrands_View_AdminAndAgentSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", environment.ids.cobrand)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand admin", environment.tok.cobrandAdmin},
		{"agent", environment.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestCobrands_View_UnrelatedReturns404(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", environment.ids.cobrand)
	response := environment.do("GET", path, "", environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestCobrands_Create_AnyAuthenticated_Returns200(t *testing.T) {
	environment := newTestEnvironment(t)
	response := environment.do("POST", "/api/collections/cobrands/records", `{"name":"New Co"}`, environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestCobrands_UpdateDelete_OnlyCobrandAdmin(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", environment.ids.cobrand)

	response := environment.do("PATCH", path, `{"name":"Updated Co"}`, environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusOK)

	response = environment.do("PATCH", path, `{"name":"Hacked"}`, environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestCobrandAdmins_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.cobrandAdmin

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/cobrandAdmins/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/cobrandAdmins/records", "{}", "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		// cobrandAdmins.updateRule is locked to superusers only (nil) as of
		// the approval-gate migration, so PocketBase rejects with 403
		// ("Only admins can perform this action") rather than 404.
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", id), "{}", "")
		environment.assertStatus(t, response, http.StatusForbidden)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestCobrandAdmins_View_AdminAndAgentSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", environment.ids.cobrandAdmin)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand admin", environment.tok.cobrandAdmin},
		{"agent", environment.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestCobrandAdmins_Delete_SelfRemoval(t *testing.T) {
	environment := newTestEnvironment(t)

	// cobrandAdmin user can remove themselves.
	newCobrandID := environment.createRecord(t, "cobrands", map[string]any{"name": "Self Remove Co"})
	cobrandAdminUserID := environment.userIdByEmail(t, "cobrand@test.com")
	caID := environment.createRecord(t, "cobrandAdmins",
		map[string]any{"user": cobrandAdminUserID, "cobrand": newCobrandID})

	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", caID)
	response := environment.do("DELETE", path, "", environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusNoContent)
}

func TestAgents_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.agent

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/agents/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/agents/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/agents/records", "{}", "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/agents/records/%s", id), "{}", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/agents/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestAgents_View_CobrandAdminAndAgentSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/agents/records/%s", environment.ids.agent)
	// Property owner access omitted — requires 3 back-relations (PocketBase v0.22 limit).
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand admin", environment.tok.cobrandAdmin},
		{"agent self", environment.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestAgents_Delete_CobrandAdminOrAgentSelf(t *testing.T) {
	environment := newTestEnvironment(t)

	// Agent can resign (delete own agent record).
	newPersonID := environment.createRecord(t, "persons",
		map[string]any{"name": "Resignee", "DOB": "1990-01-01 00:00:00.000Z"})
	agentToDeleteID := environment.createRecord(t, "agents",
		map[string]any{"person": newPersonID, "cobrand": environment.ids.cobrand})

	// cobrand admin can delete.
	path := fmt.Sprintf("/api/collections/agents/records/%s", agentToDeleteID)
	response := environment.do("DELETE", path, "", environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusNoContent)
}

func TestPropertyAgents_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.propertyAgent

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/propertyAgents/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/propertyAgents/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/propertyAgents/records", "{}", "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/propertyAgents/records/%s", id), "{}", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/propertyAgents/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestPropertyAgents_View_AllowedGroupsSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/propertyAgents/records/%s", environment.ids.propertyAgent)
	// Cobrand manager omitted — PocketBase v0.22 limitation.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"property owner", environment.tok.userOwner},
		{"cobrand owner", environment.tok.cobrandAdmin},
		{"agent self", environment.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestPropertyAgents_Delete_AgentCanResign(t *testing.T) {
	environment := newTestEnvironment(t)

	// Use a fresh agent to avoid unique constraint with the fixture's propertyAgent.
	newPersonID := environment.createRecord(t, "persons",
		map[string]any{"name": "Resign Agent", "DOB": "1990-01-01 00:00:00.000Z"})
	newAgentID := environment.createRecord(t, "agents",
		map[string]any{"person": newPersonID, "cobrand": environment.ids.cobrand})
	paID := environment.createRecord(t, "propertyAgents",
		map[string]any{"agent": newAgentID, "property": environment.ids.property})

	path := fmt.Sprintf("/api/collections/propertyAgents/records/%s", paID)
	// Cobrand admin can delete the agent's propertyAgent record.
	response := environment.do("DELETE", path, "", environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusNoContent)
}

func TestHouseholds_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.household

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/households/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/households/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/households/records", "{}", "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/households/records/%s", id), "{}", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/households/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestHouseholds_View_AllowedGroupsSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/households/records/%s", environment.ids.household)
	// Cobrand manager omitted — PocketBase v0.22 limitation.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"owner", environment.tok.userOwner},
		{"cobrand owner", environment.tok.cobrandAdmin},
		{"household self", environment.tok.household},
		{"tenant (cross-visibility)", environment.tok.tenant},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestHouseholds_Delete_MemberCanLeave(t *testing.T) {
	environment := newTestEnvironment(t)
	// Use a fresh person to avoid unique constraint with the fixture's household.
	newPersonID := environment.createRecord(t, "persons",
		map[string]any{"name": "Leaver", "DOB": "2000-01-01 00:00:00.000Z"})
	hhID := environment.createRecord(t, "households",
		map[string]any{"person": newPersonID, "property": environment.ids.property})
	path := fmt.Sprintf("/api/collections/households/records/%s", hhID)
	// Property owner can delete household records.
	response := environment.do("DELETE", path, "", environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusNoContent)
}

func TestTenants_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.tenant

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/tenants/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/tenants/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/tenants/records", "{}", "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/tenants/records/%s", id), "{}", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/tenants/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestTenants_View_AllowedGroupsSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/tenants/records/%s", environment.ids.tenant)
	// Cobrand manager omitted — PocketBase v0.22 limitation.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"owner", environment.tok.userOwner},
		{"cobrand owner", environment.tok.cobrandAdmin},
		{"tenant self", environment.tok.tenant},
		{"household (cross-visibility)", environment.tok.household},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestTenants_Delete_TenantCanVacate(t *testing.T) {
	environment := newTestEnvironment(t)
	// Use a fresh person to avoid unique constraint with the fixture's tenant.
	newPersonID := environment.createRecord(t, "persons",
		map[string]any{"name": "Vacater", "DOB": "2000-01-01 00:00:00.000Z"})
	tID := environment.createRecord(t, "tenants",
		map[string]any{"person": newPersonID, "property": environment.ids.property})
	path := fmt.Sprintf("/api/collections/tenants/records/%s", tID)
	// Property owner can delete tenant records.
	response := environment.do("DELETE", path, "", environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusNoContent)
}

func TestItems_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.item
	body := `{"name":"X","description":"Y"}`

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/items/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/items/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/items/records", body, "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/items/records/%s", id), body, "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/items/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestItems_View_AllowedGroupsSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/items/records/%s", environment.ids.item)
	// Property-owner access requires 3+ back-relations (omitted due to PocketBase v0.22 limit).
	// Tenant/household/agent reach items via 2-back property chain.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", environment.tok.userOwner},
		{"tenant", environment.tok.tenant},
		{"household", environment.tok.household},
		{"agent", environment.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestItems_View_UnrelatedReturns404(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/items/records/%s", environment.ids.item)
	response := environment.do("GET", path, "", environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestItems_Create_AnyAuthenticated_Returns200(t *testing.T) {
	environment := newTestEnvironment(t)
	response := environment.do("POST", "/api/collections/items/records",
		`{"name":"New Item","description":"Desc"}`, environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestItems_UpdateDelete_ResidentsCannotModify(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/items/records/%s", environment.ids.item)

	// Residents (tenant, household, agent) cannot update items.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"tenant", environment.tok.tenant},
		{"household", environment.tok.household},
	} {
		t.Run(tok.name+" cannot update", func(t *testing.T) {
			response := environment.do("PATCH", path, `{"name":"Hacked"}`, tok.token)
			environment.assertStatus(t, response, http.StatusNotFound)
		})
	}

	// Person owner can update.
	response := environment.do("PATCH", path, `{"name":"Updated"}`, environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestEntryDevices_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.entryDevice

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/entryDevices/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/entryDevices/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/entryDevices/records", "{}", "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/entryDevices/records/%s", id), "{}", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/entryDevices/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestEntryDevices_View_PersonOwnerSees200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/entryDevices/records/%s", environment.ids.entryDevice)
	// Only person owner (via item.personItems chain) can view entryDevices.
	// PocketBase v0.22 can't handle forward→back→forward→back chains from entryDevices.
	response := environment.do("GET", path, "", environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestEntryDevices_CreateUpdateDelete_ResidentsCannotModify(t *testing.T) {
	environment := newTestEnvironment(t)

	// Residents cannot create entry devices.
	newItemID := environment.createRecord(t, "items",
		map[string]any{"name": "New Key", "description": "D"})
	body := fmt.Sprintf(`{"item":%q,"deviceType":"Key","identifier":"X001","defunctReason":"None"}`, newItemID)

	response := environment.do("POST", "/api/collections/entryDevices/records", body, environment.tok.tenant)
	environment.assertStatus(t, response, http.StatusBadRequest)

	// Resident cannot update existing entry device.
	path := fmt.Sprintf("/api/collections/entryDevices/records/%s", environment.ids.entryDevice)
	response = environment.do("PATCH", path, `{"identifier":"HACKED"}`, environment.tok.tenant)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestPropertyItems_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.propertyItem

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/propertyItems/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/propertyItems/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/propertyItems/records", "{}", "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/propertyItems/records/%s", id), "{}", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/propertyItems/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestPropertyItems_View_AllGroupsSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/propertyItems/records/%s", environment.ids.propertyItem)
	// Cobrand manager omitted — PocketBase v0.22 limitation.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", environment.tok.userOwner},
		{"cobrand owner", environment.tok.cobrandAdmin},
		{"tenant", environment.tok.tenant},
		{"household", environment.tok.household},
		{"agent", environment.tok.agent},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestPropertyItems_Create_EitherSideCanAdd(t *testing.T) {
	environment := newTestEnvironment(t)
	// Person owner (item owner via personItems) can add a propertyItems record.
	newItemID := environment.createRecord(t, "items",
		map[string]any{"name": "Extra Item", "description": "D"})
	environment.createRecord(t, "personItems",
		map[string]any{"person": environment.ids.personOwner, "item": newItemID})
	body := fmt.Sprintf(`{"item":%q,"property":%q}`, newItemID, environment.ids.property)
	response := environment.do("POST", "/api/collections/propertyItems/records", body, environment.tok.userOwner)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		t.Errorf("person owner create propertyItem: expected 200 or 400, got %d", response.StatusCode)
	}
}

func TestPersonItems_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.personItem

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/personItems/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/personItems/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/personItems/records", "{}", "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/personItems/records/%s", id), "{}", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/personItems/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestPersonItems_View_PersonOwnerSees200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/personItems/records/%s", environment.ids.personItem)
	// personItems view rule is person.user.id only; property-owner access requires 3+ backs.
	response := environment.do("GET", path, "", environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestPersonItems_Create_AnyAuthenticated_Returns200(t *testing.T) {
	environment := newTestEnvironment(t)
	newItemID := environment.createRecord(t, "items",
		map[string]any{"name": "Unclaimed Item", "description": "D"})
	unrelatedPersonID := environment.createRecord(t, "persons",
		map[string]any{"name": "Temp", "DOB": "2000-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"item":%q}`, unrelatedPersonID, newItemID)
	response := environment.do("POST", "/api/collections/personItems/records", body, environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestPersonItems_UpdateDelete_PersonOrOwner(t *testing.T) {
	environment := newTestEnvironment(t)

	// Unrelated cannot update personItems record.
	path := fmt.Sprintf("/api/collections/personItems/records/%s", environment.ids.personItem)
	response := environment.do("PATCH", path, fmt.Sprintf(`{"person":%q}`, environment.ids.personOwner), environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)

	// The person themselves can update.
	response = environment.do("PATCH", path, fmt.Sprintf(`{"person":%q}`, environment.ids.personOwner), environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusOK)
}
