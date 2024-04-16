package logger

import (
	"go.uber.org/zap"
)

var Sugar zap.SugaredLogger

func Initialize() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()

	Sugar = *logger.Sugar()

	return nil
}
