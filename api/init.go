package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/woshikedayaa/waveform-backend/api/controllers"
	"github.com/woshikedayaa/waveform-backend/api/middlewares"
	"github.com/woshikedayaa/waveform-backend/docs"
)

// 这个文件取代 router 包
// 初始化路由函数
func InitRouter() *gin.Engine {
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

	// 用于接收硬件数据的路由组
	hardwareGroup := engine.Group("/hardware")
	// 接收硬件数据的路由
	hardwareGroup.GET("/receive", controllers.ReceiveHardwareData())

	// 用于和前端交互的路由组
	viewGroup := engine.Group("/view")
	// 与前端交互的WebSocket路由
	viewGroup.GET("/ws", controllers.HandleWebSocket())
	// 处理波形图数据保存的路由
	viewGroup.POST("/save", controllers.SaveWave())
	// 获取历史波形记录
	viewGroup.GET("/history", controllers.GetHistory())

	// 暂时用不到（暂留）
	apiGroup := engine.Group("/api")
	apiGroup.GET("/latest", controllers.GetWaveFromByHttp())

	return engine
}
