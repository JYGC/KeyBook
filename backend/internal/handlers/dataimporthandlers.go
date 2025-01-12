package handlers

import (
	"encoding/json"
	"keybook/backend/internal/services"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

type IDataImportHandlers interface {
	ImportCsv(context echo.Context) error
}

type DataImportHandlers struct {
	dataImportServices services.IDataImportServices
}

func (d DataImportHandlers) ImportCsv(context echo.Context) error {
	loggedInUser := apis.RequestInfo(context).AuthRecord
	csvContent := apis.RequestInfo(context).Data
	csvContentJson, csvContentJsonErr := json.Marshal(csvContent)
	if csvContentJsonErr != nil {
		return csvContentJsonErr
	}

	d.dataImportServices.ProcessImportData(loggedInUser.Id, csvContentJson)
	return nil
}

func NewDataImportHandlers(dataImportServices services.IDataImportServices) IDataImportHandlers {
	dataImportHandlers := DataImportHandlers{}
	dataImportHandlers.dataImportServices = dataImportServices
	return dataImportHandlers
}

func RegisterDataImportHandlersToRouter(router *echo.Echo, dataImportHandlers IDataImportHandlers) {
	router.POST("/dataimport/csv", dataImportHandlers.ImportCsv)
}
