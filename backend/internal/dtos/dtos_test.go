package dtos_test

import (
	"encoding/json"
	"testing"

	"keybook/backend/internal/dtos"
)

func marshalToStringMap(t *testing.T, value any) map[string]string {
	t.Helper()
	marshaledJson, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var marshaledFields map[string]string
	if err := json.Unmarshal(marshaledJson, &marshaledFields); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	return marshaledFields
}

func assertJsonField(t *testing.T, marshaledFields map[string]string, fieldName, want string) {
	t.Helper()
	if got := marshaledFields[fieldName]; got != want {
		t.Errorf("json field %q: want %q, got %q", fieldName, want, got)
	}
}

func TestPersonDto_FieldMapping(t *testing.T) {
	personDto := dtos.PersonDto{
		Id:           "p1",
		Name:         "Alice",
		DOB:          "1990-01-01 00:00:00.000Z",
		User:         "u1",
		ProfileImage: "img.jpg",
	}
	marshaledFields := marshalToStringMap(t, personDto)
	assertJsonField(t, marshaledFields, "id", "p1")
	assertJsonField(t, marshaledFields, "name", "Alice")
	assertJsonField(t, marshaledFields, "DOB", "1990-01-01 00:00:00.000Z")
	assertJsonField(t, marshaledFields, "user", "u1")
	assertJsonField(t, marshaledFields, "profileImage", "img.jpg")
}

func TestPersonDto_Unmarshal(t *testing.T) {
	rawJson := `{"id":"p2","name":"Bob","DOB":"1985-06-15 00:00:00.000Z","user":"u2","profileImage":"bob.png"}`
	var personDto dtos.PersonDto
	if err := json.Unmarshal([]byte(rawJson), &personDto); err != nil {
		t.Fatal(err)
	}
	if personDto.Id != "p2" {
		t.Errorf("Id: want p2, got %s", personDto.Id)
	}
	if personDto.Name != "Bob" {
		t.Errorf("Name: want Bob, got %s", personDto.Name)
	}
	if personDto.DOB != "1985-06-15 00:00:00.000Z" {
		t.Errorf("DOB: want 1985-06-15 00:00:00.000Z, got %s", personDto.DOB)
	}
	if personDto.User != "u2" {
		t.Errorf("User: want u2, got %s", personDto.User)
	}
	if personDto.ProfileImage != "bob.png" {
		t.Errorf("ProfileImage: want bob.png, got %s", personDto.ProfileImage)
	}
}

func TestPersonDto_RequiredFields(t *testing.T) {
	var personDto dtos.PersonDto
	if personDto.Id != "" || personDto.Name != "" || personDto.DOB != "" || personDto.User != "" || personDto.ProfileImage != "" {
		t.Error("zero-value PersonDto should have all empty string fields")
	}
}

func TestPropertyDto_FieldMapping(t *testing.T) {
	propertyDto := dtos.PropertyDto{
		Id:      "prop1",
		Address: "1 Test St",
	}
	marshaledFields := marshalToStringMap(t, propertyDto)
	assertJsonField(t, marshaledFields, "id", "prop1")
	assertJsonField(t, marshaledFields, "address", "1 Test St")
}

func TestPropertyDto_Unmarshal(t *testing.T) {
	rawJson := `{"id":"prop2","address":"2 Sample Ave"}`
	var propertyDto dtos.PropertyDto
	if err := json.Unmarshal([]byte(rawJson), &propertyDto); err != nil {
		t.Fatal(err)
	}
	if propertyDto.Id != "prop2" {
		t.Errorf("Id: want prop2, got %s", propertyDto.Id)
	}
	if propertyDto.Address != "2 Sample Ave" {
		t.Errorf("Address: want 2 Sample Ave, got %s", propertyDto.Address)
	}
}

func TestCobrandDto_FieldMapping(t *testing.T) {
	cobrandDto := dtos.CobrandDto{Id: "cb1", Name: "Acme Co"}
	marshaledFields := marshalToStringMap(t, cobrandDto)
	assertJsonField(t, marshaledFields, "id", "cb1")
	assertJsonField(t, marshaledFields, "name", "Acme Co")
}

func TestCobrandDto_Unmarshal(t *testing.T) {
	rawJson := `{"id":"cb2","name":"Beta LLC"}`
	var cobrandDto dtos.CobrandDto
	if err := json.Unmarshal([]byte(rawJson), &cobrandDto); err != nil {
		t.Fatal(err)
	}
	if cobrandDto.Id != "cb2" || cobrandDto.Name != "Beta LLC" {
		t.Errorf("CobrandDto unmarshal: got {%s %s}", cobrandDto.Id, cobrandDto.Name)
	}
}

