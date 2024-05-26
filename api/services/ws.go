package services

import (
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/logf"
	"go.uber.org/zap"
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

const WebsocketMaxFailCount = 24

// WS 对websocket的简单封装 写 service 只用关心数据
type WS struct {
	logger  *zap.Logger
	conn    *websocket.Conn
	timeout time.Duration

	id        int64
	closeChan chan struct{}
	closed    bool
	failCount int
}

// HandleWebsocketForWaveform 处理来自前端的websocket 处理波形图的
func HandleWebsocketForWaveform(conn *websocket.Conn, timeout time.Duration) {
	w := &WS{
		logger:  logf.Open("service/ws"),
		conn:    conn,
		timeout: timeout,
	}
	go w.Serve()
	// 这里发送波形数据
	defer w.Close()
	for w.WriteReadAble() {

	}
	//
}

func (w *WS) Closed() bool {
	return w.closed
}

func (w *WS) WriteReadAble() bool {
	return !w.closed
}

func (w *WS) Close() {
	defer func() {
		// do nothing
		recover()
	}()
	w.logger.Debug("websocket 被手动关闭", zap.Int64("ID", w.id))
	w.closeChan <- struct{}{}
	w.closed = true
}

// todo 封装写和读的方法 实现自动处理超时

func (w *WS) WriteText(data []byte) error {
	err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout))
	if err != nil {
		return err
	}
	return w.conn.WriteMessage(websocket.TextMessage, data)
}

func (w *WS) Read() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WS) Ping() error {
	err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout))
	if err != nil {
		return err
	}
	return w.conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (w *WS) Serve() {
	if w.logger == nil {
		w.logger = zap.NewNop()
	}
	if w.conn == nil {
		w.logger.Warn("websocket 连接是一个空指针 nil")
		return
	}
	if w.timeout <= 0 {
		w.logger.Warn("超时时间没有设置 设置为默认 10s")
		w.timeout = 10 * time.Second
	}
	w.closed = false
	w.id = time.Now().UnixMilli()
	w.closeChan = make(chan struct{})

	w.conn.SetPongHandler(func(appData string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(w.timeout))
		w.logger.Debug("pong", zap.Int64("ID", w.id))
		return nil
	})

	w.conn.SetCloseHandler(func(code int, text string) error {
		// 防止直接调用 ws.conn.close 这里把上层的 close写入
		w.closed = true
		message := websocket.FormatCloseMessage(code, "")
		// 这里相较于官方方法 添加了自己定义的时间
		if err := w.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(w.timeout)); err != nil {
			return err
		}
		return nil
	})

	// 每一秒 ping 一下
	ticker := time.NewTicker(time.Second)
	w.logger.Debug("开始处理websocket 连接", zap.Int64("ID", w.id))
	for {
		select {
		// ping
		case _ = <-ticker.C:
			if w.closed {
				continue
			}
			err := w.Ping()
			if err != nil {
				w.failCount++
				w.logger.Error("ping 客户端发生错误,将重新尝试 ping ", zap.Int64("ID", w.id), zap.Error(err))
			}
			if w.failCount >= WebsocketMaxFailCount {
				w.logger.Error("错误次数达到最大次数 将断开连接 ", zap.Int64("ID", w.id))
				w.Close()
				// 这里是还继续的循环 因为close需要走这个循环
				continue
			}
			w.logger.Debug("ping", zap.Int64("ID", w.id))
		// 关闭
		case _ = <-w.closeChan:
			close(w.closeChan)
			_ = w.conn.Close()
			ticker.Stop()
			w.logger.Debug("websocket 连接被远端或者手动关闭", zap.Int64("ID", w.id))
			return
		}
	}
}
