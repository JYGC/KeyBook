package dtos_test

import (
	"encoding/json"
	"testing"

	"keybook/backend/internal/dtos"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func marshal(t *testing.T, v any) map[string]string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	return m
}

func assertField(t *testing.T, m map[string]string, key, want string) {
	t.Helper()
	if got := m[key]; got != want {
		t.Errorf("json field %q: want %q, got %q", key, want, got)
	}
}

// ── PersonDto ─────────────────────────────────────────────────────────────────

func TestPersonDto_FieldMapping(t *testing.T) {
	dto := dtos.PersonDto{
		Id:           "p1",
		Name:         "Alice",
		DOB:          "1990-01-01 00:00:00.000Z",
		User:         "u1",
		ProfileImage: "img.jpg",
	}
	m := marshal(t, dto)
	assertField(t, m, "id", "p1")
	assertField(t, m, "name", "Alice")
	assertField(t, m, "DOB", "1990-01-01 00:00:00.000Z")
	assertField(t, m, "user", "u1")
	assertField(t, m, "profileImage", "img.jpg")
}

func TestPersonDto_Unmarshal(t *testing.T) {
	raw := `{"id":"p2","name":"Bob","DOB":"1985-06-15 00:00:00.000Z","user":"u2","profileImage":"bob.png"}`
	var dto dtos.PersonDto
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		t.Fatal(err)
	}
	if dto.Id != "p2" {
		t.Errorf("Id: want p2, got %s", dto.Id)
	}
	if dto.Name != "Bob" {
		t.Errorf("Name: want Bob, got %s", dto.Name)
	}
	if dto.DOB != "1985-06-15 00:00:00.000Z" {
		t.Errorf("DOB: want 1985-06-15 00:00:00.000Z, got %s", dto.DOB)
	}
	if dto.User != "u2" {
		t.Errorf("User: want u2, got %s", dto.User)
	}
	if dto.ProfileImage != "bob.png" {
		t.Errorf("ProfileImage: want bob.png, got %s", dto.ProfileImage)
	}
}

func TestPersonDto_RequiredFields(t *testing.T) {
	var dto dtos.PersonDto
	if dto.Id != "" || dto.Name != "" || dto.DOB != "" || dto.User != "" || dto.ProfileImage != "" {
		t.Error("zero-value PersonDto should have all empty string fields")
	}
}

// ── PropertyDto ───────────────────────────────────────────────────────────────

func TestPropertyDto_FieldMapping(t *testing.T) {
	dto := dtos.PropertyDto{
		Id:      "prop1",
		Address: "1 Test St",
	}
	m := marshal(t, dto)
	assertField(t, m, "id", "prop1")
	assertField(t, m, "address", "1 Test St")
}

func TestPropertyDto_Unmarshal(t *testing.T) {
	raw := `{"id":"prop2","address":"2 Sample Ave"}`
	var dto dtos.PropertyDto
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		t.Fatal(err)
	}
	if dto.Id != "prop2" {
		t.Errorf("Id: want prop2, got %s", dto.Id)
	}
	if dto.Address != "2 Sample Ave" {
		t.Errorf("Address: want 2 Sample Ave, got %s", dto.Address)
	}
}

// ── CobrandDto ────────────────────────────────────────────────────────────────

func TestCobrandDto_FieldMapping(t *testing.T) {
	dto := dtos.CobrandDto{Id: "cb1", Name: "Acme Co"}
	m := marshal(t, dto)
	assertField(t, m, "id", "cb1")
	assertField(t, m, "name", "Acme Co")
}

func TestCobrandDto_Unmarshal(t *testing.T) {
	raw := `{"id":"cb2","name":"Beta LLC"}`
	var dto dtos.CobrandDto
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		t.Fatal(err)
	}
	if dto.Id != "cb2" || dto.Name != "Beta LLC" {
		t.Errorf("CobrandDto unmarshal: got {%s %s}", dto.Id, dto.Name)
	}
}

// ── CobrandAdminDto ───────────────────────────────────────────────────────────

// CobrandAdminDto mixes string and bool fields, so it is marshaled to
// map[string]any here rather than reusing the map[string]string helper above.
func TestCobrandAdminDto_FieldMapping(t *testing.T) {
	dto := dtos.CobrandAdminDto{Id: "ca1", User: "u1", Cobrand: "cb1", Approved: true}
	data, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	if m["id"] != "ca1" {
		t.Errorf("id: want ca1, got %v", m["id"])
	}
	if m["user"] != "u1" {
		t.Errorf("user: want u1, got %v", m["user"])
	}
	if m["cobrand"] != "cb1" {
		t.Errorf("cobrand: want cb1, got %v", m["cobrand"])
	}
	if m["approved"] != true {
		t.Errorf("approved: want true, got %v", m["approved"])
	}
}

func TestCobrandAdminDto_ApprovedDefaultsFalse(t *testing.T) {
	dto := dtos.CobrandAdminDto{Id: "ca2", User: "u2", Cobrand: "cb2"}
	if dto.Approved {
		t.Error("Approved: want false zero-value, got true")
	}
}

// ── CobrandPropertyManagerDto ─────────────────────────────────────────────────

