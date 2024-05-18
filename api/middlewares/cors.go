// 跨域中间件

package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:     false,
		AllowOrigins:        []string{"*"},
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
