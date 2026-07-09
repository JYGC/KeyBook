package services

import (
	"errors"

	"keybook/backend/internal/dtos"
)

type IPropertyOwnerService interface {
	EnsureNoDuplicatePersonOwner(existing []dtos.PersonPropertyOwnerDto, personId string) error
	EnsureNoDuplicateCobrandOwner(existing []dtos.CobrandPropertyOwnerDto, cobrandId string) error
}

type PropertyOwnerService struct{}

func NewPropertyOwnerService() IPropertyOwnerService {
	return &PropertyOwnerService{}
}

func (s *PropertyOwnerService) EnsureNoDuplicatePersonOwner(existing []dtos.PersonPropertyOwnerDto, personId string) error {
	for _, o := range existing {
		if o.Person == personId {
			return errors.New("person is already an owner via this property owner record")
		}
	}
	return nil
}

func (s *PropertyOwnerService) EnsureNoDuplicateCobrandOwner(existing []dtos.CobrandPropertyOwnerDto, cobrandId string) error {
	for _, o := range existing {
		if o.Cobrand == cobrandId {
			return errors.New("cobrand is already an owner via this property owner record")
		}
	}
	return nil
}
