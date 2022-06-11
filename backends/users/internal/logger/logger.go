package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(level zapcore.Level) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	cfg.Encoding = "json"
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{l}, nil
}