func TestCobrandPropertyManagerDto_FieldMapping(t *testing.T) {
	dto := dtos.CobrandPropertyManagerDto{Id: "cpm1", Cobrand: "cb1", Property: "prop1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "cpm1")
	assertField(t, m, "cobrand", "cb1")
	assertField(t, m, "property", "prop1")
}

// ── CobrandPropertyOwnerDto ───────────────────────────────────────────────────

func TestCobrandPropertyOwnerDto_FieldMapping(t *testing.T) {
	dto := dtos.CobrandPropertyOwnerDto{Id: "cpo1", Cobrand: "cb1", PropertyOwner: "po1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "cpo1")
	assertField(t, m, "cobrand", "cb1")
	assertField(t, m, "propertyOwner", "po1")
}

// ── PropertyOwnerDto ──────────────────────────────────────────────────────────

func TestPropertyOwnerDto_FieldMapping(t *testing.T) {
	dto := dtos.PropertyOwnerDto{Id: "po1", Property: "prop1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "po1")
	assertField(t, m, "property", "prop1")
}

// ── PersonPropertyOwnerDto ────────────────────────────────────────────────────

func TestPersonPropertyOwnerDto_FieldMapping(t *testing.T) {
	dto := dtos.PersonPropertyOwnerDto{Id: "ppo1", Person: "p1", PropertyOwner: "po1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "ppo1")
	assertField(t, m, "person", "p1")
	assertField(t, m, "propertyOwner", "po1")
}

// ── AgentDto ──────────────────────────────────────────────────────────────────

func TestAgentDto_FieldMapping(t *testing.T) {
	dto := dtos.AgentDto{Id: "ag1", Person: "p1", Cobrand: "cb1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "ag1")
	assertField(t, m, "person", "p1")
	assertField(t, m, "cobrand", "cb1")
}

// ── PropertyAgentDto ──────────────────────────────────────────────────────────

func TestPropertyAgentDto_FieldMapping(t *testing.T) {
	dto := dtos.PropertyAgentDto{Id: "pa1", Agent: "ag1", Property: "prop1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "pa1")
	assertField(t, m, "agent", "ag1")
	assertField(t, m, "property", "prop1")
}

// ── HouseholdDto ──────────────────────────────────────────────────────────────

func TestHouseholdDto_FieldMapping(t *testing.T) {
	dto := dtos.HouseholdDto{Id: "hh1", Person: "p1", Property: "prop1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "hh1")
	assertField(t, m, "person", "p1")
	assertField(t, m, "property", "prop1")
}

// ── TenantDto ─────────────────────────────────────────────────────────────────

func TestTenantDto_FieldMapping(t *testing.T) {
	dto := dtos.TenantDto{Id: "tn1", Person: "p1", Property: "prop1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "tn1")
	assertField(t, m, "person", "p1")
	assertField(t, m, "property", "prop1")
}

// ── ItemDto ───────────────────────────────────────────────────────────────────

func TestItemDto_FieldMapping(t *testing.T) {
	dto := dtos.ItemDto{Id: "it1", Name: "Key", Description: "Front door key", Picture: "key.jpg"}
	m := marshal(t, dto)
	assertField(t, m, "id", "it1")
	assertField(t, m, "name", "Key")
	assertField(t, m, "description", "Front door key")
	assertField(t, m, "picture", "key.jpg")
}

func TestItemDto_Unmarshal(t *testing.T) {
	raw := `{"id":"it2","name":"Fob","description":"Gate fob","picture":"fob.png"}`
	var dto dtos.ItemDto
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		t.Fatal(err)
	}
	if dto.Id != "it2" || dto.Name != "Fob" || dto.Description != "Gate fob" || dto.Picture != "fob.png" {
		t.Errorf("ItemDto unmarshal: got {%s %s %s %s}", dto.Id, dto.Name, dto.Description, dto.Picture)
	}
}

// ── PropertyItemDto ───────────────────────────────────────────────────────────

func TestPropertyItemDto_FieldMapping(t *testing.T) {
	dto := dtos.PropertyItemDto{Id: "pi1", Item: "it1", Property: "prop1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "pi1")
	assertField(t, m, "item", "it1")
	assertField(t, m, "property", "prop1")
}

// ── PersonItemDto ─────────────────────────────────────────────────────────────

func TestPersonItemDto_FieldMapping(t *testing.T) {
	dto := dtos.PersonItemDto{Id: "pei1", Person: "p1", Item: "it1"}
	m := marshal(t, dto)
	assertField(t, m, "id", "pei1")
	assertField(t, m, "person", "p1")
	assertField(t, m, "item", "it1")
}

// ── EntryDeviceDto ────────────────────────────────────────────────────────────

func TestEntryDeviceDto_FieldMapping(t *testing.T) {
	dto := dtos.EntryDeviceDto{
		Id:            "ed1",
		Item:          "it1",
		DeviceType:    "Key",
		Identifier:    "KEY-001",
		DefunctReason: "None",
	}
	m := marshal(t, dto)
	assertField(t, m, "id", "ed1")
	assertField(t, m, "item", "it1")
	assertField(t, m, "deviceType", "Key")
	assertField(t, m, "identifier", "KEY-001")
	assertField(t, m, "defunctReason", "None")
}

func TestEntryDeviceDto_Unmarshal(t *testing.T) {
	raw := `{"id":"ed2","item":"it2","deviceType":"Fob","identifier":"FOB-002","defunctReason":"Lost"}`
	var dto dtos.EntryDeviceDto
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		t.Fatal(err)
	}
	if dto.Id != "ed2" || dto.Item != "it2" || dto.DeviceType != "Fob" ||
		dto.Identifier != "FOB-002" || dto.DefunctReason != "Lost" {
		t.Errorf("EntryDeviceDto unmarshal got unexpected values")
	}
}
