package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLogger(level string) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var l zapcore.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		l = zapcore.InfoLevel
	}
	config.Level.SetLevel(l)

	logger, _ := config.Build()
	Logger = logger.Sugar()
}

// Add this function to get the logger
func GetLogger() *zap.SugaredLogger {
	if Logger == nil {
		InitLogger("info")
	}
	return Logger
}