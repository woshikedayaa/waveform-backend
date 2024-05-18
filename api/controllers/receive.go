package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
	"github.com/xtaci/kcp-go/v5"
	"log"
	"net/http"
	"sync"
	"time"
)

// 存储来自示波器硬件的数据 -> ws发送到前端
var wsDataBuffer []struct {
	Timestamp time.Time
	Data      []byte
}

// 存储来自示波器硬件的数据 -> savewave保存到数据库
var saveWaveDataBuffer []struct {
	Timestamp time.Time
	Data      []byte
}

var mu sync.Mutex

// ReceiveHardwareData
// @Summary      接收来自硬件的数据
// @Description  该接口用于接收来自硬件的二进制数据
// @Accept       multipart/form-data
// @Success      200  {string}  "Started receiving data"
// @Router      /hardware/receive [get]
func ReceiveHardwareData() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取配置信息
		kcpConfig := config.G().Kcp

		// 创建 KCP 会话
		sess, err := kcp.DialWithOptions(kcpConfig.ListenAddress, nil, kcpConfig.Sndwnd, kcpConfig.Rcvwnd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, resp.Error("Failed to connect to KCP server"))
			return
		}
		// 函数结束时关闭连接
		defer sess.Close()

		// 配置 KCP 参数
		sess.SetWindowSize(kcpConfig.Sndwnd, kcpConfig.Rcvwnd)
		sess.SetMtu(kcpConfig.Mtu)
		sess.SetNoDelay(kcpConfig.Nodelay, kcpConfig.Interval, kcpConfig.Resend, kcpConfig.Nc)

		// 启动协程异步读取数据
		go func() {
			// 分配内存
			buf := make([]byte, 4096)
			for {
				n, err := sess.Read(buf)
				if err != nil {
					log.Println("Failed to read from KCP server:", err)
					return
				}

				// 保存数据到 buffer
				mu.Lock()
				entry := struct {
					Timestamp time.Time
					Data      []byte
				}{
					Timestamp: time.Now(),
					Data:      append([]byte(nil), buf[:n]...),
				}
				wsDataBuffer = append(wsDataBuffer, entry)
				saveWaveDataBuffer = append(saveWaveDataBuffer, entry)
				mu.Unlock()
			}
		}()

		c.JSON(http.StatusOK, resp.Success("Started receiving data"))
	}
}
