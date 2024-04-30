package handlers

import (
	"errors"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Register(c echo.Context) error {
	reqBody := models.RegisterRequest{}

	if err := c.Bind(&reqBody); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user := models.User{
		UserID:   models.UserID(uuid.New().String()),
		Login:    reqBody.Login,
		Password: reqBody.Password,
	}

	ctx := c.Request().Context()
	token, err := h.service.Register(ctx, user)

	if err != nil {
		//логин уже занят
		if errors.Is(err, errs.ErrAlreadyExists) {
			return c.NoContent(http.StatusConflict)
		}
		//ошибка сервера
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("Authorization", string(token))

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
