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

	if userID == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	ctx := c.Request().Context()

	withdrawals, err := h.service.Withdrawals(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotExists) {
			return c.NoContent(http.StatusNoContent)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return c.JSON(http.StatusOK, withdrawals)
}
