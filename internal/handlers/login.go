package handlers

import (
	"errors"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Login(c echo.Context) error {
	reqBody := models.RegisterRequest{}

	if err := c.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := models.User{
		Login:    reqBody.Login,
		Password: reqBody.Password,
	}

	ctx := c.Request().Context()

	//если получили токен, значит успешно нашли пользователя в базе и записали его id в Claims
	token, err := h.service.Login(ctx, user)
	if err != nil {
		if errors.Is(err, errs.ErrNotExists) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set("Authorization", string(token))

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
