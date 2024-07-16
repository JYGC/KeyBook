package handlers

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"net/http"

	"github.com/labstack/echo/v5"
)

type IHistoryHandlers interface {
	AddHistories(context echo.Context) error
	GetHistories(context echo.Context) error
}

type HistoryHandlers struct {
	historyRepository repositories.IHistoryRepository
}

func (h HistoryHandlers) AddHistories(context echo.Context) error {
	test2 := []interface{}{}
	test2 = append(test2, dtos.AddDeviceHistoriesDto{
		ActionDescription: "testinsfsa",
		DateSpecified:     "dssdgsdg",
	})
	h.historyRepository.AddHistories("personhistory", test2)
	return nil
}

func (h HistoryHandlers) GetHistories(context echo.Context) error {
	retval, getHistoriesErr := h.historyRepository.GetHistories("personhistory")
	if getHistoriesErr != nil {
		return context.JSON(http.StatusInternalServerError, getHistoriesErr)
	}
	return context.JSON(http.StatusOK, retval)
}

func NewHistoryHandlers(historyRepository repositories.IHistoryRepository) IHistoryHandlers {
	return HistoryHandlers{
		historyRepository,
	}
}

func RegisterHistoryHandlersToRouter(router *echo.Echo, historyHandlers IHistoryHandlers) {
	router.POST("/api/history/add", historyHandlers.AddHistories)
	router.GET("/api/history/get", historyHandlers.GetHistories)
}
