package contract_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"

	_ "keybook/backend/migrations"
)

// tokens holds auth tokens for each user type in the fixture set.
type tokens struct {
	admin           string
	userOwner       string
	cobrandAdmin    string
	unapprovedAdmin string
	manager         string
	tenant          string
	household       string
	agent           string
	unrelated       string
}

// ids holds record IDs for all created fixture records.
type ids struct {
	property               string
	propertyOwner          string
	personOwner            string
	personPropertyOwner    string
	cobrand                string
	cobrandAdmin           string
	unapprovedCobrandAdmin string
	cobrandPropertyOwner   string
	cobrandPropertyMgr     string
	personTenant           string
	personHousehold        string
	personAgent            string
	tenant                 string
	household              string
	agent                  string
	propertyAgent          string
	item                   string
	personItem             string
	propertyItem           string
	entryDevice            string
}

type testEnvironment struct {
	e   *echo.Echo
	tok tokens
	ids ids
	app *tests.TestApp
}

// newTestEnvironment boots a fresh PocketBase instance with our migrations applied,
// creates one admin and seven regular users, then builds the full standard
// fixture set described in design.md §7.
func newTestEnvironment(t *testing.T) *testEnvironment {
	t.Helper()

	testApp, err := tests.NewTestApp(t.TempDir())
	if err != nil {
		t.Fatalf("new test app: %v", err)
	}
	t.Cleanup(testApp.Cleanup)

	// Create PocketBase admin directly in DB.
	admin := &models.Admin{Email: "admin@test.com"}
	admin.SetPassword("Admin12345678!")
	if err := testApp.Dao().SaveAdmin(admin); err != nil {
		t.Fatalf("save admin: %v", err)
	}

	e, err := apis.InitApi(testApp)
	if err != nil {
		t.Fatalf("init api: %v", err)
	}
	testApp.OnBeforeServe().Trigger(&core.ServeEvent{App: testApp, Router: e}) //nolint:errcheck

	environment := &testEnvironment{e: e, app: testApp}

	environment.tok.admin = environment.adminToken(t, "admin@test.com", "Admin12345678!")

	for _, u := range []struct {
		email  string
		pw     string
		target *string
	}{
		{"owner@test.com", "Pass12345678!", &environment.tok.userOwner},
		{"cobrand@test.com", "Pass12345678!", &environment.tok.cobrandAdmin},
		{"unapproved@test.com", "Pass12345678!", &environment.tok.unapprovedAdmin},
		{"manager@test.com", "Pass12345678!", &environment.tok.manager},
		{"tenant@test.com", "Pass12345678!", &environment.tok.tenant},
		{"household@test.com", "Pass12345678!", &environment.tok.household},
		{"agent@test.com", "Pass12345678!", &environment.tok.agent},
		{"unrelated@test.com", "Pass12345678!", &environment.tok.unrelated},
	} {
		environment.createUser(t, u.email, u.pw)
		*u.target = environment.userToken(t, u.email, u.pw)
	}

	userOwnerID := environment.userIdByEmail(t, "owner@test.com")
	userCobrandAdminID := environment.userIdByEmail(t, "cobrand@test.com")
	userUnapprovedAdminID := environment.userIdByEmail(t, "unapproved@test.com")
	userManagerID := environment.userIdByEmail(t, "manager@test.com")
	userTenantID := environment.userIdByEmail(t, "tenant@test.com")
	userHouseholdID := environment.userIdByEmail(t, "household@test.com")
	userAgentID := environment.userIdByEmail(t, "agent@test.com")

	environment.ids.personOwner = environment.createRecord(t, "persons",
		map[string]any{"name": "Alice Owner", "DOB": "1990-01-01 00:00:00.000Z", "user": userOwnerID})

	environment.ids.personTenant = environment.createRecord(t, "persons",
		map[string]any{"name": "Bob Tenant", "DOB": "1992-05-15 00:00:00.000Z", "user": userTenantID})

	environment.ids.personHousehold = environment.createRecord(t, "persons",
		map[string]any{"name": "Carol Household", "DOB": "1988-03-20 00:00:00.000Z", "user": userHouseholdID})

	environment.ids.personAgent = environment.createRecord(t, "persons",
		map[string]any{"name": "Dave Agent", "DOB": "1985-07-10 00:00:00.000Z", "user": userAgentID})

	environment.ids.property = environment.createRecord(t, "properties",
		map[string]any{"address": "1 Test Street"})

	environment.ids.propertyOwner = environment.createRecord(t, "propertyOwners",
		map[string]any{"property": environment.ids.property})

	environment.ids.personPropertyOwner = environment.createRecord(t, "personPropertyOwners",
		map[string]any{"person": environment.ids.personOwner, "propertyOwner": environment.ids.propertyOwner})

	environment.ids.cobrand = environment.createRecord(t, "cobrands",
		map[string]any{"name": "Test Co"})

	environment.ids.cobrandAdmin = environment.createRecord(t, "cobrandAdmins",
		map[string]any{"user": userCobrandAdminID, "cobrand": environment.ids.cobrand, "approved": true})

	environment.ids.unapprovedCobrandAdmin = environment.createRecord(t, "cobrandAdmins",
		map[string]any{"user": userUnapprovedAdminID, "cobrand": environment.ids.cobrand, "approved": false})

	environment.ids.cobrandPropertyOwner = environment.createRecord(t, "cobrandPropertyOwners",
		map[string]any{"cobrand": environment.ids.cobrand, "propertyOwner": environment.ids.propertyOwner})

	environment.ids.cobrandPropertyMgr = environment.createRecord(t, "cobrandPropertyManagers",
		map[string]any{"cobrand": environment.ids.cobrand, "property": environment.ids.property})

	// We need a second cobrand for the manager user; instead, use a separate
	// cobrand that manages the property via cobrandPropertyManagers. The manager
	// user (userManagerID) needs to be a cobrand admin of a cobrand that has a
	// cobrandPropertyManagers record for our property. Override cobrandPropertyMgr
	// to use a new "manager cobrand" with userManagerID as admin.
	managerCobrandID := environment.createRecord(t, "cobrands",
		map[string]any{"name": "Manager Co"})
	environment.createRecord(t, "cobrandAdmins",
		map[string]any{"user": userManagerID, "cobrand": managerCobrandID})
	// Replace the previous cobrandPropertyMgr with one for the manager cobrand.
	environment.ids.cobrandPropertyMgr = environment.createRecord(t, "cobrandPropertyManagers",
		map[string]any{"cobrand": managerCobrandID, "property": environment.ids.property})

	environment.ids.tenant = environment.createRecord(t, "tenants",
		map[string]any{"person": environment.ids.personTenant, "property": environment.ids.property})

	environment.ids.household = environment.createRecord(t, "households",
		map[string]any{"person": environment.ids.personHousehold, "property": environment.ids.property})

	environment.ids.agent = environment.createRecord(t, "agents",
		map[string]any{"person": environment.ids.personAgent, "cobrand": environment.ids.cobrand})

	environment.ids.propertyAgent = environment.createRecord(t, "propertyAgents",
		map[string]any{"agent": environment.ids.agent, "property": environment.ids.property})

	environment.ids.item = environment.createRecord(t, "items",
		map[string]any{"name": "Front Door Key", "description": "Main key"})

	environment.ids.personItem = environment.createRecord(t, "personItems",
		map[string]any{"person": environment.ids.personOwner, "item": environment.ids.item})

	environment.ids.propertyItem = environment.createRecord(t, "propertyItems",
		map[string]any{"item": environment.ids.item, "property": environment.ids.property})

	environment.ids.entryDevice = environment.createRecord(t, "entryDevices",
		map[string]any{
			"item":          environment.ids.item,
			"deviceType":    "Key",
			"identifier":    "KEY-001",
			"defunctReason": "None",
		})

	return environment
}

