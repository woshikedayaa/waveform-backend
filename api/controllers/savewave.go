package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/models"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
	"net/http"
	"sync"
	"time"
)

var clearTimer *time.Timer
var clearTimerMutex sync.Mutex

// 接收客户端发送的JSON数据
type SaveRequest struct {
	Name string `json:"name" binding:"required"`
}

// SaveWave
// @Summary 保存波形数据
// @Description 该接口用于保存固定时间内的波形数据
// @Accept json
// @Produce json
// @Router /view/save [post]
func SaveWave() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SaveRequest
		// 解析JSON数据并绑定到SaveRequest类型的变量req上
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, resp.Error("Invalid request payload"))
			return
		}
		// 计算保存波形的时间范围
		now := time.Now()
		startTime := now.Add(-1 * time.Second)
		endTime := now.Add(1 * time.Second)

		mu.Lock()
		defer mu.Unlock()

		// 筛选出时间戳在指定范围内的数据，并合并到dataToSave中
		var dataToSave []byte
		for _, entry := range saveWaveDataBuffer {
			if entry.Timestamp.After(startTime) && entry.Timestamp.Before(endTime) {
				dataToSave = append(dataToSave, entry.Data...)
			}
		}

		saveWave := models.SaveWave{
			Name:      req.Name,
			Timestamp: now.Unix(),
			Data:      dataToSave,
		}
		// 保存到数据库
		if err := config.DB().Create(&saveWave).Error; err != nil {
			c.JSON(http.StatusInternalServerError, resp.Error("Failed to save wave data"))
			return
		}

		c.JSON(http.StatusOK, resp.Success("Wave data saved successfully"))
		// 重置定时器
		resetClearTimer()
	}
}

// 重置清空缓冲区的定时器
// 在无请求的时候，定时清空数据
func resetClearTimer() {
	clearTimerMutex.Lock()
	defer clearTimerMutex.Unlock()

	if clearTimer != nil {
		clearTimer.Stop()
	}

	clearTimer = time.AfterFunc(5*time.Second, func() {
		mu.Lock()
		defer mu.Unlock()
		saveWaveDataBuffer = []struct {
			Timestamp time.Time
			Data      []byte
		}{}
	})
}
