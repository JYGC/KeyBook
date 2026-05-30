package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

// ── persons ───────────────────────────────────────────────────────────────────

func TestPersons_Unauthenticated_Returns403(t *testing.T) {
	env := newTestEnv(t)

	cases := []struct{ method, path string }{
		{"GET", "/api/collections/persons/records"},
		{"GET", fmt.Sprintf("/api/collections/persons/records/%s", env.ids.personOwner)},
		{"POST", "/api/collections/persons/records"},
		{"PATCH", fmt.Sprintf("/api/collections/persons/records/%s", env.ids.personOwner)},
		{"DELETE", fmt.Sprintf("/api/collections/persons/records/%s", env.ids.personOwner)},
	}
	for _, c := range cases {
		t.Run(c.method+" "+c.path, func(t *testing.T) {
			resp := env.do(c.method, c.path, `{"name":"X","DOB":"1990-01-01 00:00:00.000Z"}`, "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
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

	cases := []struct{ method, path string }{
		{"GET", "/api/collections/properties/records"},
		{"GET", fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)},
		{"POST", "/api/collections/properties/records"},
		{"PATCH", fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)},
		{"DELETE", fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)},
	}
	for _, c := range cases {
		t.Run(c.method, func(t *testing.T) {
			resp := env.do(c.method, c.path, `{"address":"X"}`, "")
			env.assertStatus(t, resp, http.StatusForbidden)
		})
	}
}

func TestProperties_View_AllSixGroupsSee200(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)

	groups := []struct {
		name  string
		token string
	}{
		{"person owner", env.tok.userOwner},
		{"cobrand owner", env.tok.cobrandAdmin},
		{"cobrand manager", env.tok.manager},
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

func TestProperties_Delete_OwnersReturn204(t *testing.T) {
	env := newTestEnv(t)

	// Create a fresh property so deletion doesn't affect other tests.
	pid := env.createRecord(t, "properties", map[string]any{"address": "To Delete"})
	// Create ownership chain so person owner can delete it.
	powID := env.createRecord(t, "propertyOwners", map[string]any{"property": pid})
	env.createRecord(t, "personPropertyOwners",
		map[string]any{"person": env.ids.personOwner, "propertyOwner": powID})

	path := fmt.Sprintf("/api/collections/properties/records/%s", pid)
	resp := env.do("DELETE", path, "", env.tok.userOwner)
	env.assertStatus(t, resp, http.StatusNoContent)
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
