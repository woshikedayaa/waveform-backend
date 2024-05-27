// 注意： 这个文件 是在开发环境启用的

//go:build !deploy

package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/woshikedayaa/waveform-backend/api/middlewares"
	"github.com/woshikedayaa/waveform-backend/api/swag"
	"github.com/woshikedayaa/waveform-backend/config"
)

func ginConfigure() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()

	//配置中间件
	engine.Use(
		gin.Recovery(),
		middlewares.Logging(),
	)
	if config.G().Server.Http.Cors.Enabled {
		engine.Use(middlewares.Cors())
	}

	// Swagger API文档路由
	// 初始化Swagger
	swag.SwaggerInfo.Title = "Waveform Backend API"
	swag.SwaggerInfo.Description = "这是一个示波器后端API的文档."
	swag.SwaggerInfo.Version = "0.1"
	swag.SwaggerInfo.Host = "localhost:8080"
	swag.SwaggerInfo.BasePath = "/"

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return engine
}
