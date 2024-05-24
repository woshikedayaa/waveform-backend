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
	// 配置数据库
	err = dao.InitDataBase()
	if err != nil {
		logger.Fatal("can not init database ", zap.Error(err))
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
	// 这里就不用zap的 zap.string 了
	logger.Warn("HTTP 服务器将会启动在 " + ipPort)
	if err = router.Run(ipPort); err != nil {
		panic(err)
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