func (environment *testEnvironment) do(method, path, body, token string) *http.Response {
	var bodyReader *strings.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	} else {
		bodyReader = strings.NewReader("")
	}
	request := httptest.NewRequest(method, path, bodyReader)
	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	responseRecorder := httptest.NewRecorder()
	environment.e.ServeHTTP(responseRecorder, request)
	return responseRecorder.Result()
}

func (environment *testEnvironment) assertStatus(t *testing.T, got *http.Response, want int) {
	t.Helper()
	if got.StatusCode != want {
		body := new(bytes.Buffer)
		body.ReadFrom(got.Body) //nolint:errcheck
		t.Errorf("expected status %d, got %d — body: %s", want, got.StatusCode, body.String())
	}
}

// adminToken authenticates as a PocketBase admin and returns the JWT.
func (environment *testEnvironment) adminToken(t *testing.T, email, pw string) string {
	t.Helper()
	response := environment.do("POST", "/api/admins/auth-with-password",
		fmt.Sprintf(`{"identity":%q,"password":%q}`, email, pw), "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("admin auth failed: %d", response.StatusCode)
	}
	return extractToken(t, response)
}

// userToken authenticates as a `users` collection record and returns the JWT.
func (environment *testEnvironment) userToken(t *testing.T, email, pw string) string {
	t.Helper()
	response := environment.do("POST", "/api/collections/users/auth-with-password",
		fmt.Sprintf(`{"identity":%q,"password":%q}`, email, pw), "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("user auth failed for %s: %d", email, response.StatusCode)
	}
	return extractToken(t, response)
}

