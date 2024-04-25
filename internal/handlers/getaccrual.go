package handlers

import (
	"errors"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *Handler) GetOrderAccrualInfo(c echo.Context) error {
	value := c.Param("number")
	orderID, _ := strconv.Atoi(value)

	ctx := c.Request().Context()

	orderAccrual, err := h.service.GetOrderAccrualInfo(ctx, models.OrderID(orderID))
	if err != nil {
		if errors.Is(err, errs.ErrNotExists) {
			return echo.NewHTTPError(http.StatusNoContent, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echo.NewHTTPError(http.StatusOK, orderAccrual)
}
