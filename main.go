package main

import (
	"fmt"
	"github.com/woshikedayaa/waveform-backend/api"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/dao"
	"github.com/woshikedayaa/waveform-backend/logf"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

// Version 在构建的时候注入
var Version string

func main() {
	if len(os.Args) > 1 {
		job := os.Args[1]
		DoJob(job)
		return
	}
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = logf.LoggerInit()
	if err != nil {
		panic(err)
	}

	logger := logf.Open("main")
	// 配置数据库
	err = dao.InitDataBase()
	if err != nil {
		logger.Fatal("无法初始化数据库", zap.Error(err))
		return
	}
	// 配置路由
	router := api.InitRouter()

	// start!
	ipPort := strings.Join(
		[]string{
			config.G().Server.Http.Addr,
			strconv.Itoa(config.G().Server.Http.Port),
		}, ":")
	logger.Sugar().Infof("HTTP 服务器将启动在 %s", ipPort)
	if err = router.Run(ipPort); err != nil {
		logger.Fatal("http服务器启动失败", zap.Error(err))
	}
}

func DoJob(job string) {
	switch job {
	case "config":
		fmt.Print(config.GetExampleConfig())
	case "version":
		fmt.Println(Version)
	default:
		fmt.Printf("不支持的操作 %s\n", job)
	}
	return
}
