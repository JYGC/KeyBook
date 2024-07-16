package services

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"time"
)

type IPropertyServices interface {
	AddPropertyIfNotExists(
		loggedInUserId string,
		propertyAddress string,
		startOfOwnership time.Time,
	) (dtos.PropertyIdAddressDto, error)
}

type PropertyServices struct {
	personRepository   repositories.IPersonRepository
	propertyRepository repositories.IPropertyRepository
}

func (p PropertyServices) AddPropertyIfNotExists(
	loggedInUserId string,
	propertyAddress string,
	startOfOwnership time.Time,
) (dtos.PropertyIdAddressDto, error) {
	properties, getPropertyByNameErr := p.propertyRepository.GetPropertiesManagedByUser(loggedInUserId, propertyAddress)
	if getPropertyByNameErr != nil {
		return dtos.PropertyIdAddressDto{}, getPropertyByNameErr
	}
	if len(properties) > 0 {
		return properties[0], nil
	}
	return p.propertyRepository.AddNewProperty(loggedInUserId, propertyAddress)
}

func NewPropertyServices(
	personRepository repositories.IPersonRepository,
	propertyRepository repositories.IPropertyRepository,
) IPropertyServices {
	propertyServices := PropertyServices{}
	propertyServices.personRepository = personRepository
	propertyServices.propertyRepository = propertyRepository
	return propertyServices
}
