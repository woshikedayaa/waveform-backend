package api

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/controllers"
)

// InitRouter 这个文件取代 router 包
// 初始化路由函数
func InitRouter() *gin.Engine {
	engine := ginConfigure()

	// test
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "this is ok!",
		})
	})
	// 用于和前端交互的路由组
	viewGroup := engine.Group("/view")
	// 与前端交互的WebSocket路由
	viewGroup.GET("/ws", controllers.WebSocketController())

	// 暂时用不到（暂留）
	apiGroup := engine.Group("/api")
	apiGroup.GET("/latest", controllers.GetWaveFromByHttp())

	return engine
}
