package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
	"net/http"
)

func SaveTemporary() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.DefaultQuery("ts", "")
		if len(id) == 0 {
			c.JSON(http.StatusOK, resp.Fail("必须指定一个时间戳"))
			return
		}
		// todo
	}
}
