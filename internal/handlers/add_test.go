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

func TestHandler_AddOrder(t *testing.T) {
	logger.Initialize()

	type fields struct {
		server *echo.Echo
		ctrl   *gomock.Controller
	}

	tests := []struct {
		name           string
		orderID        string
		userID         string
		error          error
		expectedStatus int
	}{
		{
			name:           "Add order. Success.",
			orderID:        "12345678903",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          nil,
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Add order. Already Load This User.",
			orderID:        "12345678903",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          errs.ErrOrderAlreadyLoadThisUser,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Add order. Order Load Another User.",
			orderID:        "12345678903",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          errs.ErrOrderLoadAnotherUser,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Add order. Error Order Num.",
			orderID:        "123454",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          errs.ErrOrderNum,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Add order. Internal Server Error.",
			orderID:        "12345678903",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          errs.ErrInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Add order. Body Error.",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Add order. orderID Error.",
			orderID:        "12345678903Dasf4",
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			error:          errs.ErrOrderNum,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Add order. User Unauthorized.",
			orderID:        "12345678903",
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
			reqBody := strings.NewReader(tt.orderID)
			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", reqBody)
			rec := httptest.NewRecorder()

			c := s.server.NewContext(req, rec)

			//заглушка для интерфейса Storager, который реализуется в сервисном слое
			m := mocks.NewMockStorager(s.ctrl)

			//ожидание от интерфейса, реализуемого в сервисе
			m.EXPECT().AddOrder(c.Request().Context(), models.OrderID(tt.orderID), models.UserID(tt.userID)).Return(tt.error).AnyTimes()

			serv := service.New(m)
			handler := New(serv)

			//t.Log("ord.", models.OrderID(ord))
			//t.Log("userID.", models.UserID(tt.userID))

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
			errAdd := echojwt.WithConfig(config)(handler.AddOrder)(c)
			assert.NoError(t, errAdd)

			assert.Equal(t, tt.expectedStatus, rec.Code, "Код ответа не совпадает с ожидаемым")

			//if assert.NotNil(t, errAdd) {
			//	he, ok := errAdd.(*echo.HTTPError)
			//	if ok {
			//		assert.Equal(t, tt.ExpectedStatus, he.Code, "Код ответа не совпадает с ожидаемым")
			//	}
			//}

			t.Log("=============================================================================================>")
		})
	}
}
