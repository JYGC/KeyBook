package services_test

import (
	"testing"

	"keybook/backend/internal/dtos"
	"keybook/backend/internal/services"
)

func TestPropertyOwnerService_EnsureNoDuplicatePersonOwner(t *testing.T) {
	svc := services.NewPropertyOwnerService()

	existing := []dtos.PersonPropertyOwnerDto{
		{Id: "ppo1", Person: "person1", PropertyOwner: "po1"},
	}

	tests := []struct {
		name     string
		owners   []dtos.PersonPropertyOwnerDto
		personId string
		wantErr  bool
	}{
		{"first person owner", []dtos.PersonPropertyOwnerDto{}, "person1", false},
		{"different person", existing, "person2", false},
		{"duplicate person owner", existing, "person1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.EnsureNoDuplicatePersonOwner(tt.owners, tt.personId)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureNoDuplicatePersonOwner(_, %q) error = %v, wantErr %v", tt.personId, err, tt.wantErr)
			}
		})
	}
}

func TestPropertyOwnerService_EnsureNoDuplicateCobrandOwner(t *testing.T) {
	svc := services.NewPropertyOwnerService()

	existing := []dtos.CobrandPropertyOwnerDto{
		{Id: "cpo1", Cobrand: "cobrand1", PropertyOwner: "po1"},
	}

	tests := []struct {
		name      string
		owners    []dtos.CobrandPropertyOwnerDto
		cobrandId string
		wantErr   bool
	}{
		{"first cobrand owner", []dtos.CobrandPropertyOwnerDto{}, "cobrand1", false},
		{"different cobrand", existing, "cobrand2", false},
		{"duplicate cobrand owner", existing, "cobrand1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.EnsureNoDuplicateCobrandOwner(tt.owners, tt.cobrandId)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureNoDuplicateCobrandOwner(_, %q) error = %v, wantErr %v", tt.cobrandId, err, tt.wantErr)
			}
		})
	}
}
