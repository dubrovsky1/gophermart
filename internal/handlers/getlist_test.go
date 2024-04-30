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

func TestHandler_GetOrderList(t *testing.T) {
	logger.Initialize()

	type fields struct {
		server *echo.Echo
		ctrl   *gomock.Controller
	}

	tests := []struct {
		name           string
		userID         string
		error          error
		orders         []models.Order
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:   "Get list. Success.",
			userID: "80600602-efa9-47b5-9919-68d6d982f8be",
			error:  nil,
			orders: []models.Order{
				{
					OrderID: 12345678903,
					Status:  "PROCESSED",
					Accrual: 500,
					Upload:  "2024-12-13T15:15:45+03:00",
				},
				{
					OrderID: 2377225624,
					Status:  "PROCESSING",
					Accrual: 0,
					Upload:  "2024-12-10T15:15:45+03:00",
				},
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   `[{"number": 12345678903,"status": "PROCESSED","accrual": 500,"uploaded_at": "2024-12-13T15:15:45+03:00"},{"number": 2377225624,"status": "PROCESSING","uploaded_at": "2024-12-10T15:15:45+03:00"}]`,
		},
		{
			name:           "Get list. Not exists orders.",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          errs.ErrNotExists,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Get list. Internal server error.",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          errs.ErrInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Get list. Internal server error.",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          errs.ErrInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Get list. Unauthorized.",
			error:          nil,
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
			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", reqBody)
			rec := httptest.NewRecorder()

			c := s.server.NewContext(req, rec)

			//заглушка для интерфейса Storager, который реализуется в сервисном слое
			m := mocks.NewMockStorager(s.ctrl)

			//ожидание от интерфейса, реализуемого в сервисе
			m.EXPECT().GetOrderList(c.Request().Context(), models.UserID(tt.userID)).Return(tt.orders, tt.error).AnyTimes()

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
			errAdd := echojwt.WithConfig(config)(handler.GetOrderList)(c)
			assert.NoError(t, errAdd)

			assert.Equal(t, tt.expectedStatus, rec.Code, "Код ответа не совпадает с ожидаемым")

			if rec.Code == http.StatusOK {
				assert.JSONEq(t, tt.expectedJSON, string(rec.Body.Bytes()), "Тело ответа не совпадает с ожидаемым")
			}

			t.Log("=============================================================================================>")
		})
	}
}
