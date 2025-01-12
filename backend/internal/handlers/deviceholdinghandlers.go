package handlers

import (
	"encoding/json"
	"keybook/backend/internal/services"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

type IDeviceHoldingHandlers interface {
	ListDevicesHeld(context echo.Context) error
}

type DeviceHoldingHandlers struct {
	dataImportServices services.IDataImportServices
}

func (d DeviceHoldingHandlers) ListDevicesHeld(context echo.Context) error {
	loggedInUser := apis.RequestInfo(context).AuthRecord
	csvContent := apis.RequestInfo(context).Data
	csvContentJson, csvContentJsonErr := json.Marshal(csvContent)
	if csvContentJsonErr != nil {
		return csvContentJsonErr
	}

	d.dataImportServices.ProcessImportData(loggedInUser.Id, csvContentJson)
	return nil
}

func NewDeviceHandlers(dataImportServices services.IDataImportServices) IDeviceHoldingHandlers {
	deviceHandlers := DeviceHoldingHandlers{}
	deviceHandlers.dataImportServices = dataImportServices
	return deviceHandlers
}

func RegisterDeviceHandlersToRouter(router *echo.Echo, deviceHandlers IDeviceHoldingHandlers) {
	router.POST("/device/importcsv", deviceHandlers.ListDevicesHeld)
}
