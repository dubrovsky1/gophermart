package handlers

import (
	"errors"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetOrderList(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.Claims)
	userID := claims.UserID

	if userID == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	ctx := c.Request().Context()

	orders, err := h.service.GetOrderList(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotExists) {
			return c.NoContent(http.StatusNoContent)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	return c.JSON(http.StatusOK, orders)
}
