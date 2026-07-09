package services_test

import (
	"testing"

	"keybook/backend/internal/services"
)

func TestPersonService_ValidatePerson(t *testing.T) {
	svc := services.NewPersonService()

	tests := []struct {
		name    string
		pName   string
		dob     string
		wantErr bool
	}{
		{"valid", "Alice", "1990-01-01 00:00:00.000Z", false},
		{"empty name", "", "1990-01-01 00:00:00.000Z", true},
		{"empty dob", "Alice", "", true},
		{"both empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidatePerson(tt.pName, tt.dob)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePerson(%q, %q) error = %v, wantErr %v", tt.pName, tt.dob, err, tt.wantErr)
			}
		})
	}
}
