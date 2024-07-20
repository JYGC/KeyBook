package services

import (
	"encoding/json"
	"fmt"
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"strings"
	"time"
)

type IDataImportServices interface {
	ProcessImportData(loggedInUserId string, importDateJson []byte) error
}

type DataImportServices struct {
	personRepository     repositories.IPersonRepository
	propertyServices     IPropertyServices
	personServices       IPersonServices
	deviceServices       IDeviceServices
	personDeviceServices IPersonDeviceServices
}

func (d DataImportServices) getPersonNamesFromImportDataDto(
	importDataDto dtos.AddPropertyDeviceAndHistoriesDto,
) []string {
	var personNames []string
	personNamesMap := make(map[string]string)
	for _, dpdh := range importDataDto.DevicesPersonDevicesAndHistories {
		for _, pdh := range dpdh.PersonDeviceHistories {
			if strings.TrimSpace(pdh.DeviceHolder) != "" {
				personNamesMap[pdh.DeviceHolder] = pdh.DeviceHolder
			}
		}
	}
	for personName := range personNamesMap {
		personNames = append(personNames, personName)
	}
	return personNames
}

func (d DataImportServices) getDevicesFromImportDataDto(
	importDataDto dtos.AddPropertyDeviceAndHistoriesDto,
) (
	[]dtos.NewDeviceDto,
	error,
) {
	var importDevices []dtos.NewDeviceDto
	for _, dpdh := range importDataDto.DevicesPersonDevicesAndHistories {
		importDevice := dtos.NewDeviceDto{
			Name:          dpdh.Name,
			Identifier:    dpdh.Identifier,
			Type:          dpdh.Type,
			DefunctReason: dpdh.DefunctReason,
		}
		for _, dh := range dpdh.DeviceHistories {
			dateSpecifiedTime, parseDateErr := time.Parse(
				time.RFC3339,
				dh.DateSpecified,
			)
			if parseDateErr != nil {
				return nil, parseDateErr
			}
			importDevice.Histories = append(importDevice.Histories, dtos.NewDeviceHistoryDto{
				Description:    dh.ActionDescription,
				StatedDateTime: dateSpecifiedTime,
			})
		}
		importDevices = append(importDevices, importDevice)
	}
	return importDevices, nil
}

func (d DataImportServices) getPersonDevicesAndHistoriesFromImportDataDto(
	importDataDto dtos.AddPropertyDeviceAndHistoriesDto,
	devices []dtos.DeviceDto,
	persons []dtos.PersonIdNameDto,
) (
	[]dtos.NewPersonDeviceAndHistoriesDto,
	error,
) {
	var importPersonDevices []dtos.NewPersonDeviceAndHistoriesDto

	devicesMap := make(map[string]dtos.DeviceDto)
	for _, device := range devices {
		devicesMap[device.Identifier] = device
	}

	personsMap := make(map[string]dtos.PersonIdNameDto)
	for _, person := range persons {
		personsMap[person.Name] = person
	}

	for _, dpdh := range importDataDto.DevicesPersonDevicesAndHistories {
		importPersonDevice := dtos.NewPersonDeviceAndHistoriesDto{
			Device: devicesMap[dpdh.Identifier].Id,
			Person: personsMap[dpdh.CurrentHolder].Id,
		}
		for _, pdh := range dpdh.PersonDeviceHistories {
			dateSpecifiedTime, parseDateErr := time.Parse(
				time.RFC3339,
				pdh.DateSpecified,
			)
			if parseDateErr != nil {
				return nil, parseDateErr
			}
			importPersonDevice.Histories = append(
				importPersonDevice.Histories,
				dtos.NewPersonDeviceHistoryDto{
					Device:         devicesMap[dpdh.Identifier].Id,
					Person:         personsMap[pdh.DeviceHolder].Id,
					StatedDateTime: dateSpecifiedTime,
					Description:    pdh.ActionDescription,
				},
			)
		}
		importPersonDevices = append(
			importPersonDevices,
			importPersonDevice,
		)
	}

	return importPersonDevices, nil
}

func (d DataImportServices) ProcessImportData(loggedInUserId string, importDateJson []byte) error {
	var importDataDto dtos.AddPropertyDeviceAndHistoriesDto
	json.Unmarshal(importDateJson, &importDataDto)

	property, addPropertyErr := d.propertyServices.AddPropertyIfNotExistsReturnEmptyIfExists(
		loggedInUserId,
		importDataDto.PropertyAddress,
		time.Now(),
	)
	if addPropertyErr != nil {
		return addPropertyErr
	}
	if property.Id == "" {
		return fmt.Errorf("property exists")
	}

	personNames := d.getPersonNamesFromImportDataDto(importDataDto)
	newPersons, addPersonErr := d.personServices.AddNewPersonsIfNotExists(
		loggedInUserId,
		property.Id,
		personNames,
	)
	if addPersonErr != nil {
		return addPersonErr
	}

	devices, getDevicesErr := d.getDevicesFromImportDataDto(importDataDto)
	if getDevicesErr != nil {
		return getDevicesErr
	}
	newDevices, addDeviceErr := d.deviceServices.AddNewDevicesIfNotExists(
		loggedInUserId,
		property.Id,
		devices,
	)
	if addDeviceErr != nil {
		return addDeviceErr
	}

	newPersonDevices, getPersonDeviceErr := d.getPersonDevicesAndHistoriesFromImportDataDto(
		importDataDto,
		newDevices,
		newPersons,
	)
	if getPersonDeviceErr != nil {
		fmt.Printf("err: %v\n", getPersonDeviceErr)
		return getPersonDeviceErr
	}

	_, addPersonDeviceErr := d.personDeviceServices.AddNewPersonDeviceIfNotExists(
		loggedInUserId,
		property.Id,
		newPersonDevices,
	)
	if addPersonDeviceErr != nil {
		return addPersonDeviceErr
	}

	return nil
}

func NewDataImportServices(
	personRepository repositories.IPersonRepository,
	propertyServices IPropertyServices,
	personServices IPersonServices,
	deviceServices IDeviceServices,
	personDeviceServices IPersonDeviceServices,
) IDataImportServices {
	return DataImportServices{
		personRepository,
		propertyServices,
		personServices,
		deviceServices,
		personDeviceServices,
	}
}
