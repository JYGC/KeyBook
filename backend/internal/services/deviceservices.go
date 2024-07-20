package services

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
)

type IDeviceServices interface {
	AddNewDevicesIfNotExists(
		loggedInUserId string,
		propertyId string,
		inboundDevices []dtos.NewDeviceDto,
	) (
		[]dtos.DeviceDto,
		error,
	)
}

type DeviceServices struct {
	deviceRepository repositories.IDeviceRepository
	propertyServices IPropertyServices
}

func (d DeviceServices) AddNewDevicesIfNotExists(
	loggedInUserId string,
	propertyId string,
	inboundDevices []dtos.NewDeviceDto,
) (
	[]dtos.DeviceDto,
	error,
) {
	if checkPropertyOwnershipErr := d.propertyServices.ErrorIfPropertyNotOwnedByUserOrCantCheck(
		loggedInUserId,
		propertyId,
	); checkPropertyOwnershipErr != nil {
		return nil, checkPropertyOwnershipErr
	}

	existingDevices, existingDevicesErr := d.deviceRepository.GetDeviceForPropertyId(
		propertyId,
		inboundDevices,
	)
	if existingDevicesErr != nil {
		return nil, existingDevicesErr
	}
	existingDevicesMap := make(map[string]dtos.DeviceDto)
	for _, existingDevice := range existingDevices {
		existingDevicesMap[existingDevice.Identifier] = existingDevice
	}
	var newDevices []dtos.NewDeviceDto
	for _, inboundDevice := range inboundDevices {
		if _, isExistingDevice := existingDevicesMap[inboundDevice.Identifier]; !isExistingDevice {
			newDevices = append(newDevices, inboundDevice)
		}
	}

	if len(newDevices) > 0 {
		return d.deviceRepository.AddNewDevicesToProperty(propertyId, newDevices)
	}
	return existingDevices, nil
}

func NewDeviceServices(
	deviceRepository repositories.IDeviceRepository,
	propertyServices IPropertyServices,
) IDeviceServices {
	return &DeviceServices{
		deviceRepository,
		propertyServices,
	}
}
