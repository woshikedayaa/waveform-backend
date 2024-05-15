package logf

import (
	"errors"
	"fmt"
	"github.com/woshikedayaa/waveform-backend/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LoggerWrapper struct {
	logger  *zap.Logger
	pkgName string
}

func LoggerInit() error {
	var (
		err error

		encoderConfig zapcore.EncoderConfig
		encoder       zapcore.Encoder
		writer        zapcore.WriteSyncer
		format        string = strings.ToLower(config.G().Log.Format)
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
	} else {
		return errors.New(fmt.Sprintf("logf: Unsupported format %s", format))
	}
	//
	// 配置 Writer

	writer, err = getWriter(config.G().Log.Output)
	// 配置 core

	level, err := zapcore.ParseLevel(config.G().Log.Level)
	if err != nil {
		return errors.New("logf: " + err.Error())
	}
	core = zapcore.NewCore(encoder, writer, level)

	// 配置 options
	err = optionsInit()
	if err != nil {
		return errors.New("logf: " + err.Error())
	}

	// 结束
	return nil
}

func getWriter(ss []string) (zapcore.WriteSyncer, error) {
	var writers []zapcore.WriteSyncer
	for i := 0; i < len(ss); i++ {
		target := ss[i]
		switch target {
		case "stdout":
			writers = append(writers, zapcore.WriteSyncer(os.Stdout))
		case "stderr":
			writers = append(writers, zapcore.WriteSyncer(os.Stderr))
		default:
			dir := filepath.Dir(target)
			err := os.MkdirAll(dir, 0755) // rwxrw-rw-
			if err != nil {
				return nil, errors.New("logf: " + err.Error())
			}

			file, err := os.OpenFile(target, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644) // rw-r--r--
			if err != nil {
				return nil, errors.New("logf: " + err.Error())
			}
			writers = append(writers, zapcore.WriteSyncer(file))
		}
	}
	return zapcore.NewMultiWriteSyncer(writers...), nil
}

func optionsInit() error {
	// error output
	ew, err := getWriter(config.G().Log.ErrOutput)
	if err != nil {
		return err
	}
	options = append(options, zap.ErrorOutput(ew))
	return nil
}
