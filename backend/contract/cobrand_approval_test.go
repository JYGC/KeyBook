package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

// These tests cover the cobrandAdmins.approved gate (design.md section 2.4).
// env.ids.cobrand has two admins: userCobrandAdmin (approved = true, the
// standard fixture) and userUnapprovedAdmin (approved = false). Both admin
// the same cobrand, which co-owns env.ids.property via cobrandPropertyOwner —
// this lets each test prove the gate checks the ACTING admin's own approved
// value rather than merely "is any admin of this cobrand approved."

// ── list/view unaffected by approval ─────────────────────────────────────────

func TestCobrandApproval_UnapprovedAdminCanStillViewOwnCobrand(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)
	resp := env.do("GET", path, "", env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusOK)
}

// ── cobrandAdmins.update is locked to superusers ─────────────────────────────

func TestCobrandApproval_NoAuthenticatedUserCanUpdateCobrandAdmins(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", env.ids.cobrandAdmin)
	// Even the approved admin whose own record this is cannot patch it — the
	// update rule is locked to superusers only (nil), which PocketBase
	// rejects with 403 rather than the usual 404 "not visible" response.
	resp := env.do("PATCH", path, `{"approved":true}`, env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusForbidden)
}

// Note: "inviting another admin requires the inviter to be an approved
// existing admin" is enforced by the OnRecordBeforeCreateRequest("cobrandAdmins")
// hook in cmd/keybook.go, calling CobrandService.EnsureInviterIsApprovedAdmin.
// This contract test harness (newTestEnv) builds its PocketBase app via
// tests.NewTestApp + apis.InitApi directly, the same way every other contract
// test file in this package does — it never calls cmd.startBackend(), so none
// of keybook.go's custom hooks are registered here (this is true of the
// pre-existing Phase 1 hooks too, e.g. cobrandAdmins' own duplicate-admin
// check). That behavior is instead covered end-to-end by
// TestCobrandApplicationService_InviterMustBeApprovedAdmin in
// internal/application/cobrand_app_test.go, which exercises
// CobrandApplicationService.AddCobrandAdmin directly.

// ── cobrandAdmins delete: removing ANOTHER admin requires approval; removing
// your OWN record (resignation) does not ─────────────────────────────────────

