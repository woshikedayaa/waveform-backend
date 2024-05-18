package logf

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	core    zapcore.Core
	options []zap.Option //用于存储日志配置的选项
)

func Open(name string) *zap.Logger {
	var logger *zap.Logger

	// 使用 zap.New 和 配置好的 core / options 创建一个 zap.Logger
	logger = zap.New(core, options...)
	// 返回一个配置完成、具有名称的 zap.Logger 实例
	return logger.Named(name)
}
