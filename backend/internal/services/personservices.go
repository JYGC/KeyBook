package services

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
)

type IPersonServices interface {
	AddNewPersonsIfNotExists(
		loggedInUserId string,
		propertyId string,
		personNames []string,
	) (
		[]dtos.PersonIdNameDto,
		error,
	)
}

type PersonServices struct {
	personRepository repositories.IPersonRepository
	propertyServices IPropertyServices
}

func (p PersonServices) AddNewPersonsIfNotExists(
	loggedInUserId string,
	propertyId string,
	personNames []string,
) (
	[]dtos.PersonIdNameDto,
	error,
) {
	if checkPropertyOwnershipErr := p.propertyServices.ErrorIfPropertyNotOwnedByUserOrCantCheck(
		loggedInUserId,
		propertyId,
	); checkPropertyOwnershipErr != nil {
		return nil, checkPropertyOwnershipErr
	}

	existingPersons, existingPersonsErr := p.personRepository.GetPersonsForPropertyIdByPersonNames(
		propertyId,
		personNames,
	)
	if existingPersonsErr != nil {
		return nil, existingPersonsErr
	}
	existingPersonNamesMap := make(map[string]string)
	for _, existingPerson := range existingPersons {
		existingPersonNamesMap[existingPerson.Name] = existingPerson.Name
	}
	var newPersonNames []string
	for _, personName := range personNames {
		if _, isExistingPersonName := existingPersonNamesMap[personName]; !isExistingPersonName {
			newPersonNames = append(newPersonNames, personName)
		}
	}

	if len(newPersonNames) > 0 {
		return p.personRepository.AddNewPersonsToProperty(propertyId, newPersonNames)
	}
	return existingPersons, nil
}

func NewPersonServices(
	personRepository repositories.IPersonRepository,
	propertyServices IPropertyServices,
) IPersonServices {
	return &PersonServices{
		personRepository,
		propertyServices,
	}
}
