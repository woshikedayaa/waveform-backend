package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/logf"
	"go.uber.org/zap"
	"time"
)

func Logging() gin.HandlerFunc {
	logger := logf.Open("HTTP")
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		// handle
		c.Next()

		//
		latency := time.Since(start)
		status := c.Writer.Status()
		internalErr := c.Errors.ByType(gin.ErrorTypePrivate).String()
		client := c.ClientIP()

		if len(internalErr) != 0 {
			logger.Error(path,
				zap.String("method", method),
				zap.Int("status", status),
				zap.Duration("latency", latency),
				zap.String("remote", client),
				zap.String("error", internalErr),
			)
		} else {
			logger.Info(path,
				zap.String("method", method),
				zap.Int("status", status),
				zap.Duration("latency", latency),
				zap.String("remote", client),
			)
		}
	}
}
