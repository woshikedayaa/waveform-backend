package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/services"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
	"net/http"
	"strconv"
)

func GetWaveFormByWebsocket() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetWaveFromByHttp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err     error
			samples = c.DefaultQuery("sample", "4096")
			counts  = c.DefaultQuery("count", "1024")
			count   = 1024
			sample  = 4096

			points []services.Point
		)
		count, err = strconv.Atoi(counts)
		if err != nil {
			c.JSON(http.StatusOK, resp.Fail("count not is a number or it is too bigger"))
			return
		}
		sample, err = strconv.Atoi(samples)
		if err != nil {
			c.JSON(http.StatusOK, resp.Fail("sample not is a number or it is too bigger"))
			return
		}

		points, err = services.GetLatestWave(sample, count)
		if err != nil {
			c.JSON(http.StatusOK, resp.Error(err.Error()))
			return
		}
		c.JSON(http.StatusOK, resp.Success(points))
		return
	}
}
