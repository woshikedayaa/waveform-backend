package main

import (
	"fmt"
	"github.com/woshikedayaa/waveform-backend/api"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/logf"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	if slices.Contains(os.Args, "config") {
		fmt.Print(config.GetExampleConfig())
		return
	}

	configFileExist := true
	err := config.InitConfig()
	if err != nil {
		if os.IsNotExist(err) {
			configFileExist = false
		} else {
			panic(err)
		}
	}
	err = logf.LoggerInit()
	if err != nil {
		panic(err)
	}
	logger := logf.Open("main")
	// 检查配置文件是否存在 现在用logger打印出来提示用户
	if !configFileExist {
		logger.Warn("无法找到config.yaml，将使用默认配置。请参考文档进行配置")
	}
	// 配置路由
	router := api.InitRouter()

	// start!
	if err = router.Run(strings.Join(
		[]string{
			config.G().Server.Http.Addr,
			strconv.Itoa(config.G().Server.Http.Port),
		}, ":")); err != nil {
		panic(err)
	}
}
