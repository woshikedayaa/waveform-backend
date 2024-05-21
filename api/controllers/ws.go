package controllers

import (
	"github.com/gorilla/websocket"

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
