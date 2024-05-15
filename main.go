package main

import (
	"github.com/woshikedayaa/waveform-backend/config"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
}