func TestCobrandApproval_UnapprovedAdminCannotRemoveAnotherAdmin(t *testing.T) {
	env := newTestEnv(t)
	// Unapproved admin attempts to remove the approved admin's record.
	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", env.ids.cobrandAdmin)
	resp := env.do("DELETE", path, "", env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestCobrandApproval_ApprovedAdminCanRemoveAnotherAdmin(t *testing.T) {
	env := newTestEnv(t)
	resp := env.do("DELETE",
		fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", env.ids.unapprovedCobrandAdmin),
		"", env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusNoContent)
}

func TestCobrandApproval_UnapprovedAdminCanStillResignSelf(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", env.ids.unapprovedCobrandAdmin)
	resp := env.do("DELETE", path, "", env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusNoContent)
}

// ── cobrands update/delete ───────────────────────────────────────────────────

func TestCobrandApproval_UnapprovedAdminCannotUpdateCobrand(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)
	resp := env.do("PATCH", path, `{"name":"Hacked"}`, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestCobrandApproval_ApprovedAdminCanStillUpdateCobrand(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", env.ids.cobrand)
	resp := env.do("PATCH", path, `{"name":"Renamed Co"}`, env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusOK)
}

// ── properties update/delete: the acting admin's OWN approval must gate this,
// not a co-admin's — env.ids.cobrand has one approved and one unapproved admin ──

func TestCobrandApproval_UnapprovedAdminCannotUpdateProperty(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)
	resp := env.do("PATCH", path, `{"address":"Hacked St"}`, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusNotFound)
}

func TestCobrandApproval_ApprovedAdminCanStillUpdateProperty(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", env.ids.property)
	resp := env.do("PATCH", path, `{"address":"1 Renamed Street"}`, env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusOK)
}

// ── propertyOwners update/delete ─────────────────────────────────────────────

func TestCobrandApproval_UnapprovedAdminCannotUpdatePropertyOwner(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/propertyOwners/records/%s", env.ids.propertyOwner)
	resp := env.do("PATCH", path, fmt.Sprintf(`{"property":%q}`, env.ids.property), env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusNotFound)
}

// ── cobrandPropertyManagers create: only the owning cobrand's APPROVED admin ─

func TestCobrandApproval_UnapprovedAdminCannotAppointPropertyManager(t *testing.T) {
	env := newTestEnv(t)
	newCobrandID := env.createRecord(t, "cobrands", map[string]any{"name": "Candidate Mgr Co"})
	body := fmt.Sprintf(`{"cobrand":%q,"property":%q}`, newCobrandID, env.ids.property)
	resp := env.do("POST", "/api/collections/cobrandPropertyManagers/records", body, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

func TestCobrandApproval_ApprovedAdminCanStillAppointPropertyManager(t *testing.T) {
	env := newTestEnv(t)
	newCobrandID := env.createRecord(t, "cobrands", map[string]any{"name": "Candidate Mgr Co 2"})
	body := fmt.Sprintf(`{"cobrand":%q,"property":%q}`, newCobrandID, env.ids.property)
	resp := env.do("POST", "/api/collections/cobrandPropertyManagers/records", body, env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusOK)
}

// ── cobrandPropertyOwners create is gated; delete (resignation) is not ──────

func TestCobrandApproval_UnapprovedAdminCannotClaimPropertyOwnership(t *testing.T) {
	env := newTestEnv(t)
	freshPropID := env.createRecord(t, "properties", map[string]any{"address": "Claim St"})
	freshPOID := env.createRecord(t, "propertyOwners", map[string]any{"property": freshPropID})
	body := fmt.Sprintf(`{"cobrand":%q,"propertyOwner":%q}`, env.ids.cobrand, freshPOID)
	resp := env.do("POST", "/api/collections/cobrandPropertyOwners/records", body, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

func TestCobrandApproval_UnapprovedAdminCanStillResignPropertyOwnership(t *testing.T) {
	env := newTestEnv(t)
	// Seed the cobrandPropertyOwners record directly via the admin API — the
	// gate only applies to the create path, not to resigning an existing one.
	freshPropID := env.createRecord(t, "properties", map[string]any{"address": "Resign St 2"})
	freshPOID := env.createRecord(t, "propertyOwners", map[string]any{"property": freshPropID})
	cpoID := env.createRecord(t, "cobrandPropertyOwners",
		map[string]any{"cobrand": env.ids.cobrand, "propertyOwner": freshPOID})
	path := fmt.Sprintf("/api/collections/cobrandPropertyOwners/records/%s", cpoID)
	resp := env.do("DELETE", path, "", env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusNoContent)
}

// ── agents create/update, and delete's admin-removes-agent clause ───────────

func TestCobrandApproval_UnapprovedAdminCannotRegisterAgent(t *testing.T) {
	env := newTestEnv(t)
	newPersonID := env.createRecord(t, "persons", map[string]any{"name": "Candidate Agent", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"cobrand":%q}`, newPersonID, env.ids.cobrand)
	resp := env.do("POST", "/api/collections/agents/records", body, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

func TestCobrandApproval_ApprovedAdminCanStillRegisterAgent(t *testing.T) {
	env := newTestEnv(t)
	newPersonID := env.createRecord(t, "persons", map[string]any{"name": "Candidate Agent 2", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"cobrand":%q}`, newPersonID, env.ids.cobrand)
	resp := env.do("POST", "/api/collections/agents/records", body, env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestCobrandApproval_UnapprovedAdminCannotRemoveAgent(t *testing.T) {
	env := newTestEnv(t)
	path := fmt.Sprintf("/api/collections/agents/records/%s", env.ids.agent)
	resp := env.do("DELETE", path, "", env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusNotFound)
}

// ── propertyAgents create: the agent's cobrand admin must be approved ───────

func TestCobrandApproval_UnapprovedAdminCannotAssignAgentToProperty(t *testing.T) {
	env := newTestEnv(t)
	newPersonID := env.createRecord(t, "persons", map[string]any{"name": "Fresh Agent Person", "DOB": "1990-01-01 00:00:00.000Z"})
	newAgentID := env.createRecord(t, "agents", map[string]any{"person": newPersonID, "cobrand": env.ids.cobrand})
	body := fmt.Sprintf(`{"agent":%q,"property":%q}`, newAgentID, env.ids.property)
	resp := env.do("POST", "/api/collections/propertyAgents/records", body, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

// ── tenants / households create is gated by the managing/owning cobrand's
// approval ────────────────────────────────────────────────────────────────────

func TestCobrandApproval_UnapprovedAdminCannotAddTenant(t *testing.T) {
	env := newTestEnv(t)
	newPersonID := env.createRecord(t, "persons", map[string]any{"name": "Candidate Tenant", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"property":%q}`, newPersonID, env.ids.property)
	resp := env.do("POST", "/api/collections/tenants/records", body, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

func TestCobrandApproval_ApprovedAdminCanStillAddTenant(t *testing.T) {
	env := newTestEnv(t)
	newPersonID := env.createRecord(t, "persons", map[string]any{"name": "Candidate Tenant 2", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"property":%q}`, newPersonID, env.ids.property)
	resp := env.do("POST", "/api/collections/tenants/records", body, env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusOK)
}

func TestCobrandApproval_UnapprovedAdminCannotAddHouseholdMember(t *testing.T) {
	env := newTestEnv(t)
	newPersonID := env.createRecord(t, "persons", map[string]any{"name": "Candidate Household", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"property":%q}`, newPersonID, env.ids.property)
	resp := env.do("POST", "/api/collections/households/records", body, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

// ── propertyItems create is gated ────────────────────────────────────────────

func TestCobrandApproval_UnapprovedAdminCannotAddPropertyItem(t *testing.T) {
	env := newTestEnv(t)
	newItemID := env.createRecord(t, "items", map[string]any{"name": "Candidate Item", "description": "D"})
	body := fmt.Sprintf(`{"item":%q,"property":%q}`, newItemID, env.ids.property)
	resp := env.do("POST", "/api/collections/propertyItems/records", body, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusBadRequest)
}

func TestCobrandApproval_ApprovedAdminCanStillAddPropertyItem(t *testing.T) {
	env := newTestEnv(t)
	newItemID := env.createRecord(t, "items", map[string]any{"name": "Candidate Item 2", "description": "D"})
	body := fmt.Sprintf(`{"item":%q,"property":%q}`, newItemID, env.ids.property)
	resp := env.do("POST", "/api/collections/propertyItems/records", body, env.tok.cobrandAdmin)
	env.assertStatus(t, resp, http.StatusOK)
}

// ── personPropertyOwners create is gated ─────────────────────────────────────

func TestCobrandApproval_UnapprovedAdminCannotLinkPersonOwner(t *testing.T) {
	env := newTestEnv(t)
	newPersonID := env.createRecord(t, "persons", map[string]any{"name": "Candidate Person Owner", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"propertyOwner":%q}`, newPersonID, env.ids.propertyOwner)
	resp := env.do("POST", "/api/collections/personPropertyOwners/records", body, env.tok.unapprovedAdmin)
	env.assertStatus(t, resp, http.StatusBadRequest)
}
