package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/services"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
)

// ReceiveHardwareDataController
// @Summary      接收来自硬件的数据
// @Description  该接口用于接收来自硬件的二进制数据
// @Accept       multipart/form-data
// @Success      200  {string}  "Started receiving data"
// @Router      /hardware/receive [post]
func ReceiveHardwareDataController() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := services.ReceiveHardwareData()
		if err != nil {
			c.JSON(http.StatusBadGateway, resp.Error(fmt.Sprintf("failed to connect to kcp server err: %s", err)))
			return
		}
		c.JSON(http.StatusOK, resp.Success("started receiving data"))
	}
}
