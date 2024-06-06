package services

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/pkg/wave"
	"github.com/woshikedayaa/waveform-backend/pkg/ws"
	"go.uber.org/zap"
	"time"
)

type webSocket struct{}

var WebSocket webSocket

// HandleWebsocketForWaveform 处理来自前端的websocket 处理波形图的
func (webSocket) HandleWebsocketForWaveform(conn *websocket.Conn, timeout time.Duration) {
	var (
		w           = ws.HandleWs(conn, timeout)
		ticker      = time.NewTicker(time.Second)
		offset      int
		buffer      []*wave.FullData
		readyToSave []*wave.FullData
	)
	defer w.Close()
	defer ticker.Stop()
	for w.WriteReadAble() {
		select {
		case _ = <-ticker.C:
			// todo 保存到全局变量 方便保存 （可能有）
			//
			// 这里只是测试用 生成随机的数据
			f := wave.ParseRawData(wave.RandomData(1024), 1, 1024)
			//
			buffer = append(buffer, f)

			err := w.WriteJson(gin.H{
				"data": f,
				"idx":  offset + len(buffer),
			})
			if err != nil {
				w.Error("通过 websocket 写入数据的时候出现错误", err, zap.Int("data_id", offset+len(buffer)))
				return
			}
		case r := <-w.ReadChan():
			if r.MessageType == websocket.TextMessage {
				const (
					ActionSave = iota
				)
				const (
					TypeSaveTemporary = iota
					TypeSavePersistent
				)
				type Data struct {
					Action int            `json:"action"`
					Typ    int            `json:"typ"`
					Arg    map[string]any `json:"arg"`
				}
				data := &Data{}
				err := json.Unmarshal(r.Data, data)
				if err != nil {
					w.Error("解析来自前端的数据出现错误", err)
					return
				}
				switch data.Action {
				case ActionSave:
					idx, ok := data.Arg["idx"].(int)
					if !ok {
						continue
					}
					// 这里就删除现在的buffer了 释放内存 只保留已经决定了要临时保留的
					if data.Typ == TypeSaveTemporary {
						if len(buffer) <= idx-offset {
							continue
						}
						readyToSave = append(readyToSave, buffer[idx-offset])
						// 防止内存泄漏
						for i := 0; i < idx-offset; i++ {
							buffer[i] = nil
						}
						//
						buffer = append([]*wave.FullData{}, buffer[idx-offset+1:]...)
						offset = idx + 1
					} else if data.Typ == TypeSavePersistent {
						// 这里保存到数据库
					} else {
						// do nothing
					}
				}
			}
			//
		}
	}
	//
}
