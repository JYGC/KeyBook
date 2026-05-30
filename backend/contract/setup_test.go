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
	admin        string
	userOwner    string
	cobrandAdmin string
	manager      string
	tenant       string
	household    string
	agent        string
	unrelated    string
}

// ids holds record IDs for all created fixture records.
type ids struct {
	property              string
	propertyOwner         string
	personOwner           string
	personPropertyOwner   string
	cobrand               string
	cobrandAdmin          string
	cobrandPropertyOwner  string
	cobrandPropertyMgr    string
	personTenant          string
	personHousehold       string
	personAgent           string
	tenant                string
	household             string
	agent                 string
	propertyAgent         string
	item                  string
	personItem            string
	propertyItem          string
	entryDevice           string
}

type testEnv struct {
	e      *echo.Echo
	tok    tokens
	ids    ids
	app    *tests.TestApp
}

// newTestEnv boots a fresh PocketBase instance with our migrations applied,
// creates one admin and seven regular users, then builds the full standard
// fixture set described in design.md §7.
func newTestEnv(t *testing.T) *testEnv {
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

	env := &testEnv{e: e, app: testApp}

	// ── Auth tokens ──────────────────────────────────────────────────────────

	env.tok.admin = env.adminToken(t, "admin@test.com", "Admin12345678!")

	for _, u := range []struct {
		email  string
		pw     string
		target *string
	}{
		{"owner@test.com", "Pass12345678!", &env.tok.userOwner},
		{"cobrand@test.com", "Pass12345678!", &env.tok.cobrandAdmin},
		{"manager@test.com", "Pass12345678!", &env.tok.manager},
		{"tenant@test.com", "Pass12345678!", &env.tok.tenant},
		{"household@test.com", "Pass12345678!", &env.tok.household},
		{"agent@test.com", "Pass12345678!", &env.tok.agent},
		{"unrelated@test.com", "Pass12345678!", &env.tok.unrelated},
	} {
		env.createUser(t, u.email, u.pw)
		*u.target = env.userToken(t, u.email, u.pw)
	}

	// ── Look up user IDs ─────────────────────────────────────────────────────

	userOwnerID := env.userIdByEmail(t, "owner@test.com")
	userCobrandAdminID := env.userIdByEmail(t, "cobrand@test.com")
	userManagerID := env.userIdByEmail(t, "manager@test.com")
	userTenantID := env.userIdByEmail(t, "tenant@test.com")
	userHouseholdID := env.userIdByEmail(t, "household@test.com")
	userAgentID := env.userIdByEmail(t, "agent@test.com")

	// ── Persons ──────────────────────────────────────────────────────────────

	env.ids.personOwner = env.createRecord(t, "persons",
		map[string]any{"name": "Alice Owner", "DOB": "1990-01-01 00:00:00.000Z", "user": userOwnerID})

	env.ids.personTenant = env.createRecord(t, "persons",
		map[string]any{"name": "Bob Tenant", "DOB": "1992-05-15 00:00:00.000Z", "user": userTenantID})

	env.ids.personHousehold = env.createRecord(t, "persons",
		map[string]any{"name": "Carol Household", "DOB": "1988-03-20 00:00:00.000Z", "user": userHouseholdID})

	env.ids.personAgent = env.createRecord(t, "persons",
		map[string]any{"name": "Dave Agent", "DOB": "1985-07-10 00:00:00.000Z", "user": userAgentID})

	// ── Properties ───────────────────────────────────────────────────────────

	env.ids.property = env.createRecord(t, "properties",
		map[string]any{"address": "1 Test Street"})

	// ── Ownership chain ──────────────────────────────────────────────────────

	env.ids.propertyOwner = env.createRecord(t, "propertyOwners",
		map[string]any{"property": env.ids.property})

	env.ids.personPropertyOwner = env.createRecord(t, "personPropertyOwners",
		map[string]any{"person": env.ids.personOwner, "propertyOwner": env.ids.propertyOwner})

	// ── Cobrands ─────────────────────────────────────────────────────────────

	env.ids.cobrand = env.createRecord(t, "cobrands",
		map[string]any{"name": "Test Co"})

	env.ids.cobrandAdmin = env.createRecord(t, "cobrandAdmins",
		map[string]any{"user": userCobrandAdminID, "cobrand": env.ids.cobrand})

	env.ids.cobrandPropertyOwner = env.createRecord(t, "cobrandPropertyOwners",
		map[string]any{"cobrand": env.ids.cobrand, "propertyOwner": env.ids.propertyOwner})

	env.ids.cobrandPropertyMgr = env.createRecord(t, "cobrandPropertyManagers",
		map[string]any{"cobrand": env.ids.cobrand, "property": env.ids.property})

	// We need a second cobrand for the manager user; instead, use a separate
	// cobrand that manages the property via cobrandPropertyManagers. The manager
	// user (userManagerID) needs to be a cobrand admin of a cobrand that has a
	// cobrandPropertyManagers record for our property. Override cobrandPropertyMgr
	// to use a new "manager cobrand" with userManagerID as admin.
	managerCobrandID := env.createRecord(t, "cobrands",
		map[string]any{"name": "Manager Co"})
	env.createRecord(t, "cobrandAdmins",
		map[string]any{"user": userManagerID, "cobrand": managerCobrandID})
	// Replace the previous cobrandPropertyMgr with one for the manager cobrand.
	env.ids.cobrandPropertyMgr = env.createRecord(t, "cobrandPropertyManagers",
		map[string]any{"cobrand": managerCobrandID, "property": env.ids.property})

	// ── Residents ────────────────────────────────────────────────────────────

	env.ids.tenant = env.createRecord(t, "tenants",
		map[string]any{"person": env.ids.personTenant, "property": env.ids.property})

	env.ids.household = env.createRecord(t, "households",
		map[string]any{"person": env.ids.personHousehold, "property": env.ids.property})

	// ── Agents ───────────────────────────────────────────────────────────────

	env.ids.agent = env.createRecord(t, "agents",
		map[string]any{"person": env.ids.personAgent, "cobrand": env.ids.cobrand})

	env.ids.propertyAgent = env.createRecord(t, "propertyAgents",
		map[string]any{"agent": env.ids.agent, "property": env.ids.property})

	// ── Items ─────────────────────────────────────────────────────────────────

	env.ids.item = env.createRecord(t, "items",
		map[string]any{"name": "Front Door Key", "description": "Main key"})

	env.ids.personItem = env.createRecord(t, "personItems",
		map[string]any{"person": env.ids.personOwner, "item": env.ids.item})

	env.ids.propertyItem = env.createRecord(t, "propertyItems",
		map[string]any{"item": env.ids.item, "property": env.ids.property})

	env.ids.entryDevice = env.createRecord(t, "entryDevices",
		map[string]any{
			"item":          env.ids.item,
			"deviceType":    "Key",
			"identifier":    "KEY-001",
			"defunctReason": "None",
		})

	return env
}

