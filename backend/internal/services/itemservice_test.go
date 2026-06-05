package services_test

import (
	"testing"

	"keybook/backend/internal/services"
)

func TestItemService_ValidateItem(t *testing.T) {
	svc := services.NewItemService()

	tests := []struct {
		name    string
		iName   string
		wantErr bool
	}{
		{"valid", "Front Door Key", false},
		{"empty name", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateItem(tt.iName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateItem(%q) error = %v, wantErr %v", tt.iName, err, tt.wantErr)
			}
		})
	}
}

func TestItemService_ValidateEntryDeviceTransition(t *testing.T) {
	svc := services.NewItemService()

	tests := []struct {
		name                 string
		currentDefunctReason string
		newDefunctReason     string
		wantErr              bool
	}{
		{"active stays active", "", "", false},
		{"active goes defunct with reason", "", "Lost", false},
		{"defunct stays defunct", "Lost", "Lost", false},
		{"defunct updates reason", "Lost", "Damaged", false},
		{"defunct reactivated — not allowed", "Lost", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateEntryDeviceTransition(tt.currentDefunctReason, tt.newDefunctReason)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEntryDeviceTransition(%q, %q) error = %v, wantErr %v",
					tt.currentDefunctReason, tt.newDefunctReason, err, tt.wantErr)
			}
		})
	}
}
