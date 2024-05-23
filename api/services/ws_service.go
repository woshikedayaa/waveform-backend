package services

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/logf"
	"go.uber.org/zap"
	"time"
)

// SendWebSocketData 处理 WebSocket 连接的业务逻辑
func SendWebSocketData(conn *websocket.Conn) {
	// 初始化 logger
	logger := logf.Open("WsSendData")
	// 定义数据发送频率（0.5s）
	ticker := time.NewTicker(500 * time.Millisecond)
	// 函数结束时停止ticker
	defer ticker.Stop()
	// 函数结束时关闭连接
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Error closing connection: %v", zap.Error(err))
		}
	}()
	for {
		select {
		case <-ticker.C:
			mu.Lock()
			if len(dataBuffer) == 0 {
				mu.Unlock()
				continue
			}
			// 如果 dataBuffer 中有数据，向前端发送，并清空缓存
			data := make([]byte, len(dataBuffer))
			copy(data, dataBuffer)
			dataBuffer = []byte{}
			mu.Unlock()

			err := conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				logger.Error("Failed to write to WebSocket:", zap.Error(err))
				return
			}
		}
	}
}

// ReceiveWebSocketMessage 接收前端发送的消息
func ReceiveWebSocketMessage(conn *websocket.Conn) {
	logger := logf.Open("WsReceiveMessage")
	// 函数结束时关闭连接
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Error closing connection: %v", zap.Error(err))
		}
	}()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Error("Error reading from WebSocket:", zap.Error(err))
			return
		}
		logger.Info(fmt.Sprintf("Received message: %s", msg))
		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			logger.Error("Error writing to WebSocket:", zap.Error(err))
			return
		}
	}
}
