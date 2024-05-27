package config

import (
	"fmt"
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	err := InitConfig()
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	fmt.Printf("%#v", G())
}

func TestCheckConfig(t *testing.T) {
	err := getExampleCorrectConfig().check()
	if err != nil {
		t.Fatal(err)
	}
	err = getExampleMistakeConfig().check()
	if err == nil {
		t.Fatal()
	} else {
		fmt.Println(err.Error())
	}
}

func getExampleCorrectConfig() *Config {
	return &Config{
		Server: &Servers{
			Http: &HttpServer{
				Addr: "0.0.0.0",
				Port: 8080,
			},
			Kcp: &KcpServer{
				Addr: "0.0.0.0",
				Port: 8080,
			},
		},
		DB: &DB{
			Driver: "sqlite",
			DbName: "waveform",
		},
		Log: &Log{
			Output:    nil,
			ErrOutput: nil,
			Level:     "info",
			Format:    "json",
		},
	}
}

func getExampleMistakeConfig() *Config {
	return &Config{
		Server: &Servers{},
		DB:     &DB{},
		Log: &Log{
			Level: "unknown",
		},
	}
}
