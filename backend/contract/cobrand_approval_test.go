package contract_test

import (
	"fmt"
	"net/http"
	"testing"
)

// These tests cover the cobrandAdmins.approved gate (design.md section 2.4).
// environment.ids.cobrand has two admins: userCobrandAdmin (approved = true, the
// standard fixture) and userUnapprovedAdmin (approved = false). Both admin
// the same cobrand, which co-owns environment.ids.property via cobrandPropertyOwner —
// this lets each test prove the gate checks the ACTING admin's own approved
// value rather than merely "is any admin of this cobrand approved."

func TestCobrandApproval_UnapprovedAdminCanStillViewOwnCobrand(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", environment.ids.cobrand)
	response := environment.do("GET", path, "", environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusOK)
}

// ── cobrandAdmins.update is locked to superusers ─────────────────────────────

func TestCobrandApproval_NoAuthenticatedUserCanUpdateCobrandAdmins(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", environment.ids.cobrandAdmin)
	// Even the approved admin whose own record this is cannot patch it — the
	// update rule is locked to superusers only (nil), which PocketBase
	// rejects with 403 rather than the usual 404 "not visible" response.
	response := environment.do("PATCH", path, `{"approved":true}`, environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusForbidden)
}

// Note: "inviting another admin requires the inviter to be an approved
// existing admin" is enforced by the OnRecordBeforeCreateRequest("cobrandAdmins")
// hook in cmd/keybook.go, calling CobrandService.EnsureInviterIsApprovedAdmin.
// This contract test harness (newTestEnvironment) builds its PocketBase app via
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
	environment := newTestEnvironment(t)
	// Unapproved admin attempts to remove the approved admin's record.
	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", environment.ids.cobrandAdmin)
	response := environment.do("DELETE", path, "", environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestCobrandApproval_ApprovedAdminCanRemoveAnotherAdmin(t *testing.T) {
	environment := newTestEnvironment(t)
	response := environment.do("DELETE",
		fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", environment.ids.unapprovedCobrandAdmin),
		"", environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusNoContent)
}

func TestCobrandApproval_UnapprovedAdminCanStillResignSelf(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrandAdmins/records/%s", environment.ids.unapprovedCobrandAdmin)
	response := environment.do("DELETE", path, "", environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusNoContent)
}

func TestCobrandApproval_UnapprovedAdminCannotUpdateCobrand(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", environment.ids.cobrand)
	response := environment.do("PATCH", path, `{"name":"Hacked"}`, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestCobrandApproval_ApprovedAdminCanStillUpdateCobrand(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/cobrands/records/%s", environment.ids.cobrand)
	response := environment.do("PATCH", path, `{"name":"Renamed Co"}`, environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusOK)
}

// ── properties update/delete: the acting admin's OWN approval must gate this,
// not a co-admin's — environment.ids.cobrand has one approved and one unapproved admin ──

func TestCobrandApproval_UnapprovedAdminCannotUpdateProperty(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", environment.ids.property)
	response := environment.do("PATCH", path, `{"address":"Hacked St"}`, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusNotFound)
}

func TestCobrandApproval_ApprovedAdminCanStillUpdateProperty(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/properties/records/%s", environment.ids.property)
	response := environment.do("PATCH", path, `{"address":"1 Renamed Street"}`, environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestCobrandApproval_UnapprovedAdminCannotUpdatePropertyOwner(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/propertyOwners/records/%s", environment.ids.propertyOwner)
	response := environment.do("PATCH", path, fmt.Sprintf(`{"property":%q}`, environment.ids.property), environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusNotFound)
}

// ── cobrandPropertyManagers create: only the owning cobrand's APPROVED admin ─

func TestCobrandApproval_UnapprovedAdminCannotAppointPropertyManager(t *testing.T) {
	environment := newTestEnvironment(t)
	newCobrandID := environment.createRecord(t, "cobrands", map[string]any{"name": "Candidate Mgr Co"})
	body := fmt.Sprintf(`{"cobrand":%q,"property":%q}`, newCobrandID, environment.ids.property)
	response := environment.do("POST", "/api/collections/cobrandPropertyManagers/records", body, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusBadRequest)
}

func TestCobrandApproval_ApprovedAdminCanStillAppointPropertyManager(t *testing.T) {
	environment := newTestEnvironment(t)
	newCobrandID := environment.createRecord(t, "cobrands", map[string]any{"name": "Candidate Mgr Co 2"})
	body := fmt.Sprintf(`{"cobrand":%q,"property":%q}`, newCobrandID, environment.ids.property)
	response := environment.do("POST", "/api/collections/cobrandPropertyManagers/records", body, environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusOK)
}

// ── cobrandPropertyOwners create is gated; delete (resignation) is not ──────

func TestCobrandApproval_UnapprovedAdminCannotClaimPropertyOwnership(t *testing.T) {
	environment := newTestEnvironment(t)
	freshPropertyId := environment.createRecord(t, "properties", map[string]any{"address": "Claim St"})
	freshPropertyOwnerId := environment.createRecord(t, "propertyOwners", map[string]any{"property": freshPropertyId})
	body := fmt.Sprintf(`{"cobrand":%q,"propertyOwner":%q}`, environment.ids.cobrand, freshPropertyOwnerId)
	response := environment.do("POST", "/api/collections/cobrandPropertyOwners/records", body, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusBadRequest)
}

func TestCobrandApproval_UnapprovedAdminCanStillResignPropertyOwnership(t *testing.T) {
	environment := newTestEnvironment(t)
	// Seed the cobrandPropertyOwners record directly via the admin API — the
	// gate only applies to the create path, not to resigning an existing one.
	freshPropertyId := environment.createRecord(t, "properties", map[string]any{"address": "Resign St 2"})
	freshPropertyOwnerId := environment.createRecord(t, "propertyOwners", map[string]any{"property": freshPropertyId})
	cobrandPropertyOwnerId := environment.createRecord(t, "cobrandPropertyOwners",
		map[string]any{"cobrand": environment.ids.cobrand, "propertyOwner": freshPropertyOwnerId})
	path := fmt.Sprintf("/api/collections/cobrandPropertyOwners/records/%s", cobrandPropertyOwnerId)
	response := environment.do("DELETE", path, "", environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusNoContent)
}

// ── agents create/update, and delete's admin-removes-agent clause ───────────

func TestCobrandApproval_UnapprovedAdminCannotRegisterAgent(t *testing.T) {
	environment := newTestEnvironment(t)
	newPersonID := environment.createRecord(t, "persons", map[string]any{"name": "Candidate Agent", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"cobrand":%q}`, newPersonID, environment.ids.cobrand)
	response := environment.do("POST", "/api/collections/agents/records", body, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusBadRequest)
}

func TestCobrandApproval_ApprovedAdminCanStillRegisterAgent(t *testing.T) {
	environment := newTestEnvironment(t)
	newPersonID := environment.createRecord(t, "persons", map[string]any{"name": "Candidate Agent 2", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"cobrand":%q}`, newPersonID, environment.ids.cobrand)
	response := environment.do("POST", "/api/collections/agents/records", body, environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestCobrandApproval_UnapprovedAdminCannotRemoveAgent(t *testing.T) {
	environment := newTestEnvironment(t)
	path := fmt.Sprintf("/api/collections/agents/records/%s", environment.ids.agent)
	response := environment.do("DELETE", path, "", environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusNotFound)
}

// ── propertyAgents create: the agent's cobrand admin must be approved ───────

func TestCobrandApproval_UnapprovedAdminCannotAssignAgentToProperty(t *testing.T) {
	environment := newTestEnvironment(t)
	newPersonID := environment.createRecord(t, "persons", map[string]any{"name": "Fresh Agent Person", "DOB": "1990-01-01 00:00:00.000Z"})
	newAgentID := environment.createRecord(t, "agents", map[string]any{"person": newPersonID, "cobrand": environment.ids.cobrand})
	body := fmt.Sprintf(`{"agent":%q,"property":%q}`, newAgentID, environment.ids.property)
	response := environment.do("POST", "/api/collections/propertyAgents/records", body, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusBadRequest)
}

// ── tenants / households create is gated by the managing/owning cobrand's
// approval ────────────────────────────────────────────────────────────────────

func TestCobrandApproval_UnapprovedAdminCannotAddTenant(t *testing.T) {
	environment := newTestEnvironment(t)
	newPersonID := environment.createRecord(t, "persons", map[string]any{"name": "Candidate Tenant", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"property":%q}`, newPersonID, environment.ids.property)
	response := environment.do("POST", "/api/collections/tenants/records", body, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusBadRequest)
}

func TestCobrandApproval_ApprovedAdminCanStillAddTenant(t *testing.T) {
	environment := newTestEnvironment(t)
	newPersonID := environment.createRecord(t, "persons", map[string]any{"name": "Candidate Tenant 2", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"property":%q}`, newPersonID, environment.ids.property)
	response := environment.do("POST", "/api/collections/tenants/records", body, environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestCobrandApproval_UnapprovedAdminCannotAddHouseholdMember(t *testing.T) {
	environment := newTestEnvironment(t)
	newPersonID := environment.createRecord(t, "persons", map[string]any{"name": "Candidate Household", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"property":%q}`, newPersonID, environment.ids.property)
	response := environment.do("POST", "/api/collections/households/records", body, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusBadRequest)
}

func TestCobrandApproval_UnapprovedAdminCannotAddPropertyItem(t *testing.T) {
	environment := newTestEnvironment(t)
	newItemID := environment.createRecord(t, "items", map[string]any{"name": "Candidate Item", "description": "D"})
	body := fmt.Sprintf(`{"item":%q,"property":%q}`, newItemID, environment.ids.property)
	response := environment.do("POST", "/api/collections/propertyItems/records", body, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusBadRequest)
}

func TestCobrandApproval_ApprovedAdminCanStillAddPropertyItem(t *testing.T) {
	environment := newTestEnvironment(t)
	newItemID := environment.createRecord(t, "items", map[string]any{"name": "Candidate Item 2", "description": "D"})
	body := fmt.Sprintf(`{"item":%q,"property":%q}`, newItemID, environment.ids.property)
	response := environment.do("POST", "/api/collections/propertyItems/records", body, environment.tok.cobrandAdmin)
	environment.assertStatus(t, response, http.StatusOK)
}

func TestCobrandApproval_UnapprovedAdminCannotLinkPersonOwner(t *testing.T) {
	environment := newTestEnvironment(t)
	newPersonID := environment.createRecord(t, "persons", map[string]any{"name": "Candidate Person Owner", "DOB": "1990-01-01 00:00:00.000Z"})
	body := fmt.Sprintf(`{"person":%q,"propertyOwner":%q}`, newPersonID, environment.ids.propertyOwner)
	response := environment.do("POST", "/api/collections/personPropertyOwners/records", body, environment.tok.unapprovedAdmin)
	environment.assertStatus(t, response, http.StatusBadRequest)
}
