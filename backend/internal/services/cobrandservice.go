package services

import (
	"errors"

	"keybook/backend/internal/dtos"
)

type ICobrandService interface {
	ValidateCobrand(name string) error
	EnsureAdminIsUnique(admins []dtos.CobrandAdminDto, userId string) error
	EnsureInviterIsApprovedAdmin(admins []dtos.CobrandAdminDto, inviterUserId string) error
}

type CobrandService struct{}

func NewCobrandService() ICobrandService {
	return &CobrandService{}
}

func (s *CobrandService) ValidateCobrand(name string) error {
	if name == "" {
		return errors.New("cobrand name is required")
	}
	return nil
}

func (s *CobrandService) EnsureAdminIsUnique(admins []dtos.CobrandAdminDto, userId string) error {
	for _, admin := range admins {
		if admin.User == userId {
			return errors.New("user is already an admin of this cobrand")
		}
	}
	return nil
}

// A cobrand with no admins yet is exempt: the first admin has nobody to be
// approved by, so bootstrapping would otherwise be impossible.
func (s *CobrandService) EnsureInviterIsApprovedAdmin(admins []dtos.CobrandAdminDto, inviterUserId string) error {
	if len(admins) == 0 {
		return nil
	}
	for _, admin := range admins {
		if admin.User == inviterUserId {
			if !admin.Approved {
				return errors.New("admin must be approved before adding another admin")
			}
			return nil
		}
	}
	return errors.New("only an existing admin may add another admin to this cobrand")
}
