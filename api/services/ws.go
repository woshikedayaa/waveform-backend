package services

import (
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/pkg/wave"
	"github.com/woshikedayaa/waveform-backend/pkg/ws"
	"time"
)

//
//// SendWebSocketData 处理 WebSocket 连接的业务逻辑
//func SendWebSocketData(conn *websocket.Conn) {
//	// 初始化 logger
//	logger := logf.Open("WsSendData")
//	// 定义数据发送频率（0.5s）
//	ticker := time.NewTicker(500 * time.Millisecond)
//	// 函数结束时停止ticker
//	defer ticker.Stop()
//	// 函数结束时关闭连接
//	defer func() {
//		if err := conn.Close(); err != nil {
//			logger.Error("Error closing connection", zap.Error(err))
//		}
//	}()
//	for {
//		select {
//		case <-ticker.C:
//			mu.Lock()
//			if len(dataBuffer) == 0 {
//				mu.Unlock()
//				continue
//			}
//			// 如果 dataBuffer 中有数据，向前端发送，并清空缓存
//			data := make([]byte, len(dataBuffer))
//			copy(data, dataBuffer)
//			dataBuffer = []byte{}
//			mu.Unlock()
//
//			err := conn.WriteMessage(websocket.BinaryMessage, data)
//			if err != nil {
//				logger.Error("Failed to write to WebSocket", zap.Error(err))
//				return
//			}
//		}
//	}
//}
//
//// ReceiveWebSocketMessage 接收前端发送的消息
//func ReceiveWebSocketMessage(conn *websocket.Conn) {
//	logger := logf.Open("WsReceiveMessage")
//	// 函数结束时关闭连接
//	defer func() {
//		if err := conn.Close(); err != nil {
//			logger.Error("Error closing connection: %v", zap.Error(err))
//		}
//	}()
//	for {
//		_, msg, err := conn.ReadMessage()
//		if err != nil {
//			logger.Error("Error reading from WebSocket:", zap.Error(err))
//			return
//		}
//		logger.Debug(fmt.Sprintf("Received message: %s", msg))
//
//		err = conn.WriteMessage(websocket.TextMessage, msg)
//		if err != nil {
//			logger.Error("Error writing to WebSocket:", zap.Error(err))
//			return
//		}
//	}
//}

type webSocket struct{}

var WebSocket webSocket

// HandleWebsocketForWaveform 处理来自前端的websocket 处理波形图的
func (webSocket) HandleWebsocketForWaveform(conn *websocket.Conn, timeout time.Duration) {

	w := ws.HandleWs(conn, timeout)
	defer w.Close()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for w.WriteReadAble() {
		select {
		case _ = <-ticker.C:
			// todo 保存到全局变量 方便保存 （可能有）
			//
			// 这里只是测试用 生成随机的数据
			f := wave.ParseRawData(wave.RandomData(1024), 1, 1024)
			//
			err := w.WriteJson(f)
			if err != nil {
				w.Error("通过 websocket 写入数据的时候出现错误", err)
				w.Close()
				return
			}
		case r := <-w.ReadChan():
			_ = r
		}
	}
	//
}
