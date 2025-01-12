package services

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
)

type IPersonDeviceServices interface {
	AddNewPersonDeviceIfNotExists(
		loggedInUserId string,
		propertyId string,
		personDevices []dtos.NewPersonDeviceAndHistoriesDto,
	) (
		[]dtos.PersonDeviceDto,
		error,
	)
}

type PersonDeviceServices struct {
	personDeviceRepository repositories.IPersonDeviceRepository
	propertyServices       IPropertyServices
	personServices         IPersonServices
	deviceServices         IDeviceServices
}

func (p PersonDeviceServices) AddNewPersonDeviceIfNotExists(
	loggedInUserId string,
	propertyId string,
	personDevices []dtos.NewPersonDeviceAndHistoriesDto,
) (
	[]dtos.PersonDeviceDto,
	error,
) {
	persons, getPersonsErr := p.personServices.ValidateOwnerGetPersonsByProperty(
		loggedInUserId,
		propertyId,
	)
	if getPersonsErr != nil {
		return nil, getPersonsErr
	}
	personsMap := make(map[string]dtos.PersonIdNameDto)
	for _, person := range persons {
		personsMap[person.Id] = person
	}
	devices, getDevicesErr := p.deviceServices.ValidateOwnerGetDevicesByProperty(
		loggedInUserId,
		propertyId,
	)
	if getDevicesErr != nil {
		return nil, getDevicesErr
	}
	devicesMap := make(map[string]dtos.DeviceDto)
	for _, device := range devices {
		devicesMap[device.Id] = device
	}

	var personDevicesToAdd []dtos.NewPersonDeviceAndHistoriesDto
	for _, personDevice := range personDevices {
		_, isDeviceOfProperty := devicesMap[personDevice.Device]
		_, isPersonOfProperty := personsMap[personDevice.Person]
		if isDeviceOfProperty && isPersonOfProperty {
			personDevicesToAdd = append(personDevicesToAdd, personDevice)
		}
	}

	existingPersonDevices, getPersonDeviceErr := p.personDeviceRepository.GetPersonDevicesForProperty(
		personDevicesToAdd,
	)
	if getPersonDeviceErr != nil {
		return nil, getPersonDeviceErr
	}
	existingPersonDevicesMap := make(map[string]dtos.PersonDeviceDto)
	for _, existingPersonDevice := range existingPersonDevices {
		existingPersonDevicesMap[existingPersonDevice.Device] = existingPersonDevice
	}
	var newPersonDevicesToAdd []dtos.NewPersonDeviceAndHistoriesDto
	for _, personDeviceToAdd := range personDevicesToAdd {
		if _, isExistingPersonDevice := existingPersonDevicesMap[personDeviceToAdd.Device]; !isExistingPersonDevice {
			newPersonDevicesToAdd = append(newPersonDevicesToAdd, personDeviceToAdd)
		}
	}

	if len(newPersonDevicesToAdd) > 0 {
		return p.personDeviceRepository.AddNewPersonDevices(newPersonDevicesToAdd)
	}
	return nil, nil
}

func NewPersonDeviceServices(
	personDeviceRepository repositories.IPersonDeviceRepository,
	propertyServices IPropertyServices,
	personServices IPersonServices,
	deviceServices IDeviceServices,
) IPersonDeviceServices {
	return PersonDeviceServices{
		personDeviceRepository,
		propertyServices,
		personServices,
		deviceServices,
	}
}
