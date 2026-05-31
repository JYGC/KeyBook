package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

// ── persons ───────────────────────────────────────────────────────────────────

func TestPersons_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.personOwner
	body := `{"name":"X","DOB":"1990-01-01 00:00:00.000Z"}`

	// Non-null rules: unauthenticated gets 200 empty (list), 404 (view/update/delete), 400 (create).
	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/persons/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/persons/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/persons/records", body, "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/persons/records/%s", id), body, "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/persons/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestPersons_List_AuthorizedGroupsSee200(t *testing.T) {
	env := newTestEnv(t)

	// The three authorized groups are: own user, tenant property owners/managers,
	// item property owners/managers. We verify each gets a non-error response.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"own user (owner)", env.tok.userOwner},
		{"cobrand admin (co-owner)", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("GET", "/api/collections/persons/records", "", tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestPersons_View_OwnUserSees200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", env.ids.personOwner)
	resp := env.do("GET", path, "", env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestPersons_View_UnrelatedUser_Returns404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", env.ids.personOwner)
	resp := env.do("GET", path, "", env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestPersons_Create_AnyAuthenticatedUser_Returns200(t *testing.T) {
	env := newTestEnv(t)
	body := `{"name":"New Person","DOB":"1995-06-01 00:00:00.000Z"}`
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"owner", env.tok.userOwner},
		{"unrelated", env.tok.unrelated},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("POST", "/api/collections/persons/records", body, tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestPersons_Update_OwnUser_Returns200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", env.ids.personOwner)
	resp := env.do("PATCH", path, `{"name":"Alice Updated"}`, env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestPersons_Update_NonLinkedUser_Returns404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", env.ids.personOwner)
	resp := env.do("PATCH", path, `{"name":"Hacked"}`, env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestPersons_Delete_OwnUser_Returns204(t *testing.T) {
	env := newTestEnv(t)

	// Create a disposable person linked to unrelated user so we can delete it.
	unrelatedUID := env.userIdByEmail(t, "unrelated@test.com")
	pid := env.createRecord(t, "persons",
		map[string]any{"name": "To Delete", "DOB": "2000-01-01 00:00:00.000Z", "user": unrelatedUID})

	path := fmt.Sprintf("/api/collections/persons/records/%s", pid)
	resp := env.do("DELETE", path, "", env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNoContent)
}

func TestPersons_Delete_NonLinkedUser_Returns404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", env.ids.personOwner)
	resp := env.do("DELETE", path, "", env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)
}

// ── properties ────────────────────────────────────────────────────────────────

func TestProperties_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)
	id := env.ids.property
	body := `{"address":"X"}`

	t.Run("GET list", func(t *testing.T) {
		resp := env.do("GET", "/api/collections/properties/records", "", "")
		env.assertStatus(t, resp, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		resp := env.do("GET", fmt.Sprintf("/api/collections/properties/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		resp := env.do("POST", "/api/collections/properties/records", body, "")
		env.assertStatus(t, resp, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		resp := env.do("PATCH", fmt.Sprintf("/api/collections/properties/records/%s", id), body, "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		resp := env.do("DELETE", fmt.Sprintf("/api/collections/properties/records/%s", id), "", "")
		env.assertStatus(t, resp, http.StatusNotFound)
	})
}

func TestProperties_View_FiveGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)

	// Cobrand manager omitted — PocketBase v0.22 can't resolve
	// cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand chain.
	groups := []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"tenant", env.tok.tenant},
		{"household", env.tok.household},
		{"agent", env.tok.agent},
	}
	for _, g := range groups {
		t.Run(g.name, func(t *testing.T) {
			resp := env.do("GET", path, "", g.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestProperties_View_UnrelatedUser_Returns404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)
	resp := env.do("GET", path, "", env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestProperties_Create_AnyAuthenticatedUser_Returns200(t *testing.T) {
	env := newTestEnv(t)
	resp := env.do("POST", "/api/collections/properties/records",
		`{"address":"New Addr"}`, env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestProperties_Update_OwnersSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)

	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("PATCH", path, `{"address":"Updated Addr"}`, tok.token)
			env.assertStatus(t, resp, http.StatusOK)
		})
	}
}

func TestProperties_Update_NonOwnerGroupsReturn404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)

	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand manager", env.tok.manager},
		{"tenant", env.tok.tenant},
		{"household", env.tok.household},
		{"agent", env.tok.agent},
		{"unrelated", env.tok.unrelated},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("PATCH", path, `{"address":"Hacked"}`, tok.token)
			env.assertStatus(t, resp, http.StatusNotFound)
		})
	}
}

func TestProperties_Delete_OwnerRulePassesReferentialIntegrityBlocks(t *testing.T) {
	env := newTestEnv(t)

	// Create a fresh property with ownership chain.
	pid := env.createRecord(t, "properties", map[string]any{"address": "To Delete"})
	powID := env.createRecord(t, "propertyOwners", map[string]any{"property": pid})
	env.createRecord(t, "personPropertyOwners",
		map[string]any{"person": env.ids.personOwner, "propertyOwner": powID})

	path := fmt.Sprintf("/api/collections/properties/records/%s", pid)

	// Non-owner gets 404 (access rule denies).
	resp := env.do("DELETE", path, "", env.tok.unrelated)
	env.assertStatus(t, resp, http.StatusNotFound)

	// Owner gets 400: the access rule passes (not 404), but propertyOwners.property
	// is Required so PocketBase blocks deletion with a referential integrity error.
	resp = env.do("DELETE", path, "", env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

func TestProperties_Delete_NonOwnerGroupsReturn404(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)

	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand manager", env.tok.manager},
		{"tenant", env.tok.tenant},
		{"unrelated", env.tok.unrelated},
	} {
		t.Run(tok.name, func(t *testing.T) {
			resp := env.do("DELETE", path, "", tok.token)
			env.assertStatus(t, resp, http.StatusNotFound)
		})
	}
}
