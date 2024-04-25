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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()) //неверный формат запроса
	}

	orderID, err := strconv.Atoi(string(body))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()) //неверный формат запроса
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.Claims)
	userID := claims.UserID

	ctx := c.Request().Context()

	err = h.service.AddOrder(ctx, models.OrderID(orderID), userID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderAlreadyLoadThisUser) {
			return echo.NewHTTPError(http.StatusOK, err.Error()) //номер заказа уже был загружен этим пользователем
		} else if errors.Is(err, errs.ErrOrderLoadAnotherUser) {
			return echo.NewHTTPError(http.StatusConflict, err.Error()) //номер заказа уже был загружен другим пользователем
		} else if errors.Is(err, errs.ErrOrderNum) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()) //неверный формат номера заказа
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error()) //внутренняя ошибка сервера
	}

	return echo.NewHTTPError(http.StatusAccepted)
}
