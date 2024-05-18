package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woshikedayaa/waveform-backend/api/services"
	"github.com/woshikedayaa/waveform-backend/pkg/resp"
)

type WaveFormParam struct {
	sample int
	count  int
}

func (w *WaveFormParam) GetFromUrl(c *gin.Context) error {
	var (
		err     error
		samples = c.DefaultQuery("sample", "4096")
		counts  = c.DefaultQuery("count", "1024")
	)
	w.count, err = strconv.Atoi(counts)
	if err != nil {
		return errors.New("count not is a number or it is too bigger")
	}
	w.sample, err = strconv.Atoi(samples)
	if err != nil {
		return errors.New("sample not is a number or it is too bigger")
	}
	return nil
}

func GetWaveFormByWebsocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
			wp  = &WaveFormParam{}
		)
		err = wp.GetFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, resp.Fail(err.Error()))
			return
		}

	}
}

func GetWaveFromByHttp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err    error
			wp     = &WaveFormParam{}
			points []services.Point
		)
		err = wp.GetFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, resp.Fail(err.Error()))
			return
		}

		points, err = services.WaveForm.GetLatestWave(wp.sample, wp.count)
		if err != nil {
			c.JSON(http.StatusOK, resp.Error(err.Error()))
			return
		}
		c.JSON(http.StatusOK, resp.Success(points))
		return
	}
}
