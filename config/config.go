package config

import (
	_ "embed"
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

// 使用 //go:embed 注解将 config_full.yaml 嵌入到程序中作为字符串变量configFull
//
//go:embed config_full.yaml
var configFull string

// 全部配置
type Config struct {
	Server Server `json:"server" yaml:"server"`
	Log    Log    `json:"log" yaml:"log"`
	Kcp    Kcp    `json:"kcp" yaml:"kcp"`
}

// 服务器监听的地址和端口
type Server struct {
	Port int    `json:"port" yaml:"port"`
	Addr string `json:"addr" yaml:"addr"`
}

// 日志配置
type Log struct {
	Output    []string `json:"output" yaml:"output"`
	ErrOutput []string `json:"errOutput" yaml:"errOutput"`
	Level     string   `json:"level" yaml:"level"`
	Format    string   `json:"format" yaml:"format"`
}

// KCP 协议配置（已包括调优配置）
type Kcp struct {
	ListenAddress string `json:"listen_address" yaml:"listen_address"`
	Mode          string `json:"mode" yaml:"mode"`
	Crypt         string `json:"crypt" yaml:"crypt"`
	Sndwnd        int    `json:"sndwnd" yaml:"sndwnd"`
	Rcvwnd        int    `json:"rcvwnd" yaml:"rcvwnd"`
	Mtu           int    `json:"mtu" yaml:"mtu"`
	Nodelay       int    `json:"nodelay" yaml:"nodelay"`
	Interval      int    `json:"interval" yaml:"interval"`
	Resend        int    `json:"resend" yaml:"resend"`
	Nc            int    `json:"nc" yaml:"nc"`
}

// 存储解析后的配置
var config *Config = &Config{}

// 将配置存入全局，方便其他包引用
func G() *Config {
	return config
}

// 配置初始化
func InitConfig() error {
	var (
		configPath = ""
		configName = "config.yaml"
		configType = "yaml"
	)

	// 这样做是为了符合 FHS 规范（如果是windows系统，获取当前工作目录作为配置路径。否则使用Linux下符合FHS规范的路径）
	if runtime.GOOS == "windows" {
		configPath, _ = os.Getwd()
	} else {
		configPath = "/usr/local/share/etc/waveform/"
	}

	// 配置 Viper
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	viper.SetConfigFile(filepath.Join(configPath, configName))

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("config: " + err.Error())
	}
	// 将配置文件内容反序列化到config变量中
	err = viper.Unmarshal(config)
	if err != nil {
		return errors.New("config: " + err.Error())
	}
	// todo config检查

	return nil
}
