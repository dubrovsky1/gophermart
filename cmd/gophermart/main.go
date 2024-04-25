package main

import (
	"bytes"
	"github.com/dubrovsky1/gophermart/internal/app"
	"github.com/dubrovsky1/gophermart/internal/middleware/logger"
	"os/exec"
)

func main() {
	logger.Initialize()

	a := app.New()
	a.Init()

	logger.Sugar.Infow("Flags:", "-a", a.Config.Host, "-d", a.Config.ConnectionString, "-f", a.Config.AccrualAddress)

	a.Run()

	//args := []string{a.Config.Host, a.Config.ConnectionString, a.Config.AccrualAddress}
	//cmd := exec.Command(a.Config.AccrualAddress, args...)
	cmd := exec.Command("cmd", "/C", a.Config.AccrualAddress, "-a", a.Config.Host, "-d", a.Config.ConnectionString)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Start()
	if err != nil {
		logger.Sugar.Infow("Accrual start Log", "err", err)
	}
	err = cmd.Wait()
	if err != nil {
		logger.Sugar.Infow("Accrual Wait Log", "err", err)
	}

	logger.Sugar.Infow("Accrual start Log", "buf", buf.String())

}
