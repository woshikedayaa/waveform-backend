package services

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/woshikedayaa/waveform-backend/dao"
	"github.com/woshikedayaa/waveform-backend/logf"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
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
		w           = ws.HandleWs(conn, timeout, logf.Open("service/ws"))
		ticker      = time.NewTicker(time.Second)
		offset      int
		buffer      []*wave.FullData
		readyToSave []*wave.FullData
	)
	defer w.Close()
	defer ticker.Stop()
	defer func() {
		// 防止内存泄漏
		for i := 0; i < len(buffer); i++ {
			buffer[i] = nil
		}
		for i := 0; i < len(readyToSave); i++ {
			readyToSave[i] = nil
		}
	}()

	for w.WriteReadAble() {
		select {
		case _ = <-ticker.C:
			// todo 保存到全局变量 方便保存 （可能有）
			//
			// 这里只是测试用 生成随机的数据
			f := wave.ParseRawData(wave.RandomData(1024), 1, 1024)
			// 防止太大了
			if len(buffer) > 200 {
				buffer[0] = nil
				buffer = buffer[1:]
			}
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
					Arg    map[string]any `json:"args"`
				}
				data := &Data{}
				err := json.Unmarshal(r.Value[:r.Length], data)
				if err != nil {
					w.Error("解析来自前端的数据出现错误", err)
					return
				}
				// 开始解析
				switch data.Action {
				case ActionSave:
					if data.Typ == TypeSaveTemporary {
						idxf, ok := data.Arg["idx"].(float64)
						if !ok {
							continue
						}
						idx := int(idxf)
						// 这里就删除现在的buffer了 释放内存 只保留已经决定了要临时保留的
						if idx-offset < 0 || len(buffer) <= idx-offset {
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
						// 这里就退出这个 ws 连接了 这一轮的ws已经保存
						names, ok := data.Arg["names"].([]any)
						if !ok {
							return
						}
						idxes, ok := data.Arg["idxes"].([]any)
						if !ok {
							return
						}
						if len(idxes) != len(names) {
							err := w.WriteJson(resp.Fail("要保存的数量不符合,names 应该和 idxes 数量相同"))
							if err != nil && websocket.IsUnexpectedCloseError(err,
								websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
								w.Error("保存到数据库时发送发生错误，在返回至前端发送错误", err)
								return
							}
						}
						if len(idxes) == 0 {
							return
						}
						if len(idxes) < len(readyToSave) {
							w.Warn("持久化保存的数量多于临时保存的数量，将不保存")
							return
						}

						db := dao.NewWaveFormDao()
						for i := 0; i < len(names); i++ {
							// 这里保存到数据库
							name, ok := names[i].(string)
							if !ok {
								continue
							}
							idx, ok := idxes[i].(float64)
							if !ok {
								continue
							}
							err = db.Save(context.Background(), name, readyToSave[int(idx)])
							if err != nil {
								w.Error("保存到数据库时发送发生错误", err)
							}
						}
						return
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
