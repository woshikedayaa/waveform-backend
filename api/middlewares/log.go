// Logger 中间件

package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/logf"
	"go.uber.org/zap"
	"time"
)

// Logging 日志中间件函数
func Logging() gin.HandlerFunc {
	// 初始化一个用于 HTTP 请求的日志实例
	logger := logf.Open("HTTP")
	// 记录请求的基本信息
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		// handle
		c.Next()

		// 记录请求耗时
		latency := time.Since(start)
		// 记录状态码
		status := c.Writer.Status()
		// 记录内部错误
		internalErr := c.Errors.ByType(gin.ErrorTypePrivate).String()
		// 记录客户端IP
		client := c.ClientIP()

		// 根据是否存在内部错误，使用Zap记录不同级别的日志。
		// 如果有内部错误，则记录错误日志；否则，记录信息日志。
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
