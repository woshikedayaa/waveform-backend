package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/services"
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
		)
		count, err = strconv.Atoi(samples)
		if err != nil {

		}

		services.GetLatestWave(sample, count)
	}
}
