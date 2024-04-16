package main

import (
	"github.com/dubrovsky1/gophermart/internal/app"
	"github.com/dubrovsky1/gophermart/internal/middleware/logger"
)

func main() {
	logger.Initialize()

	a := app.New()
	a.Init()

	logger.Sugar.Infow("Flags:", "-a", a.Config.Host, "-d", a.Config.ConnectionString, "-f", a.Config.AccrualAddress)

	a.Run()
}
