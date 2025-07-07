package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes and returns a configured logger
func InitLogger() *zap.Logger {
	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create logger config
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig = encoderConfig

	// Create logger
	logger, err := loggerConfig.Build()
	if err != nil {
		// If we can't create the logger, use a default one
		logger, _ = zap.NewProduction()
	}

	return logger
} 