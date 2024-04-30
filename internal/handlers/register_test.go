package handlers

import (
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/middleware/logger"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/dubrovsky1/gophermart/internal/service"
	"github.com/dubrovsky1/gophermart/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Register(t *testing.T) {
	logger.Initialize()

	type fields struct {
		server *echo.Echo
		ctrl   *gomock.Controller
	}

	tests := []struct {
		name           string
		error          error
		user           string
		userMock       models.User
		expectedStatus int
	}{
		{
			name:  "Register. Success.",
			error: nil,
			user:  `{"login": "Alex","password": "123456"}`,
			userMock: models.User{
				Login:    "Alex",
				Password: "123456",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "Register. Bad Request.",
			error: nil,
			user:  `{"login": "Alex","password": "123456"`,
			userMock: models.User{
				UserID:   "80600602-efa9-47b5-9919-68d6d982f8be",
				Login:    "Alex",
				Password: "123456",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Register. Already Exists.",
			error:          errs.ErrAlreadyExists,
			user:           `{"login": "Alex","password": "123456"}`,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Register. Internal Server Error.",
			error:          errs.ErrInternalServerError,
			user:           `{"login": "Alex","password": "123456"}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	s := fields{
		server: echo.New(),
		ctrl:   gomock.NewController(t),
	}
	defer s.ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//формируем запрос
			reqBody := strings.NewReader(tt.user)
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", reqBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := s.server.NewContext(req, rec)

			//заглушка для интерфейса Storager, который реализуется в сервисном слое
			m := mocks.NewMockStorager(s.ctrl)

			serv := service.New(m)
			handler := New(serv)

			//ожидание от интерфейса, реализуемого в сервисе
			m.EXPECT().Register(c.Request().Context(), gomock.Any()).Return(rec.Header().Get("Authorization"), tt.error).AnyTimes()

			errAdd := handler.Register(c)
			assert.NoError(t, errAdd)

			assert.Equal(t, tt.expectedStatus, rec.Code, "Код ответа не совпадает с ожидаемым")

			t.Log("=============================================================================================>")
		})
	}
}
