package wave

import (
	"encoding/json"
	"errors"
)

// FullData
// Points 和 Head
type FullData struct {
	Header Head `json:"header"`
	Body   Body `json:"body"`
}

type Points struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Body []Points

func (b Body) JSON() ([]byte, error) {
	if len(b) == 0 {
		return nil, errors.New("wave: empty body")
	}

	return json.Marshal(b)
}

type Head struct {
}

func (h Head) String() string {
	// todo Head.String()
	return ""
}

// ParseRawData
// 解析的是原数据 (无头部) 把 data 解析成 points
// 返回的 FullData.Head 为空
func ParseRawData(data []byte, sample, count int) *FullData {
	if sample <= 0 || count <= 0 {
		return nil
	}
	var (
		f   = new(FullData)
		cnt = 0
	)
	for i := -1; i < len(data); i++ {
		if cnt%sample != 0 {
			for i++; i < len(data) && data[i] != '\n'; i++ {
			}
			cnt++
			continue
		}

		num := 0
		for i++; i < len(data) && data[i] != 0x0a; i++ {
			num = num*10 + int(data[i]-'0')
		}
		f.Body = append(f.Body, Points{
			X: len(f.Body),
			Y: num,
		})
		cnt++
		if len(f.Body) == count {
			return f
		}
	}
	return f
}

// Parse 解析包含头部的数据
func Parse(data []byte) *FullData {
	return nil
}
