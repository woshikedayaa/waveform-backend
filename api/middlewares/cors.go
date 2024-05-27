// 跨域中间件

package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/config"
	"time"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:     false,
		AllowOrigins:        config.G().Server.Http.Cors.Origins,
		AllowMethods:        []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowPrivateNetwork: true,
		AllowHeaders:        []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials:    true,
		MaxAge:              12 * time.Hour,
		AllowWildcard:       true,
		AllowWebSockets:     true,
		AllowFiles:          true,
	})
}
