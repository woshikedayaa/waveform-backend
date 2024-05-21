package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/logf"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// WebSocket 升级配置
var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second, // 超时时间为10秒
	//读写缓冲区1024字节
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有源的连接
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理 WebSocket连接
func HandleWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建 logger
		logger := logf.Open("WS")
		// 升级 HTTP 连接为 WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, resp.Error("Failed to Upgrader WebSocket"))
			return
		}
		// 函数结束时关闭连接
		defer conn.Close()

		// 启动协程定时发送接收到的硬件数据
		go func() {
			for {
				time.Sleep(500 * time.Millisecond) // 发送频率
				// 初始化data为字节切片
				data := []byte{}
				// 获取缓冲区中的数据
				mu.Lock()
				if len(wsDataBuffer) == 0 {
					mu.Unlock()
					continue
				}
				for _, entry := range wsDataBuffer {
					data = append(data, entry.Data...)
				}
				// 清空dataBuffer
				wsDataBuffer = []struct {
					Timestamp time.Time
					Data      []byte
				}{}
				mu.Unlock()

				// 通过 WebSocket 发送数据到前端
				err := conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					logger.Error("Failed to send message:", zap.Error(err))
					break
				}
			}
		}()

		for {
			// 读取前端发送的消息
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Error("Failed to send message:", zap.Error(err))
				break
			}
			logger.Info("收到消息为:", zap.String("message", string(msg)))

			// 在这里处理接收到的消息，可以进行反射等操作

			// 回复客户端收到的消息
			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.Error("Failed to send message:", zap.Error(err))
				break
			}
		}
	}
}
