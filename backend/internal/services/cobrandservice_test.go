package services_test

import (
	"testing"

	"keybook/backend/internal/dtos"
	"keybook/backend/internal/services"
)

func TestCobrandService_ValidateCobrand(t *testing.T) {
	svc := services.NewCobrandService()

	tests := []struct {
		name    string
		cName   string
		wantErr bool
	}{
		{"valid", "Acme Co", false},
		{"empty name", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateCobrand(tt.cName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCobrand(%q) error = %v, wantErr %v", tt.cName, err, tt.wantErr)
			}
		})
	}
}

func TestCobrandService_EnsureAdminIsUnique(t *testing.T) {
	svc := services.NewCobrandService()

	existing := []dtos.CobrandAdminDto{
		{Id: "ca1", User: "user1", Cobrand: "cobrand1"},
	}

	tests := []struct {
		name    string
		admins  []dtos.CobrandAdminDto
		userId  string
		wantErr bool
	}{
		{"first admin — empty list", []dtos.CobrandAdminDto{}, "user1", false},
		{"new user — not yet admin", existing, "user2", false},
		{"user already admin", existing, "user1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.EnsureAdminIsUnique(tt.admins, tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureAdminIsUnique(_, %q) error = %v, wantErr %v", tt.userId, err, tt.wantErr)
			}
		})
	}
}

func TestCobrandService_EnsureInviterIsApprovedAdmin(t *testing.T) {
	svc := services.NewCobrandService()

	approvedAdmin := []dtos.CobrandAdminDto{
		{Id: "ca1", User: "user1", Cobrand: "cobrand1", Approved: true},
	}
	unapprovedAdmin := []dtos.CobrandAdminDto{
		{Id: "ca1", User: "user1", Cobrand: "cobrand1", Approved: false},
	}

	tests := []struct {
		name      string
		admins    []dtos.CobrandAdminDto
		inviterId string
		wantErr   bool
	}{
		{"bootstrap — no existing admins", []dtos.CobrandAdminDto{}, "anyone", false},
		{"approved existing admin invites", approvedAdmin, "user1", false},
		{"unapproved existing admin invites", unapprovedAdmin, "user1", true},
		{"non-admin attempts to invite", approvedAdmin, "user2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.EnsureInviterIsApprovedAdmin(tt.admins, tt.inviterId)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureInviterIsApprovedAdmin(_, %q) error = %v, wantErr %v", tt.inviterId, err, tt.wantErr)
			}
		})
	}
}
