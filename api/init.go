package api

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/middlewares"
)

// 这个文件取代 router 包

func InitRouter() *gin.Engine {
	engine := gin.New()

	engine.Use(
		middlewares.Cors(),
		gin.Recovery(),
		middlewares.Logging(),
	)

	return engine
}
