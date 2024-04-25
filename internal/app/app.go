package app

import (
	"context"
	"github.com/dubrovsky1/gophermart/internal/handlers"
	"github.com/dubrovsky1/gophermart/internal/middleware/logger"
	"github.com/dubrovsky1/gophermart/internal/service"
	"github.com/dubrovsky1/gophermart/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Config  *Config
	Storage storage.Storager
	Handler *handlers.Handler
	Server  *echo.Echo
}

func New() *App {
	cfg := NewConfig(10*time.Second, 2)

	return &App{
		Config: cfg,
	}
}

func (a *App) Init() {
	var err error

	if err = a.Config.ParseFlags(); err != nil {
		log.Fatal("Parse flags error. ", err)
	}

	stor, err := storage.New(a.Config.ConnectionString, a.Config.ConnectionTimeout, a.Config.MigrateVersion)
	if err != nil {
		log.Fatal("Get storage error. ", err)
	}

	serv := service.New(*stor)
	handler := handlers.New(serv)

	server := echo.New()

	//сжатие данных
	server.Use(middleware.Gzip())

	//каждый запрос логгируем с помощью zap
	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogMethod: true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Sugar.Infoln(
				"uri", v.URI,
				"method", v.Method,
				"status", v.Status,
			)
			return nil
		},
	}))

	//регистрация и аутентификация доступна всем
	accessible := server.Group("")
	{
		accessible.POST("/api/user/register", handler.Register)
		accessible.POST("/api/user/login", handler.Login)
	}

	restricted := server.Group("")
	{
		config := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(service.Claims)
			},
			TokenLookup: "header:Authorization",
			SigningKey:  []byte(service.SecretKey),
		}
		restricted.Use(echojwt.WithConfig(config))

		restricted.POST("/api/user/orders", handler.AddOrder)
		restricted.GET("/api/user/orders", handler.GetOrderList)
		restricted.GET("/api/user/balance", handler.GetBalance)
		restricted.POST("/api/user/balance/withdraw", handler.Withdraw)
		restricted.GET("/api/user/withdrawals", handler.Withdrawals)
	}

	server.GET("/api/orders/:number", handler.GetOrderAccrualInfo)

	a.Storage = *stor
	a.Handler = handler
	a.Server = server
}

func (a *App) Run() {

	go func() {
		if err := a.Server.Start(a.Config.Host); err != http.ErrServerClosed {
			log.Fatalf("Listen and serve returned err: %v", err)
		}
	}()
	logger.Sugar.Infow("Server is listening", "host", a.Config.Host)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	if err := a.Server.Shutdown(ctx); err != nil {
		logger.Sugar.Infow("Server shutdown error", "err", err.Error())
	}

	a.Close()
	logger.Sugar.Infow("Shutting down server gracefully")
}

func (a *App) Close() {
	a.Storage.Close()
	logger.Sugar.Infow("Storage closed")
}
