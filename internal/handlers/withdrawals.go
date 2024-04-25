package handlers

import (
	"errors"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Withdrawals(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.Claims)
	userID := claims.UserID

	ctx := c.Request().Context()

	withdrawals, err := h.service.Withdrawals(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotExists) {
			return echo.NewHTTPError(http.StatusNoContent, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echo.NewHTTPError(http.StatusOK, withdrawals)
}
