package api

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/controllers"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
	"net/http"
)

// InitRouter 这个文件取代 router 包
// 初始化路由函数
func InitRouter() *gin.Engine {
	engine := ginConfigure()

	// test
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, resp.Success("pong!"))
	})
	// 用于和前端交互的路由组
	viewGroup := engine.Group("/view")
	// 与前端交互的WebSocket路由
	viewGroup.GET("/ws", controllers.WebSocketController())

	// 用于保存与获取历史数据的路由组
	// historyGroup := engine.Group("/save")

	return engine
}
