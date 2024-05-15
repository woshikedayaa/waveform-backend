package logf

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	core    zapcore.Core
	options []zap.Option
)

func Open(name string) *zap.Logger {
	var logger *zap.Logger

	logger = zap.New(core, options...)

	return logger.Named(name)
}
