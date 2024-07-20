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
) []dtos.NewDeviceDto {
	var importDevices []dtos.NewDeviceDto
	for _, dpdh := range importDataDto.DevicesPersonDevicesAndHistories {
		importDevices = append(importDevices, dtos.NewDeviceDto{
			Name:          dpdh.Name,
			Identifier:    dpdh.Identifier,
			Type:          dpdh.Type,
			DefunctReason: dpdh.DefunctReason,
		})
	}
	return importDevices
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
				dtos.NewPersonDeviceHistoriesDto{
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

	property, addPropertyErr := d.propertyServices.AddPropertyIfNotExists(
		loggedInUserId,
		importDataDto.PropertyAddress,
		time.Now(),
	)
	if addPropertyErr != nil {
		return addPropertyErr
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

	devices := d.getDevicesFromImportDataDto(importDataDto)
	newDevices, addDeviceErr := d.deviceServices.AddNewDevicesIfNotExists(
		loggedInUserId,
		property.Id,
		devices,
	)
	if addDeviceErr != nil {
		return addDeviceErr
	}

	importPersonDevices, err := d.getPersonDevicesAndHistoriesFromImportDataDto(
		importDataDto,
		newDevices,
		newPersons,
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	m, e := json.Marshal(importPersonDevices)
	if e != nil {
		fmt.Printf("e: %v\n", e)
		return e
	}
	fmt.Printf("m: %v\n", string(m))
	//_, addPersonDeviceErr := d.personDeviceServices.AddNewPersonDevices
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
