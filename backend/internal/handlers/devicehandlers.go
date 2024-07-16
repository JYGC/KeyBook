package handlers

import (
	"encoding/json"
	"fmt"
	"keybook/backend/internal/services"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

type IDeviceHandlers interface {
	ImportCsv(context echo.Context) error
}

type DeviceHandlers struct {
	dataImportServices services.IDataImportServices
}

func (d DeviceHandlers) ImportCsv(context echo.Context) error {
	loggedInUser := apis.RequestInfo(context).AuthRecord
	csvContent := apis.RequestInfo(context).Data
	csvContentJson, csvContentJsonErr := json.Marshal(csvContent)
	if csvContentJsonErr != nil {
		fmt.Printf("err: %v\n", csvContentJsonErr)
		return csvContentJsonErr
	}

	d.dataImportServices.ProcessImportData(loggedInUser.Id, csvContentJson)
	return nil
}

func NewDeviceHandlers(dataImportServices services.IDataImportServices) IDeviceHandlers {
	deviceHandlers := DeviceHandlers{}
	deviceHandlers.dataImportServices = dataImportServices
	return deviceHandlers
}

func RegisterDeviceHandlersToRouter(router *echo.Echo, deviceHandlers IDeviceHandlers) {
	router.POST("/api/device/importcsv", deviceHandlers.ImportCsv)
}
