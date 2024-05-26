package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
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
	// Swagger API文档路由
	// 初始化Swagger
	docs.SwaggerInfo.Title = "Waveform Backend API"
	docs.SwaggerInfo.Description = "这是一个示波器后端API的文档."
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
