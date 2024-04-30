package handlers

import (
	"errors"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/dubrovsky1/gophermart/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) AddOrder(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest) //неверный формат запроса
	}

	orderID, err := strconv.Atoi(string(body))
	if err != nil {
		return c.NoContent(http.StatusBadRequest) //неверный формат запроса
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.Claims)
	userID := claims.UserID

	if userID == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	ctx := c.Request().Context()

	err = h.service.AddOrder(ctx, models.OrderID(orderID), userID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderAlreadyLoadThisUser) {
			return c.NoContent(http.StatusOK) //номер заказа уже был загружен этим пользователем
		} else if errors.Is(err, errs.ErrOrderLoadAnotherUser) {
			return c.NoContent(http.StatusConflict) //номер заказа уже был загружен другим пользователем
		} else if errors.Is(err, errs.ErrOrderNum) {
			return c.NoContent(http.StatusUnprocessableEntity) //неверный формат номера заказа
		}
		return c.NoContent(http.StatusInternalServerError) //внутренняя ошибка сервера
	}

	return c.NoContent(http.StatusAccepted)
}
