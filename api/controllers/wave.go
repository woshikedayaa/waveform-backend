package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
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
		return errors.New("count 太大或者不是一个数字")
	}
	w.sample, err = strconv.Atoi(samples)
	if err != nil {
		return errors.New("sample 太大或者不是一个数字")
	}
	return nil
}