// CobrandAdminDto mixes string and bool fields, so it is marshaled to
// map[string]any here rather than reusing the map[string]string helper above.
func TestCobrandAdminDto_FieldMapping(t *testing.T) {
	cobrandAdminDto := dtos.CobrandAdminDto{Id: "ca1", User: "u1", Cobrand: "cb1", Approved: true}
	data, err := json.Marshal(cobrandAdminDto)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var marshaledFields map[string]any
	if err := json.Unmarshal(data, &marshaledFields); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	if marshaledFields["id"] != "ca1" {
		t.Errorf("id: want ca1, got %v", marshaledFields["id"])
	}
	if marshaledFields["user"] != "u1" {
		t.Errorf("user: want u1, got %v", marshaledFields["user"])
	}
	if marshaledFields["cobrand"] != "cb1" {
		t.Errorf("cobrand: want cb1, got %v", marshaledFields["cobrand"])
	}
	if marshaledFields["approved"] != true {
		t.Errorf("approved: want true, got %v", marshaledFields["approved"])
	}
}

func TestCobrandAdminDto_ApprovedDefaultsFalse(t *testing.T) {
	cobrandAdminDto := dtos.CobrandAdminDto{Id: "ca2", User: "u2", Cobrand: "cb2"}
	if cobrandAdminDto.Approved {
		t.Error("Approved: want false zero-value, got true")
	}
}

func TestCobrandPropertyManagerDto_FieldMapping(t *testing.T) {
	cobrandPropertyManagerDto := dtos.CobrandPropertyManagerDto{Id: "cpm1", Cobrand: "cb1", Property: "prop1"}
	marshaledFields := marshalToStringMap(t, cobrandPropertyManagerDto)
	assertJsonField(t, marshaledFields, "id", "cpm1")
	assertJsonField(t, marshaledFields, "cobrand", "cb1")
	assertJsonField(t, marshaledFields, "property", "prop1")
}

func TestCobrandPropertyOwnerDto_FieldMapping(t *testing.T) {
	cobrandPropertyOwnerDto := dtos.CobrandPropertyOwnerDto{Id: "cpo1", Cobrand: "cb1", PropertyOwner: "po1"}
	marshaledFields := marshalToStringMap(t, cobrandPropertyOwnerDto)
	assertJsonField(t, marshaledFields, "id", "cpo1")
	assertJsonField(t, marshaledFields, "cobrand", "cb1")
	assertJsonField(t, marshaledFields, "propertyOwner", "po1")
}

func TestPropertyOwnerDto_FieldMapping(t *testing.T) {
	propertyOwnerDto := dtos.PropertyOwnerDto{Id: "po1", Property: "prop1"}
	marshaledFields := marshalToStringMap(t, propertyOwnerDto)
	assertJsonField(t, marshaledFields, "id", "po1")
	assertJsonField(t, marshaledFields, "property", "prop1")
}

func TestPersonPropertyOwnerDto_FieldMapping(t *testing.T) {
	personPropertyOwnerDto := dtos.PersonPropertyOwnerDto{Id: "ppo1", Person: "p1", PropertyOwner: "po1"}
	marshaledFields := marshalToStringMap(t, personPropertyOwnerDto)
	assertJsonField(t, marshaledFields, "id", "ppo1")
	assertJsonField(t, marshaledFields, "person", "p1")
	assertJsonField(t, marshaledFields, "propertyOwner", "po1")
}

func TestAgentDto_FieldMapping(t *testing.T) {
	agentDto := dtos.AgentDto{Id: "ag1", Person: "p1", Cobrand: "cb1"}
	marshaledFields := marshalToStringMap(t, agentDto)
	assertJsonField(t, marshaledFields, "id", "ag1")
	assertJsonField(t, marshaledFields, "person", "p1")
	assertJsonField(t, marshaledFields, "cobrand", "cb1")
}

func TestPropertyAgentDto_FieldMapping(t *testing.T) {
	propertyAgentDto := dtos.PropertyAgentDto{Id: "pa1", Agent: "ag1", Property: "prop1"}
	marshaledFields := marshalToStringMap(t, propertyAgentDto)
	assertJsonField(t, marshaledFields, "id", "pa1")
	assertJsonField(t, marshaledFields, "agent", "ag1")
	assertJsonField(t, marshaledFields, "property", "prop1")
}

