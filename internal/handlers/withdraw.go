package handlers

import (
	"errors"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/dubrovsky1/gophermart/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Withdraw(c echo.Context) error {
	reqBody := models.Withdraw{}

	if err := c.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.Claims)
	userID := claims.UserID

	ctx := c.Request().Context()

	err := h.service.Withdraw(ctx, reqBody.OrderID, userID, reqBody.Sum)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNum) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()) //неверный формат номера заказа
		} else if errors.Is(err, errs.ErrNotEnoughFunds) {
			return echo.NewHTTPError(http.StatusPaymentRequired, err.Error()) //на счету недостаточно средств
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()) //внутренняя ошибка сервера
	}

	return echo.NewHTTPError(http.StatusOK)
}
