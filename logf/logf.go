package logf

import (
	"errors"
	"github.com/woshikedayaa/waveform-backend/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type LoggerWrapper struct {
	logger  *zap.Logger
	pkgName string
}

func LoggerInit() error {
	var (
		encoderConfig zapcore.EncoderConfig
		encoder       zapcore.Encoder
		writer        zapcore.WriteSyncer
		format        string = config.G().Log.Format
	)
	// 没有输出就直接return了
	if len(config.G().Log.Output) == 0 {
		return nil
	}

	// 配置 encoderConfig 和 encoder
	encoderConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
	}

	if format == "json" {

		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			//这里是 json 格式的 解析 毫秒
			enc.AppendInt64(t.UnixMilli())
		}

		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else if format == "console" {

		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			// 这个是 zap官方库里面的实现 稍微修改了一下
			layout := time.DateTime
			type appendTimeEncoder interface {
				AppendTimeLayout(time.Time, string)
			}
			if enc, ok := enc.(appendTimeEncoder); ok {
				enc.AppendTimeLayout(t, layout)
				return
			}
			enc.AppendString(t.Format(layout))
		}

		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else if format == "console_with_color" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			// 这个是 zap官方库里面的实现 稍微修改了一下
			layout := time.DateTime
			type appendTimeEncoder interface {
				AppendTimeLayout(time.Time, string)
			}
			if enc, ok := enc.(appendTimeEncoder); ok {
				enc.AppendTimeLayout(t, layout)
				return
			}
			enc.AppendString(t.Format(layout))
		}

		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	//
	// 配置 Writer
	var writers []zapcore.WriteSyncer
	for i := 0; i < len(config.G().Log.Output); i++ {
		target := config.G().Log.Output[i]
		switch target {
		case "stdout":
			writers = append(writers, zapcore.WriteSyncer(os.Stdout))
		case "stderr":
			writers = append(writers, zapcore.WriteSyncer(os.Stderr))
		default:
			file, err := os.OpenFile(target, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
			if err != nil {
				return errors.New("logf: " + err.Error())
			}
			writers = append(writers, zapcore.WriteSyncer(file))
		}
	}
	writer = zapcore.NewMultiWriteSyncer(writers...)
	// 配置 core

	level, err := zapcore.ParseLevel(config.G().Log.Level)
	if err != nil {
		return errors.New("logf: " + err.Error())
	}
	core = zapcore.NewCore(encoder, writer, level)

	// 这里可以添加一些自定义的 options
	// options = append(options, zap.Development())
	// 结束
	return nil
}
