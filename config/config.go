package config

import (
	_ "embed"
	"errors"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

//go:embed config_full.yaml
var configFull string

type Config struct {
	Server Server `json:"server" yaml:"server"`
	Log    Log    `json:"log" yaml:"log"`
	Device Device `json:"device" yaml:"device"`
}

type Server struct {
	Port int    `json:"port" yaml:"port"`
	Addr string `json:"addr" yaml:"addr"`
}

type Log struct {
	Output []string `json:"output" yaml:"output"`
	Level  string   `json:"level" yaml:"level"`
	Format string   `json:"format" yaml:"format"`
}

type Device struct {
	PortName string `json:"portName" yaml:"portName"`
}

var config *Config

func Get() *Config {
	return config
}

func InitConfig() error {
	var (
		configPath = ""
		configName = "config.yaml"
		configType = "yaml"
	)

	// 这样做是为了符合 FHS 规范
	if runtime.GOOS == "windows" {
		configPath, _ = os.Getwd()
	} else {
		configPath = "/usr/local/share/etc/waveform"
	}

	//
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	viper.SetConfigFile(configName)

	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("config: " + err.Error())
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return errors.New("config: " + err.Error())
	}
	return nil
}
