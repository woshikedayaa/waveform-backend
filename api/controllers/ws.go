package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/api/services"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
	"net/http"
	"time"
)

// WebSocket 升级配置
var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second, // 超时时间为10秒
	//读写缓冲区 1024 字节
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有源的连接
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketController 处理 WebSocket 连接
func WebSocketController() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 升级 HTTP 连接为 WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, resp.Error(fmt.Sprintf("无法成功将此连接升级为 websocket %s", err)))
			return
		}
		go services.WebSocket.HandleWebsocketForWaveform(conn, 10*time.Second)
	}
}
