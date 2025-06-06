package utils

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

const (
	LoggerKey = "logger"
)

func Logger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(LoggerKey).(*zap.Logger)
	if !ok {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(fmt.Errorf("Error creating backup logger %s", err))
		}
		logger.Warn("Logger not found in context. Using default")
	}
	return logger
}
