package handlers

import (
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/middleware/logger"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/dubrovsky1/gophermart/internal/service"
	"github.com/dubrovsky1/gophermart/internal/service/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_GetBalance(t *testing.T) {
	logger.Initialize()

	type fields struct {
		server *echo.Echo
		ctrl   *gomock.Controller
	}

	tests := []struct {
		name           string
		userID         string
		error          error
		balance        models.Balance
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:   "Get balance. Success.",
			userID: "80600602-efa9-47b5-9919-68d6d982f8be",
			error:  nil,
			balance: models.Balance{
				Current:   500.5,
				Withdrawn: 42,
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"current": 500.5,"withdrawn": 42}`,
		},
		{
			name:   "Get balance. Internal Server Error.",
			userID: "80600602-efa9-47b5-9919-68d6d982f8be",
			error:  errs.ErrInternalServerError,
			balance: models.Balance{
				Current:   500.5,
				Withdrawn: 42,
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:  "Get balance. Unauthorized.",
			error: nil,
			balance: models.Balance{
				Current:   500.5,
				Withdrawn: 42,
			},
			expectedStatus: http.StatusUnauthorized,
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
			reqBody := strings.NewReader("")
			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", reqBody)
			rec := httptest.NewRecorder()

			c := s.server.NewContext(req, rec)

			//заглушка для интерфейса Storager, который реализуется в сервисном слое
			m := mocks.NewMockStorager(s.ctrl)

			//ожидание от интерфейса, реализуемого в сервисе
			m.EXPECT().GetBalance(c.Request().Context(), models.UserID(tt.userID)).Return(tt.balance, tt.error).AnyTimes()

			serv := service.New(m)
			handler := New(serv)

			//аутентификация
			token, err := service.BuildJWTString(models.UserID(tt.userID))
			require.NoError(t, err)
			req.Header.Set(echo.HeaderAuthorization, token)

			config := echojwt.Config{
				NewClaimsFunc: func(c echo.Context) jwt.Claims {
					return new(service.Claims)
				},
				TokenLookup: "header:Authorization",
				SigningKey:  []byte(service.SecretKey),
			}

			//вызов функции обработчика обернутой в JWT MiddleWare
			errAdd := echojwt.WithConfig(config)(handler.GetBalance)(c)
			assert.NoError(t, errAdd)

			assert.Equal(t, tt.expectedStatus, rec.Code, "Код ответа не совпадает с ожидаемым")

			if rec.Code == http.StatusOK {
				assert.JSONEq(t, tt.expectedJSON, string(rec.Body.Bytes()), "Тело ответа не совпадает с ожидаемым")
			}

			t.Log("=============================================================================================>")
		})
	}
}
