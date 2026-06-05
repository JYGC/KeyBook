package services

import (
	"errors"

	"keybook/backend/internal/dtos"
)

type ICobrandService interface {
	ValidateCobrand(name string) error
	EnsureAdminIsUnique(admins []dtos.CobrandAdminDto, userId string) error
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
