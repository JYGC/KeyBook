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
	for _, a := range admins {
		if a.User == userId {
			return errors.New("user is already an admin of this cobrand")
		}
	}
	return nil
}

// EnsureInviterIsApprovedAdmin allows the bootstrap case (no existing admins)
// through unconditionally. Once a cobrand has at least one admin, only an
// existing admin whose own approved field is true may add another.
func (s *CobrandService) EnsureInviterIsApprovedAdmin(admins []dtos.CobrandAdminDto, inviterUserId string) error {
	if len(admins) == 0 {
		return nil
	}
	for _, a := range admins {
		if a.User == inviterUserId {
			if !a.Approved {
				return errors.New("admin must be approved before adding another admin")
			}
			return nil
		}
	}
	return errors.New("only an existing admin may add another admin to this cobrand")
}
