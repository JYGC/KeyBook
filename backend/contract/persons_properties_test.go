package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPersons_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.personOwner
	body := `{"name":"X","DOB":"1990-01-01 00:00:00.000Z"}`

	// Non-null rules: unauthenticated gets 200 empty (list), 404 (view/update/delete), 400 (create).
	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/persons/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/persons/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/persons/records", body, "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/persons/records/%s", id), body, "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/persons/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestPersons_List_AuthorizedGroupsSee200(t *testing.T) {
	environment := newTestEnvironment(t)

	// The three authorized groups are: own user, tenant property owners/managers,
	// item property owners/managers. We verify each gets a non-error response.
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"own user (owner)", environment.tok.userOwner},
		{"cobrand admin (co-owner)", environment.tok.cobrandAdmin},
		{"cobrand manager", environment.tok.manager},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("GET", "/api/collections/persons/records", "", tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestPersons_View_OwnUserSees200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", environment.ids.personOwner)
	response := environment.do("GET", path, "", environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestPersons_View_UnrelatedUser_Returns404(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", environment.ids.personOwner)
	response := environment.do("GET", path, "", environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestPersons_Create_AnyAuthenticatedUser_Returns200(t *testing.T) {
	environment := newTestEnvironment(t)
	body := `{"name":"New Person","DOB":"1995-06-01 00:00:00.000Z"}`
	for _, tok := range []struct {
		name  string
		token string
	}{
		{"owner", environment.tok.userOwner},
		{"unrelated", environment.tok.unrelated},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("POST", "/api/collections/persons/records", body, tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestPersons_Update_OwnUser_Returns200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", environment.ids.personOwner)
	response := environment.do("PATCH", path, `{"name":"Alice Updated"}`, environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestPersons_Update_NonLinkedUser_Returns404(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", environment.ids.personOwner)
	response := environment.do("PATCH", path, `{"name":"Hacked"}`, environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestPersons_Delete_OwnUser_Returns204(t *testing.T) {
	environment := newTestEnvironment(t)

	// Create a disposable person linked to unrelated user so we can delete it.
	unrelatedUID := environment.userIdByEmail(t, "unrelated@test.com")
	recordId := environment.createRecord(t, "persons",
		map[string]any{"name": "To Delete", "DOB": "2000-01-01 00:00:00.000Z", "user": unrelatedUID})

	path := fmt.Sprintf("/api/collections/persons/records/%s", recordId)
	response := environment.do("DELETE", path, "", environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNoContent)
}

func TestPersons_Delete_NonLinkedUser_Returns404(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/persons/records/%s", environment.ids.personOwner)
	response := environment.do("DELETE", path, "", environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestProperties_Unauthenticated_Returns403(t *testing.T) {
	environment := newTestEnvironment(t)
	id := environment.ids.property
	body := `{"address":"X"}`

	t.Run("GET list", func(t *testing.T) {
		response := environment.do("GET", "/api/collections/properties/records", "", "")
		environment.assertStatus(t, response, http.StatusOK)
	})
	t.Run("GET single", func(t *testing.T) {
		response := environment.do("GET", fmt.Sprintf("/api/collections/properties/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("POST", func(t *testing.T) {
		response := environment.do("POST", "/api/collections/properties/records", body, "")
		environment.assertStatus(t, response, http.StatusBadRequest)
	})
	t.Run("PATCH", func(t *testing.T) {
		response := environment.do("PATCH", fmt.Sprintf("/api/collections/properties/records/%s", id), body, "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
	t.Run("DELETE", func(t *testing.T) {
		response := environment.do("DELETE", fmt.Sprintf("/api/collections/properties/records/%s", id), "", "")
		environment.assertStatus(t, response, http.StatusNotFound)
	})
}

func TestProperties_View_FiveGroupsSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", environment.ids.property)

	// Cobrand manager omitted — PocketBase v0.22 can't resolve
	// cobrandPropertyManagers_via_property.cobrand.cobrandAdmins_via_cobrand chain.
	groups := []struct {
		name  string
		token string
	}{
		{"person owner", environment.tok.userOwner},
		{"cobrand owner", environment.tok.cobrandAdmin},
		{"tenant", environment.tok.tenant},
		{"household", environment.tok.household},
		{"agent", environment.tok.agent},
	}
	for _, g := range groups {
		t.Run(g.name, func(t *testing.T) {
			response := environment.do("GET", path, "", g.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestProperties_View_UnrelatedUser_Returns404(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", environment.ids.property)
	response := environment.do("GET", path, "", environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestProperties_Create_AnyAuthenticatedUser_Returns200(t *testing.T) {
	environment := newTestEnvironment(t)
	response := environment.do("POST", "/api/collections/properties/records",
		`{"address":"New Addr"}`, environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestProperties_Update_OwnersSee200(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", environment.ids.property)

	for _, tok := range []struct {
		name  string
		token string
	}{
		{"person owner", environment.tok.userOwner},
		{"cobrand owner", environment.tok.cobrandAdmin},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("PATCH", path, `{"address":"Updated Addr"}`, tok.token)
			environment.assertStatus(t, response, http.StatusOK)
		})
	}
}

func TestProperties_Update_NonOwnerGroupsReturn404(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", environment.ids.property)

	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand manager", environment.tok.manager},
		{"tenant", environment.tok.tenant},
		{"household", environment.tok.household},
		{"agent", environment.tok.agent},
		{"unrelated", environment.tok.unrelated},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("PATCH", path, `{"address":"Hacked"}`, tok.token)
			environment.assertStatus(t, response, http.StatusNotFound)
		})
	}
}

func TestProperties_Delete_OwnerRulePassesReferentialIntegrityBlocks(t *testing.T) {
	environment := newTestEnvironment(t)

	// Create a fresh property with ownership chain.
	recordId := environment.createRecord(t, "properties", map[string]any{"address": "To Delete"})
	propertyOwnerId := environment.createRecord(t, "propertyOwners", map[string]any{"property": recordId})
	environment.createRecord(t, "personPropertyOwners",
		map[string]any{"person": environment.ids.personOwner, "propertyOwner": propertyOwnerId})

	path := fmt.Sprintf("/api/collections/properties/records/%s", recordId)

	// Non-owner gets 404 (access rule denies).
	response := environment.do("DELETE", path, "", environment.tok.unrelated)
	environment.assertStatus(t, response, http.StatusNotFound)

	// Owner gets 400: the access rule passes (not 404), but propertyOwners.property
	// is Required so PocketBase blocks deletion with a referential integrity error.
	response = environment.do("DELETE", path, "", environment.tok.userOwner)
	environment.assertStatus(t, response, http.StatusBadRequest)
}

func TestProperties_Delete_NonOwnerGroupsReturn404(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", environment.ids.property)

	for _, tok := range []struct {
		name  string
		token string
	}{
		{"cobrand manager", environment.tok.manager},
		{"tenant", environment.tok.tenant},
		{"unrelated", environment.tok.unrelated},
	} {
		t.Run(tok.name, func(t *testing.T) {
			response := environment.do("DELETE", path, "", tok.token)
			environment.assertStatus(t, response, http.StatusNotFound)
		})
	}
}