// createUser creates a users record via the admin API.
func (environment *testEnvironment) createUser(t *testing.T, email, pw string) {
	t.Helper()
	body := fmt.Sprintf(`{"email":%q,"password":%q,"passwordConfirm":%q,"emailVisibility":true}`,
		email, pw, pw)
	response := environment.do("POST", "/api/collections/users/records", body, environment.tok.admin)
	if response.StatusCode != http.StatusOK {
		requestBodyBuffer := new(bytes.Buffer)
		requestBodyBuffer.ReadFrom(response.Body) //nolint:errcheck
		t.Fatalf("create user %s failed: %d — %s", email, response.StatusCode, requestBodyBuffer.String())
	}
}

// userIdByEmail looks up a users record ID by email via admin API.
func (environment *testEnvironment) userIdByEmail(t *testing.T, email string) string {
	t.Helper()
	response := environment.do("GET",
		fmt.Sprintf("/api/collections/users/records?filter=(email=%q)", email),
		"", environment.tok.admin)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("lookup user %s failed: %d", email, response.StatusCode)
	}
	var result struct {
		Items []struct {
			Id string `json:"id"`
		} `json:"items"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil || len(result.Items) == 0 {
		t.Fatalf("lookup user %s: no result", email)
	}
	return result.Items[0].Id
}

// createRecord creates a record in the given collection via the admin API and
// returns its ID.
func (environment *testEnvironment) createRecord(t *testing.T, collection string, fields map[string]any) string {
	t.Helper()
	body, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal %s body: %v", collection, err)
	}
	response := environment.do("POST", "/api/collections/"+collection+"/records", string(body), environment.tok.admin)
	if response.StatusCode != http.StatusOK {
		requestBodyBuffer := new(bytes.Buffer)
		requestBodyBuffer.ReadFrom(response.Body) //nolint:errcheck
		t.Fatalf("create %s failed: %d — %s", collection, response.StatusCode, requestBodyBuffer.String())
	}
	var result struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil || result.Id == "" {
		t.Fatalf("create %s: no id in response", collection)
	}
	return result.Id
}

// extractToken pulls the "token" field from an auth response body.
func extractToken(t *testing.T, response *http.Response) string {
	t.Helper()
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil || result.Token == "" {
		t.Fatalf("extract token: %v", err)
	}
	return result.Token
}
