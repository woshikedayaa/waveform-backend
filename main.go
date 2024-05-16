package main

import (
	"github.com/woshikedayaa/waveform-backend/api"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/logf"
	"strconv"
	"strings"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = logf.LoggerInit()
	if err != nil {
		panic(err)
	}
	router := api.InitRouter()

	// start!
	if err = router.Run(strings.Join(
		[]string{
			config.G().Server.Addr,
			strconv.Itoa(config.G().Server.Port),
		}, ":")); err != nil {
		panic(err)
	}
}
