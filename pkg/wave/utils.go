package wave

import (
	"go.bug.st/serial"
	"math/rand/v2"
	"slices"
)

func ScanAvailablePort() ([]string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		return nil, &OpError{
			err: err,
			op:  "scan",
		}
	}

	return ports, nil
}

func GenerateRandomData(count int) []byte {
	if count <= 0 {
		return nil
	}

	Uint8ToAscii := func(i uint8) []byte {
		var res []byte
		for i != 0 {
			cur := i % 10
			res = append(res, cur+'0')
			i /= 10
		}
		slices.Reverse(res)
		return res
	}
	var res []byte

	for {
		n := rand.Uint32()
		// 这里复用一下 这里生成的是 uint32 可以分为四个 uint8
		// 可以减少随机数生成的次数
		for i := 0; i < 4 && len(res) < count; i++ {
			cur := uint8(n)
			n >>= 8

			res = append(res, Uint8ToAscii(cur)...)
			res = append(res, '\n')
		}
		if len(res) == count {
			return res
		}
	}
}
