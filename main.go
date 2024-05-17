package main

import (
	"github.com/woshikedayaa/waveform-backend/api"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/logf"
	"os"
	"strconv"
	"strings"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = logf.LoggerInit()
	if err != nil {
		panic(err)
	}
	logger := logf.Open("main")
	// 检查是不是有config.yaml配置文件 没有就抛出一个warn
	_, err = os.Stat(config.GetDefaultConfigFilePath())
	if os.IsNotExist(err) {
		logger.Warn("无法找到config.yaml，将使用默认配置。请参考文档进行配置")
	}
	// 配置路由
	router := api.InitRouter()

	// start!
	if err = router.Run(strings.Join(
		[]string{
			config.G().Server.Addr,
			strconv.Itoa(config.G().Server.Port),
		}, ":")); err != nil {
		panic(err)
	}
}