// ── HTTP helpers ─────────────────────────────────────────────────────────────

func (env *testEnv) do(method, path, body, token string) *http.Response {
	var bodyReader *strings.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	} else {
		bodyReader = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	env.e.ServeHTTP(rr, req)
	return rr.Result()
}

func (env *testEnv) assertStatus(t *testing.T, got *http.Response, want int) {
	t.Helper()
	if got.StatusCode != want {
		body := new(bytes.Buffer)
		body.ReadFrom(got.Body) //nolint:errcheck
		t.Errorf("expected status %d, got %d — body: %s", want, got.StatusCode, body.String())
	}
}

// adminToken authenticates as a PocketBase admin and returns the JWT.
func (env *testEnv) adminToken(t *testing.T, email, pw string) string {
	t.Helper()
	resp := env.do("POST", "/api/admins/auth-with-password",
		fmt.Sprintf(`{"identity":%q,"password":%q}`, email, pw), "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("admin auth failed: %d", resp.StatusCode)
	}
	return extractToken(t, resp)
}

// userToken authenticates as a `users` collection record and returns the JWT.
func (env *testEnv) userToken(t *testing.T, email, pw string) string {
	t.Helper()
	resp := env.do("POST", "/api/collections/users/auth-with-password",
		fmt.Sprintf(`{"identity":%q,"password":%q}`, email, pw), "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("user auth failed for %s: %d", email, resp.StatusCode)
	}
	return extractToken(t, resp)
}

// createUser creates a users record via the admin API.
func (env *testEnv) createUser(t *testing.T, email, pw string) {
	t.Helper()
	body := fmt.Sprintf(`{"email":%q,"password":%q,"passwordConfirm":%q,"emailVisibility":true}`,
		email, pw, pw)
	resp := env.do("POST", "/api/collections/users/records", body, env.tok.admin)
	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body) //nolint:errcheck
		t.Fatalf("create user %s failed: %d — %s", email, resp.StatusCode, buf.String())
	}
}

// userIdByEmail looks up a users record ID by email via admin API.
func (env *testEnv) userIdByEmail(t *testing.T, email string) string {
	t.Helper()
	resp := env.do("GET",
		fmt.Sprintf("/api/collections/users/records?filter=(email=%q)", email),
		"", env.tok.admin)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("lookup user %s failed: %d", email, resp.StatusCode)
	}
	var result struct {
		Items []struct {
			Id string `json:"id"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result.Items) == 0 {
		t.Fatalf("lookup user %s: no result", email)
	}
	return result.Items[0].Id
}

// createRecord creates a record in the given collection via the admin API and
// returns its ID.
func (env *testEnv) createRecord(t *testing.T, collection string, fields map[string]any) string {
	t.Helper()
	body, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal %s body: %v", collection, err)
	}
	resp := env.do("POST", "/api/collections/"+collection+"/records", string(body), env.tok.admin)
	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body) //nolint:errcheck
		t.Fatalf("create %s failed: %d — %s", collection, resp.StatusCode, buf.String())
	}
	var result struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.Id == "" {
		t.Fatalf("create %s: no id in response", collection)
	}
	return result.Id
}

// extractToken pulls the "token" field from an auth response body.
func extractToken(t *testing.T, resp *http.Response) string {
	t.Helper()
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.Token == "" {
		t.Fatalf("extract token: %v", err)
	}
	return result.Token
}
