package config

import (
	_ "embed"
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	Output    []string `json:"output" yaml:"output"`
	ErrOutput []string `json:"errOutput" yaml:"errOutput"`
	Level     string   `json:"level" yaml:"level"`
	Format    string   `json:"format" yaml:"format"`
}

type Device struct {
	PortName string `json:"portName" yaml:"portName"`
}

var config *Config = &Config{}

func G() *Config {
	return config
}

func InitConfig() error {
	var (
		configPath = filepath.Dir(GetDefaultConfigFilePath())
		configName = "config.yaml"
		configType = "yaml"
	)

	//
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	viper.SetConfigFile(filepath.Join(configPath, configName))
	// 看有没有配置文件 如果有配置文件就从配置文件读 没有就从默认配置文件读
	_, err := os.Stat(filepath.Join(configPath, configName))
	if os.IsNotExist(err) {
		err = viper.ReadConfig(strings.NewReader(configFull))
	} else {
		err = viper.ReadInConfig()
	}
	if err != nil {
		return errors.New("config: " + err.Error())
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return errors.New("config: " + err.Error())
	}
	// todo config检查

	return nil
}

func GetDefaultConfigFilePath() string {
	var configPath string
	// 这样做是为了符合 FHS 规范
	if runtime.GOOS == "windows" {
		configPath, _ = os.Getwd()
	} else {
		configPath = "/usr/local/share/etc/waveform/"
	}
	return filepath.Join(configPath, "config.yaml")
}
