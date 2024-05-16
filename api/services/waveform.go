package services

import (
	"errors"
	"os"
)

var latest []byte

// GetLatestWave
/*
这里只是模拟一下获取数据 后面等板子出来就换成真正的实现
如果你要获取这个数据 run script/square_wave
*/
// sample	取样率 多少个数据取一个数据画图
// count	最多取多少个数据
func GetLatestWave(sample, count int) ([]Point, error) {
	if sample == 0 || count == 0 {
		return nil, errors.New("service: GetLatestWave sample and count must >= 1")
	}
	if len(latest) == 0 {
		bytes, err := os.ReadFile("test/square_wave.txt")
		if err != nil {
			return nil, err
		}
		latest = bytes
	}
	var (
		cnt = 0
		res []Point
	)
	for i := -1; i < len(latest); i++ {
		num := uint8(0)
		for i++; i < len(latest) && latest[i] != 0x0a; i++ {
			num = num*10 + (latest[i] - '0')
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

type Point struct {
	X int   `json:"x"`
	Y uint8 `json:"y"`
}
