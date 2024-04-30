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

func TestHandler_Withdraw(t *testing.T) {
	logger.Initialize()

	type fields struct {
		server *echo.Echo
		ctrl   *gomock.Controller
	}

	tests := []struct {
		name           string
		withdraw       string
		userID         string
		orderid        models.OrderID
		sum            float64
		error          error
		expectedStatus int
	}{
		{
			name:           "Withdraw. Success.",
			withdraw:       `{"order": "2377225624","sum": 751.5}`,
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			orderid:        2377225624,
			sum:            751.5,
			error:          nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Withdraw. Error Order Num.",
			withdraw:       `{"order": "0000","sum": 751.5}`,
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			orderid:        0000,
			sum:            751.5,
			error:          errs.ErrOrderNum,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Withdraw. Error Not Enough Funds.",
			withdraw:       `{"order": "2377225624","sum": 751.5}`,
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			orderid:        2377225624,
			sum:            751.5,
			error:          errs.ErrNotEnoughFunds,
			expectedStatus: http.StatusPaymentRequired,
		},
		{
			name:           "Withdraw. Internal Server Error.",
			withdraw:       `{"order": "2377225624","sum": 751.5}`,
			userID:         "80600602-efa9-47b5-9919-68d6d982f8be",
			orderid:        2377225624,
			sum:            751.5,
			error:          errs.ErrInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Withdraw. User Unauthorized.",
			withdraw:       `{"order": "2377225624","sum": 751.5}`,
			orderid:        2377225624,
			sum:            751.5,
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
			reqBody := strings.NewReader(tt.withdraw)
			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", reqBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := s.server.NewContext(req, rec)

			//заглушка для интерфейса Storager, который реализуется в сервисном слое
			m := mocks.NewMockStorager(s.ctrl)

			//ожидание от интерфейса, реализуемого в сервисе
			m.EXPECT().Withdraw(c.Request().Context(), tt.orderid, models.UserID(tt.userID), tt.sum).Return(tt.error).AnyTimes()

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
			errAdd := echojwt.WithConfig(config)(handler.Withdraw)(c)
			assert.NoError(t, errAdd)

			assert.Equal(t, tt.expectedStatus, rec.Code, "Код ответа не совпадает с ожидаемым")

			t.Log("=============================================================================================>")
		})
	}
}
