package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	Logger *zap.SugaredLogger
	atom   zap.AtomicLevel
	once   sync.Once
)

func InitLogger(level string) {
	once.Do(func() {
		atom = zap.NewAtomicLevel()
		setLogLevel(level)

		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.Level = atom

		logger, _ := config.Build()
		Logger = logger.Sugar()
	})
}

func GetLogger() *zap.SugaredLogger {
	if Logger == nil {
		InitLogger("info")
	}
	return Logger
}

func SetLogLevel(level string) {
	setLogLevel(level)
}

func setLogLevel(level string) {
	var l zapcore.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		l = zapcore.InfoLevel
	}
	atom.SetLevel(l)
}