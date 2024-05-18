package logf

import (
	"github.com/woshikedayaa/waveform-backend/config"
	"testing"
)

func TestLogfToConsole(t *testing.T) {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = LoggerInit()
	if err != nil {
		panic(err)
	}
	logger := Open("test")
	logger.Info("hello,info")
	logger.Warn("hello,warn")
	logger.Error("hello,error")
}
