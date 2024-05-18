package main

import (
	"github.com/woshikedayaa/waveform-backend/api"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/logf"
	"strconv"
	"strings"
)

func main() {
	// 初始化路由
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	// 初始化日志
	err = logf.LoggerInit()
	if err != nil {
		panic(err)
	}
	// 初始化数据库
	err = config.InitGorm()
	if err != nil {
		panic(err)
	}

	// 初始化路由
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