func TestHouseholdDto_FieldMapping(t *testing.T) {
	householdDto := dtos.HouseholdDto{Id: "hh1", Person: "p1", Property: "prop1"}
	marshaledFields := marshalToStringMap(t, householdDto)
	assertJsonField(t, marshaledFields, "id", "hh1")
	assertJsonField(t, marshaledFields, "person", "p1")
	assertJsonField(t, marshaledFields, "property", "prop1")
}

func TestTenantDto_FieldMapping(t *testing.T) {
	tenantDto := dtos.TenantDto{Id: "tn1", Person: "p1", Property: "prop1"}
	marshaledFields := marshalToStringMap(t, tenantDto)
	assertJsonField(t, marshaledFields, "id", "tn1")
	assertJsonField(t, marshaledFields, "person", "p1")
	assertJsonField(t, marshaledFields, "property", "prop1")
}

func TestItemDto_FieldMapping(t *testing.T) {
	itemDto := dtos.ItemDto{Id: "it1", Name: "Key", Description: "Front door key", Picture: "key.jpg"}
	marshaledFields := marshalToStringMap(t, itemDto)
	assertJsonField(t, marshaledFields, "id", "it1")
	assertJsonField(t, marshaledFields, "name", "Key")
	assertJsonField(t, marshaledFields, "description", "Front door key")
	assertJsonField(t, marshaledFields, "picture", "key.jpg")
}

func TestItemDto_Unmarshal(t *testing.T) {
	rawJson := `{"id":"it2","name":"Fob","description":"Gate fob","picture":"fob.png"}`
	var itemDto dtos.ItemDto
	if err := json.Unmarshal([]byte(rawJson), &itemDto); err != nil {
		t.Fatal(err)
	}
	if itemDto.Id != "it2" || itemDto.Name != "Fob" || itemDto.Description != "Gate fob" || itemDto.Picture != "fob.png" {
		t.Errorf("ItemDto unmarshal: got {%s %s %s %s}", itemDto.Id, itemDto.Name, itemDto.Description, itemDto.Picture)
	}
}

func TestPropertyItemDto_FieldMapping(t *testing.T) {
	propertyItemDto := dtos.PropertyItemDto{Id: "pi1", Item: "it1", Property: "prop1"}
	marshaledFields := marshalToStringMap(t, propertyItemDto)
	assertJsonField(t, marshaledFields, "id", "pi1")
	assertJsonField(t, marshaledFields, "item", "it1")
	assertJsonField(t, marshaledFields, "property", "prop1")
}

func TestPersonItemDto_FieldMapping(t *testing.T) {
	personItemDto := dtos.PersonItemDto{Id: "pei1", Person: "p1", Item: "it1"}
	marshaledFields := marshalToStringMap(t, personItemDto)
	assertJsonField(t, marshaledFields, "id", "pei1")
	assertJsonField(t, marshaledFields, "person", "p1")
	assertJsonField(t, marshaledFields, "item", "it1")
}

func TestEntryDeviceDto_FieldMapping(t *testing.T) {
	entryDeviceDto := dtos.EntryDeviceDto{
		Id:            "ed1",
		Item:          "it1",
		DeviceType:    "Key",
		Identifier:    "KEY-001",
		DefunctReason: "None",
	}
	marshaledFields := marshalToStringMap(t, entryDeviceDto)
	assertJsonField(t, marshaledFields, "id", "ed1")
	assertJsonField(t, marshaledFields, "item", "it1")
	assertJsonField(t, marshaledFields, "deviceType", "Key")
	assertJsonField(t, marshaledFields, "identifier", "KEY-001")
	assertJsonField(t, marshaledFields, "defunctReason", "None")
}

func TestEntryDeviceDto_Unmarshal(t *testing.T) {
	rawJson := `{"id":"ed2","item":"it2","deviceType":"Fob","identifier":"FOB-002","defunctReason":"Lost"}`
	var entryDeviceDto dtos.EntryDeviceDto
	if err := json.Unmarshal([]byte(rawJson), &entryDeviceDto); err != nil {
		t.Fatal(err)
	}
	if entryDeviceDto.Id != "ed2" || entryDeviceDto.Item != "it2" || entryDeviceDto.DeviceType != "Fob" ||
		entryDeviceDto.Identifier != "FOB-002" || entryDeviceDto.DefunctReason != "Lost" {
		t.Errorf("EntryDeviceDto unmarshal got unexpected values")
	}
}
