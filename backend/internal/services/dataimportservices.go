package services

import (
	"encoding/json"
	"fmt"
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"strings"
	"time"
)

// all db calls here must have dao := app.Dao().WithoutHooks()

type IDataImportServices interface {
	ProcessImportData(loggedInUserId string, importDateJson []byte) error
}

type DataImportServices struct {
	personRepository repositories.IPersonRepository
	propertyServices IPropertyServices
}

func (d DataImportServices) ProcessImportData(loggedInUserId string, importDateJson []byte) error {
	var importDataDto dtos.AddPropertyDeviceAndHistoriesDto
	json.Unmarshal(importDateJson, &importDataDto)

	personNames := make(map[string]string)
	for _, dpdh := range importDataDto.DevicesPersonDevicesAndHistories {
		for _, pdh := range dpdh.PersonDeviceHistories {
			if strings.TrimSpace(pdh.DeviceHolder) != "" {
				personNames[pdh.DeviceHolder] = pdh.DeviceHolder
			}
		}
	}

	fmt.Printf("loggedInUser: %v\n", loggedInUserId)
	d.propertyServices.AddPropertyIfNotExists(loggedInUserId, importDataDto.PropertyAddress, time.Now())

	return nil
}

func NewDataImportServices(
	personRepository repositories.IPersonRepository,
	propertyServices IPropertyServices,
) IDataImportServices {
	return DataImportServices{
		personRepository,
		propertyServices,
	}
}
