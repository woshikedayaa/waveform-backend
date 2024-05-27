// 注意： 这个文件 是在生产环境启用的

//go:build deploy

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/middlewares"
	"github.com/woshikedayaa/waveform-backend/config"
)

func ginConfigure() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	//配置中间件
	engine.Use(
		gin.Recovery(),
		middlewares.Logging(),
	)
	if config.G().Server.Http.Cors.Enabled {
		engine.Use(middlewares.Cors())
	}
	return engine
}
