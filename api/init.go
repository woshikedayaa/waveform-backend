package api

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/controllers"
	"github.com/woshikedayaa/waveform-backend/api/middlewares"
)

// InitRouter 这个文件取代 router 包
// 初始化路由函数
func InitRouter() *gin.Engine {
	ginConfigure()

	engine := gin.New()

	//配置中间件
	engine.Use(
		middlewares.Cors(),
		gin.Recovery(),
		middlewares.Logging(),
	)

	// 用于接收硬件数据的路由组
	// hardwareGroup := engine.Group("/hardware")

	// 用于和前端交互的路由组

	// 与前端交互的WebSocket路由

	// 处理波形图数据保存的路由

	// 暂时用不到（暂留）
	apiGroup := engine.Group("/api")
	apiGroup.GET("/latest", controllers.GetWaveFromByHttp())

	return engine
}
