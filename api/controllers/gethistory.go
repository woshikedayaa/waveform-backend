package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/models"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
	"net/http"
)

// GetHistory
// @Summary 获取历史波形数据
// @Description 该接口用于获取保存过的波形数据
// @Accept json
// @Produce json
// @Router /view/gethistory [get]
func GetHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var waves []models.SaveWave

		// 执行查询操作
		if err := config.DB().Raw("SELECT id, name, CAST(timestamp AS DATETIME) as timestamp, data FROM save_waves ORDER BY id DESC LIMIT 30").Scan(&waves).Error; err != nil {
			// 记录详细错误信息
			c.JSON(http.StatusInternalServerError, resp.Error("Failed to retrieve data"))
			// 日志分级还没写完，先用这个代替查看错误
			// print(err)
			return
		}
		c.JSON(http.StatusOK, waves)
	}
}
