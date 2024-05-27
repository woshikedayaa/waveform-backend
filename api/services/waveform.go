package services

import (
	"errors"
)

type waveForm struct{}

var WaveForm waveForm

type Point struct {
	X int   `json:"x"`
	Y uint8 `json:"y"`
}

func (waveForm) ConvDataToPoints(data []byte, sample, count int) ([]Point, error) {
	if sample == 0 || count == 0 {
		return nil, errors.New("service: GetLatestWave sample and count must >= 1")
	}
	var (
		cnt = 0
		res []Point
	)
	for i := -1; i < len(data); i++ {
		num := uint8(0)
		for i++; i < len(data) && data[i] != 0x0a; i++ {
			num = num*10 + (data[i] - '0')
		}
		if cnt%sample == 0 {
			res = append(res, Point{
				X: len(res),
				Y: num,
			})
		}
		cnt++
		if len(res) == count {
			return res, nil
		}
	}
	return res, nil
}
