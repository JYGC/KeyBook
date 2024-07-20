package services

import (
	"fmt"
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"time"
)

type IPropertyServices interface {
	ErrorIfPropertyNotOwnedByUserOrCantCheck(
		loggedInUserId string,
		propertyId string,
	) error
	AddPropertyIfNotExists(
		loggedInUserId string,
		propertyAddress string,
		startOfOwnership time.Time,
	) (
		dtos.PropertyIdAddressDto,
		error,
	)
}

type PropertyServices struct {
	propertyRepository repositories.IPropertyRepository
}

func (p PropertyServices) ErrorIfPropertyNotOwnedByUserOrCantCheck(
	loggedInUserId string,
	propertyId string,
) error {
	isPropertyBelongToUser, testOwnershipErr := p.propertyRepository.IsPropertyIdBelongToUser(
		loggedInUserId,
		propertyId,
	)
	if testOwnershipErr != nil {
		return testOwnershipErr
	}
	if !isPropertyBelongToUser {
		return fmt.Errorf(
			"userid: %s does not own propertyid: %s",
			loggedInUserId,
			propertyId,
		)
	}
	return nil
}

func (p PropertyServices) AddPropertyIfNotExists(
	loggedInUserId string,
	propertyAddress string,
	startOfOwnership time.Time,
) (
	dtos.PropertyIdAddressDto,
	error,
) {
	properties, getPropertyByNameErr := p.propertyRepository.GetPropertiesForUserByPropertyName(
		loggedInUserId,
		propertyAddress,
	)
	if getPropertyByNameErr != nil {
		return dtos.PropertyIdAddressDto{}, getPropertyByNameErr
	}
	if len(properties) > 0 {
		return properties[0], nil
	}
	return p.propertyRepository.AddNewProperty(loggedInUserId, propertyAddress)
}

func NewPropertyServices(
	propertyRepository repositories.IPropertyRepository,
) IPropertyServices {
	return PropertyServices{
		propertyRepository,
	}
}
