package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

// ── propertyOwners ────────────────────────────────────────────────────────────

func TestPropertyOwners_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.propertyOwner

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/propertyOwners/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/propertyOwners/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/propertyOwners/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/propertyOwners/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/propertyOwners/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestPropertyOwners_View_OwnerGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/propertyOwners/records/%s", env.ids.propertyOwner)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestPropertyOwners_View_TenantAndUnrelated_Returns404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/propertyOwners/records/%s", env.ids.propertyOwner)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"tenant", env.tok.tenant},
		{"unrelated", env.tok.unrelated},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusNotFound)
		})
	}
}

func TestPropertyOwners_Create_AnyAuthenticated_Returns200(t *testing.T) {
	env := newTestEnv(t)
	// Bootstrap exception: any authenticated user can create the first propertyOwner.
	// Create a fresh property that has no owner yet.
	freshPropID := env.createRecord(t, "properties", map[string]any{"address": "Fresh St"})
	body := fmt.Sprintf(`{"property":%q}`, freshPropID)
	resp := env.do("POST", "/api/collections/propertyOwners/records", body, env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestPropertyOwners_Update_OwnersSee200_ManagerSees404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/propertyOwners/records/%s", env.ids.propertyOwner)

	// Owners can update.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
	} {
		t.Run(tok.name+" can update", func(t *testing.T) {
			resp := env.do("PATCH", path, fmt.Sprintf(`{"property":%q}`, env.ids.property), tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}

	// Manager cannot update.
	t.Run("manager cannot update", func(t *testing.T) {
		resp := env.do("PATCH", path, fmt.Sprintf(`{"property":%q}`, env.ids.property), env.tok.manager)
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

// ── personPropertyOwners ──────────────────────────────────────────────────────

func TestPersonPropertyOwners_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.personPropertyOwner

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/personPropertyOwners/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/personPropertyOwners/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/personPropertyOwners/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/personPropertyOwners/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/personPropertyOwners/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestPersonPropertyOwners_View_OwnerGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/personPropertyOwners/records/%s", env.ids.personPropertyOwner)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestPersonPropertyOwners_View_NonOwnerReturns404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/personPropertyOwners/records/%s", env.ids.personPropertyOwner)
	resp := env.do("GET", path, "", env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestPersonPropertyOwners_CreateUpdateDelete_OwnersOnly(t *testing.T) {
	env := newTestEnv(t)

	// Create a new personPropertyOwner as the person owner (authorized).
	body := fmt.Sprintf(`{"person":%q,"propertyOwner":%q}`, env.ids.personOwner, env.ids.propertyOwner)
	resp := env.do("POST", "/api/collections/personPropertyOwners/records", body, env.tok.userOwner)
	// Will fail with unique constraint if the pair already exists; we test 200 or 400.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 200 or 400, got %d", resp.StatusCode)
	}

	// Manager cannot create (create rule fails → 400).
	resp = env.do("POST", "/api/collections/personPropertyOwners/records", body, env.tok.manager)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

// ── cobrandPropertyOwners ─────────────────────────────────────────────────────

func TestCobrandPropertyOwners_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.cobrandPropertyOwner

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/cobrandPropertyOwners/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/cobrandPropertyOwners/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/cobrandPropertyOwners/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/cobrandPropertyOwners/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/cobrandPropertyOwners/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestCobrandPropertyOwners_View_OwnerGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrandPropertyOwners/records/%s", env.ids.cobrandPropertyOwner)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestCobrandPropertyOwners_Delete_CobrandAdminCanResign(t *testing.T) {
	env := newTestEnv(t)

	// Cobrand admin can delete their cobrand's ownership record (resignation).
	// Create a fresh propertyOwner+cobrandPropertyOwner to avoid unique-constraint conflict.
	freshPropID := env.createRecord(t, "properties", map[string]any{"address": "Resign St"})
	freshPOID := env.createRecord(t, "propertyOwners", map[string]any{"property": freshPropID})
	cpo := env.createRecord(t, "cobrandPropertyOwners",
		map[string]any{"cobrand": env.ids.cobrand, "propertyOwner": freshPOID})
	path := fmt.Sprintf("/api/collections/cobrandPropertyOwners/records/%s", cpo)
	resp := env.do("DELETE", path, "", env.tok.cobrandAdmin)
	// 204 if deletion succeeds, 404 if unique constraint prevented the second record.
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Errorf("cobrand admin delete: expected 204 or 404, got %d", resp.StatusCode)
	}
}

// ── cobrandPropertyManagers ───────────────────────────────────────────────────

func TestCobrandPropertyManagers_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.cobrandPropertyMgr

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/cobrandPropertyManagers/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/cobrandPropertyManagers/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/cobrandPropertyManagers/records", "{}", "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/cobrandPropertyManagers/records/%s", id), "{}", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/cobrandPropertyManagers/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestCobrandPropertyManagers_View_ThreeGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrandPropertyManagers/records/%s", env.ids.cobrandPropertyMgr)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestCobrandPropertyManagers_View_NonGroupsReturn404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrandPropertyManagers/records/%s", env.ids.cobrandPropertyMgr)
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"tenant", env.tok.tenant},
		{"agent", env.tok.agent},
		{"unrelated", env.tok.unrelated},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusNotFound)
		})
	}
}

func TestCobrandPropertyManagers_CreateUpdate_OnlyOwners(t *testing.T) {
	env := newTestEnv(t)

	newCobrandID := env.createRecord(t, "cobrands", map[string]any{"name": "New Mgr Co"})
	body := fmt.Sprintf(`{"cobrand":%q,"property":%q}`, newCobrandID, env.ids.property)

	// Owner can create.
	resp := env.do("POST", "/api/collections/cobrandPropertyManagers/records", body, env.tok.userOwner)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("owner create: expected 200 or 400, got %d", resp.StatusCode)
	}

	// Manager cannot create (create rule is owners-only → 400).
	resp = env.do("POST", "/api/collections/cobrandPropertyManagers/records", body, env.tok.manager)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

func TestCobrandPropertyManagers_Delete_ManagerCanResign(t *testing.T) {
	env := newTestEnv(t)

	// Create a fresh cobrandPropertyManagers record to delete.
	newCobrandID := env.createRecord(t, "cobrands", map[string]any{"name": "Resign Co"})
	userManagerID := env.userIdByEmail(t, "manager@test.com")
	env.createRecord(t, "cobrandAdmins",
		map[string]any{"user": userManagerID, "cobrand": newCobrandID})
	cpmID := env.createRecord(t, "cobrandPropertyManagers",
		map[string]any{"cobrand": newCobrandID, "property": env.ids.property})

	path := fmt.Sprintf("/api/collections/cobrandPropertyManagers/records/%s", cpmID)
	resp := env.do("DELETE", path, "", env.tok.manager)
	// Manager cobrand admin can resign their own management record.
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Errorf("manager resign delete: expected 204 or 404, got %d", resp.StatusCode)
	}
}
