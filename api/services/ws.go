package services

import (
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/logf"
	"github.com/woshikedayaa/waveform-backend/pkg/wave"
	"go.uber.org/zap"
	"sync"
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

// WSWrapper 对websocket的简单封装 写 service 只用关心数据
type WSWrapper struct {
	logger  *zap.Logger
	conn    *websocket.Conn
	timeout time.Duration

	id        int64
	closeChan chan struct{}
	closed    bool
	failCount int
	*sync.RWMutex
}

func (w *WSWrapper) Closed() bool {
	return w.closed
}

func (w *WSWrapper) WriteReadAble() bool {
	return !w.closed
}

func (w *WSWrapper) Close() {
	defer func() {
		// do nothing
		recover()
	}()
	w.logger.Debug("websocket 被手动关闭", zap.Int64("ID", w.id))
	w.closeChan <- struct{}{}
	w.closed = true
}

// todo 封装写和读的方法 实现自动处理超时

func (w *WSWrapper) WriteText(data []byte) error {
	err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout))
	if err != nil {
		return err
	}
	w.Lock()
	defer w.Unlock()
	return w.conn.WriteMessage(websocket.TextMessage, data)
}

func (w *WSWrapper) Read() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WSWrapper) Ping() error {
	err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout))
	if err != nil {
		return err
	}
	w.Lock()
	defer w.Unlock()
	return w.conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (w *WSWrapper) Error(msg string, err error) {
	w.failCount++
	w.logger.Error(msg, zap.Error(err), zap.Int64("ID", w.id))
}

func (w *WSWrapper) Warn(msg string, filed ...zap.Field) {
	w.logger.Warn(msg, filed...)
}

//

func (w *WSWrapper) Serve() {
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
	// 这个变量用来检查ping是否收到 如果长时间没受到pong 就断开连接
	pingSent := 0
	w.logger.Debug("开始处理websocket 连接", zap.Int64("ID", w.id))
	// 用来 响应 pong
	go func() {
		_, _, _ = w.Read()
		pingSent = 0
	}()
	for w.WriteReadAble() {
		select {
		// ping
		case _ = <-ticker.C:
			if w.closed {
				continue
			}
			if pingSent > WebsocketMaxFailCount {
				w.Close()
				continue
			}
			err := w.Ping()
			if err != nil {
				w.failCount++
				w.Warn("ping失败 将尝试再次ping ", zap.Error(err))
			}
			if w.failCount >= WebsocketMaxFailCount {
				w.Error("错误次数达到最大次数，主动断开连接", nil)
				w.Close()
				// 这里是还继续的循环 因为close需要走这个循环
				continue
			}
			w.logger.Debug("ping", zap.Int64("ID", w.id))
		// 关闭
		case _ = <-w.closeChan:
			close(w.closeChan)
			w.Lock()
			_ = w.conn.Close()
			w.Unlock()
			ticker.Stop()
			w.logger.Debug("websocket 连接被远端或者手动关闭", zap.Int64("ID", w.id))
			return
		}
	}
}

func handleWs(conn *websocket.Conn, timeout time.Duration) *WSWrapper {
	w := &WSWrapper{
		logger:  logf.Open("service/ws"),
		conn:    conn,
		timeout: timeout,
		RWMutex: new(sync.RWMutex),
	}
	go w.Serve()
	return w
}

type ws struct{}

var WebSocket ws

// HandleWebsocketForWaveform 处理来自前端的websocket 处理波形图的
func (ws) HandleWebsocketForWaveform(conn *websocket.Conn, timeout time.Duration) {
	w := handleWs(conn, timeout)
	defer w.Close()
	ticker := time.NewTicker(time.Second)
	for w.WriteReadAble() {
		select {
		case _ = <-ticker.C:
			// 这里只是测试用 生成随机的数据
			data := wave.GenerateRandomData(1024)
			// todo 保存到全局变量 方便保存 （可能有）
			//
			err := w.WriteText(data)
			if err != nil {
				w.Error("通过 websocket 写入数据的时候出现错误", err)
				continue
			}
		default:
			// do nothing
		}
	}
	//
}
