package dao

import (
	"github.com/woshikedayaa/waveform-backend/config"
	"testing"
)

func TestConnect(t *testing.T) {
	config.G().DB.Driver = "sqlite"
	config.G().DB.DbName = "../waveform"

	err := InitDataBase()
	if err != nil {
		t.Fatal(err)
	}
}
