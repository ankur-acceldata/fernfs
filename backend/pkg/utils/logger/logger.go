package logger

import (
	"github.com/ankuragarwal/fernfs/backend/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(cfg *config.LoggingConfig) error {
	var err error
	var config zap.Config

	if cfg.Format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	config.Level = zap.NewAtomicLevelAt(level)

	log, err = config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(log)
	return nil
}

func Logger() *zap.Logger {
	return log
}

func Sync() error {
	return log.Sync()
} 