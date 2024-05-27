package services

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestWS_Serve(t *testing.T) {
	// a simple echo server
	cw := make(chan *WSWrapper)
	go runWsServer("/", "8080", cw)

	for ws := range cw {
		go ws.Serve()
		go func() {
			defer ws.Close()
			for ws.WriteReadAble() {
				messageType, data, err := ws.Read()
				if err != nil {
					var ce *websocket.CloseError
					if errors.As(err, &ce) && ce.Code == websocket.CloseNormalClosure {
						break
					} else {
						ws.logger.Error("read", zap.Error(err))
					}
					break
				}
				if messageType != websocket.TextMessage {
					continue
				}
				ws.logger.Info("read", zap.String("data", string(data)))

				err = ws.WriteText(data)
				if err != nil {
					ws.logger.Error("write", zap.Error(err))
				}
			}
		}()
	}

}

func runWsServer(path string, port string, cw chan *WSWrapper) {
	engine := gin.New()
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	engine.GET(path, func(c *gin.Context) {
		conn, err := up.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}
		cw <- &WSWrapper{
			logger:  zap.NewExample(),
			conn:    conn,
			timeout: 10 * time.Second,
		}

	})

	engine.Run(":" + port)
}
