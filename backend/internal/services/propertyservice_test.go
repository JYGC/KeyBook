package services_test

import (
	"testing"

	"keybook/backend/internal/dtos"
	"keybook/backend/internal/services"
)

func TestPropertyService_ValidateProperty(t *testing.T) {
	svc := services.NewPropertyService()

	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{"valid", "1 Main St", false},
		{"empty address", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateProperty(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProperty(%q) error = %v, wantErr %v", tt.address, err, tt.wantErr)
			}
		})
	}
}

func TestPropertyService_IsPersonOwner(t *testing.T) {
	svc := services.NewPropertyService()

	owners := []dtos.PersonPropertyOwnerDto{
		{Id: "ppo1", Person: "person1", PropertyOwner: "po1"},
		{Id: "ppo2", Person: "person2", PropertyOwner: "po1"},
	}

	tests := []struct {
		name     string
		personId string
		want     bool
	}{
		{"owner present", "person1", true},
		{"another owner present", "person2", true},
		{"not an owner", "person3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.IsPersonOwner(owners, tt.personId)
			if got != tt.want {
				t.Errorf("IsPersonOwner(_, %q) = %v, want %v", tt.personId, got, tt.want)
			}
		})
	}
}

func TestPropertyService_IsCobrandOwner(t *testing.T) {
	svc := services.NewPropertyService()

	owners := []dtos.CobrandPropertyOwnerDto{
		{Id: "cpo1", Cobrand: "cobrand1", PropertyOwner: "po1"},
	}

	tests := []struct {
		name      string
		cobrandId string
		want      bool
	}{
		{"owner present", "cobrand1", true},
		{"not an owner", "cobrand2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.IsCobrandOwner(owners, tt.cobrandId)
			if got != tt.want {
				t.Errorf("IsCobrandOwner(_, %q) = %v, want %v", tt.cobrandId, got, tt.want)
			}
		})
	}
}
