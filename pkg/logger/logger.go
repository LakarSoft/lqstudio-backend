package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init initializes the Zap logger
func Init(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// Sync flushes any buffered log entries
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}
