package services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/logf"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
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
const WebsocketMaxChannelBuffer = 64

type WsTLV struct {
	MessageType int
	Length      int
	Data        []byte
}

// WSWrapper 对websocket的简单封装 写 service 只用关心数据
type WSWrapper struct {
	// must order
	logger  *zap.Logger
	conn    *websocket.Conn
	timeout time.Duration

	// channel max=WebsocketMaxChannelBuffer
	readChan  chan WsTLV
	WriteChan chan WsTLV

	// do not edit
	id        int64
	closed    bool
	failCount int
	// lock
	*sync.RWMutex
}

func (w *WSWrapper) Closed() bool {
	w.RLock()
	defer w.RUnlock()
	return w.closed
}

func (w *WSWrapper) WriteReadAble() bool {
	return !w.Closed()
}

func (w *WSWrapper) Close() {
	w.Lock()
	defer w.Unlock()
	if w.closed {
		return
	}
	w.logger.Debug("websocket 被手动关闭", zap.Int64("ID", w.id))

	// 真实的close代码
	_ = w.conn.Close()
	w.closed = true
	return
}

func (w *WSWrapper) ReadChan() <-chan WsTLV {
	return w.readChan
}

// todo 封装写和读的方法 实现自动处理超时

func (w *WSWrapper) WriteJson(obj any) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return w.WriteText(bytes)
}

func (w *WSWrapper) WriteText(data []byte) error {
	return w.write(websocket.TextMessage, data)
}

func (w *WSWrapper) write(typ int, data []byte) error {
	w.WriteChan <- WsTLV{
		MessageType: typ,
		Length:      len(data),
		Data:        data,
	}
	return w.conn.SetWriteDeadline(time.Now().Add(w.timeout))
}

func (w *WSWrapper) Read() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WSWrapper) Ping() error {
	_ = w.conn.SetReadDeadline(time.Now().Add(w.timeout))
	return w.write(websocket.PingMessage, nil)
}

func (w *WSWrapper) Error(msg string, err error, field ...zap.Field) {
	w.failCount++
	field = append(field, zap.Int64("ID", w.id), zap.Error(err))
	w.logger.Error(msg, field...)
}

func (w *WSWrapper) Warn(msg string, filed ...zap.Field) {
	filed = append(filed, zap.Int64("ID", w.id))
	w.logger.Warn(msg, filed...)
}

//

func (w *WSWrapper) Serve() {
	w.id = time.Now().UnixMilli()
	if w.logger == nil {
		w.logger = zap.NewNop()
	}
	if w.conn == nil {
		w.Warn("websocket 连接是一个空指针 nil")
		return
	}
	if w.timeout <= 0 {
		w.Warn("超时时间没有设置 设置为默认 10s")
		w.timeout = 10 * time.Second
	}

	w.conn.SetPongHandler(func(appData string) error {
		_ = w.conn.SetReadDeadline(time.Now().Add(w.timeout))
		w.logger.Debug("pong", zap.Int64("ID", w.id))
		return nil
	})

	w.conn.SetCloseHandler(func(code int, text string) error {
		w.Lock()
		defer w.Unlock()
		if w.closed {
			return nil
		}
		// 这里通知上层这个连接已经 close 了
		w.closed = true
		message := websocket.FormatCloseMessage(code, "")
		// 这里相较于官方方法 添加了自己定义的时间
		if err := w.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(w.timeout)); err != nil {
			return err
		}
		return nil
	})

	// read
	go func() {
		var (
			err error
			tlv = WsTLV{}
		)
		for w.WriteReadAble() {
			_ = w.conn.SetReadDeadline(time.Now().Add(w.timeout))
			tlv.MessageType, tlv.Data, err = w.Read()
			tlv.Length = len(tlv.Data)

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					w.Error("读取消息时发生意外关闭错误", err)
					return
				} else {
					w.Warn("读取消息失败", zap.Error(err))
				}
			}
			// 这里检查一下是不是过多没处理 如果没处理并且达到了 管道的上限 就丢弃
			if len(w.WriteChan) >= WebsocketMaxChannelBuffer {
				w.Warn(fmt.Sprintf("消息过多，未处理 max= %d", WebsocketMaxChannelBuffer))
				continue
			}
			if tlv.MessageType != websocket.PingMessage && tlv.MessageType != websocket.PongMessage {
				w.readChan <- tlv
			}
		}
	}()

	// 每一秒 ping 一下
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	// 这个变量用来检查ping是否收到 如果长时间没受到pong 就断开连接
	w.logger.Debug("开始处理websocket 连接", zap.Int64("ID", w.id))
	for w.WriteReadAble() {
		select {
		// ping
		case _ = <-ticker.C:
			if w.closed {
				continue
			}
			// 做一下健康检查
			if w.failCount >= WebsocketMaxFailCount {
				w.Error("错误次数达到最大次数，主动断开连接", nil)
				w.Close()
				// 这里是还继续的循环 因为close需要走这个循环
				continue
			}
			err := w.Ping()

			if err != nil {
				w.failCount++
				w.Warn("ping失败 将尝试再次ping ", zap.Error(err))
			}
			w.logger.Debug("ping", zap.Int64("ID", w.id))
		case tlv := <-w.WriteChan:
			err := w.conn.WriteMessage(tlv.MessageType, tlv.Data)
			if err != nil {
				w.Error("write 失败", err)
			}
			if websocket.IsUnexpectedCloseError(err) {
				w.logger.Info("发生未知的关闭错误，关闭连接", zap.Int64("ID", w.id))
				w.Close()
				return
			}
		}
	}
}

func handleWs(conn *websocket.Conn, timeout time.Duration) *WSWrapper {
	w := &WSWrapper{
		logger:    logf.Open("service/ws"),
		conn:      conn,
		timeout:   timeout,
		RWMutex:   new(sync.RWMutex),
		readChan:  make(chan WsTLV, WebsocketMaxChannelBuffer),
		WriteChan: make(chan WsTLV, WebsocketMaxChannelBuffer),
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
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for w.WriteReadAble() {
		select {
		case _ = <-ticker.C:
			// todo 保存到全局变量 方便保存 （可能有）
			//
			// 这里只是测试用 生成随机的数据
			points, _ := WaveForm.ConvDataToPoints(wave.GenerateRandomData(1024), 1, 1024)
			//
			err := w.WriteJson(resp.Success(points))
			if err != nil {
				w.Error("通过 websocket 写入数据的时候出现错误", err)
				continue
			}
		}
	}
	//
}
