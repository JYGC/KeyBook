package services

import (
	"errors"

	"keybook/backend/internal/dtos"
)

type IPropertyService interface {
	ValidateProperty(address string) error
	IsPersonOwner(owners []dtos.PersonPropertyOwnerDto, personId string) bool
	IsCobrandOwner(owners []dtos.CobrandPropertyOwnerDto, cobrandId string) bool
}

type PropertyService struct{}

func NewPropertyService() IPropertyService {
	return &PropertyService{}
}

func (s *PropertyService) ValidateProperty(address string) error {
	if address == "" {
		return errors.New("address is required")
	}
	return nil
}

func (s *PropertyService) IsPersonOwner(owners []dtos.PersonPropertyOwnerDto, personId string) bool {
	for _, o := range owners {
		if o.Person == personId {
			return true
		}
	}
	return false
}

func (s *PropertyService) IsCobrandOwner(owners []dtos.CobrandPropertyOwnerDto, cobrandId string) bool {
	for _, o := range owners {
		if o.Cobrand == cobrandId {
			return true
		}
	}
	return false
}
