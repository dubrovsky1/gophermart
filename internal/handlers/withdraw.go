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
		return c.NoContent(http.StatusBadRequest)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.Claims)
	userID := claims.UserID

	if userID == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	ctx := c.Request().Context()

	err := h.service.Withdraw(ctx, reqBody.OrderID, userID, reqBody.Sum)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNum) {
			return c.NoContent(http.StatusUnprocessableEntity) //неверный формат номера заказа
		} else if errors.Is(err, errs.ErrNotEnoughFunds) {
			return c.NoContent(http.StatusPaymentRequired) //на счету недостаточно средств
		}
		return c.NoContent(http.StatusInternalServerError) //внутренняя ошибка сервера
	}

	return c.NoContent(http.StatusOK)
}
