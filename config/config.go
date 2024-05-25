package config

import (
	_ "embed"
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

// 真有GPT味道啊(指下面这句哈)（来自某人吐槽）
// 使用 go:embed 注解将 config_full.yaml 嵌入到程序中作为字符串变量configFull
//
//go:embed config_full.yaml
var configFull string

const (
	DefaultConfigFilePrefix = "config"
)

// Config 全部配置
type Config struct {
	Server Servers `json:"server" yaml:"server"`
	DB     DB      `json:"db" yaml:"db"`
	Log    Log     `json:"log" yaml:"log"`
}

type Servers struct {
	Http HttpServer `json:"http" yaml:"http"`
	Kcp  KcpServer  `json:"kcp" yaml:"kcp"`
}

type HttpServer struct {
	Addr string `json:"addr" yaml:"addr"`
	Port int    `json:"port" yaml:"port"`
}

// KcpServer KCP 协议配置（已包括调优配置）
type KcpServer struct {
	Addr     string `json:"addr" yaml:"addr"`
	Port     int    `json:"port" yaml:"port"`
	Mode     string `json:"mode" yaml:"mode"`
	Crypt    string `json:"crypt" yaml:"crypt"`
	Sndwnd   int    `json:"sndwnd" yaml:"sndwnd"`
	Rcvwnd   int    `json:"rcvwnd" yaml:"rcvwnd"`
	Mtu      int    `json:"mtu" yaml:"mtu"`
	NoDelay  int    `json:"noDelay" yaml:"noDelay"`
	Interval int    `json:"interval" yaml:"interval"`
	Resend   int    `json:"resend" yaml:"resend"`
	NC       int    `json:"nc" yaml:"nc"`
}

type DB struct {
	Driver string `json:"driver" yaml:"driver"`
	DbName string `json:"dbName" yaml:"dbName"`
}

// Log 日志配置
type Log struct {
	Output    []string `json:"output" yaml:"output"`
	ErrOutput []string `json:"errOutput" yaml:"errOutput"`
	Level     string   `json:"level" yaml:"level"`
	Format    string   `json:"format" yaml:"format"`
}

// 存储解析后的配置
var config *Config = &Config{}

// G 将配置存入全局，方便其他包引用
func G() *Config {
	return config
}

// InitConfig 配置初始化
func InitConfig() error {
	path, typ, err := findAvailAbleConfigFile()
	if err != nil {
		if os.IsNotExist(err) {
			// 没找到配置文件 从默认的读
			// 返回 os.ErrNotExist 供其他函数将会使用默认配置文件
			viper.SetConfigType("yaml")
			err = viper.ReadConfig(strings.NewReader(configFull))
			if err != nil {
				return err
			}
			err = viper.Unmarshal(config)
			if err != nil {
				return err
			}
			// 默认配置文件不需要后面的配置文件检查了
			return os.ErrNotExist
		} else {
			return err
		}
	} else {
		// 这里是找到可用的配置文件了 就从配置文件读
		// 配置 Viper
		viper.SetConfigType(typ)
		viper.AddConfigPath(GetDefaultConfigFileDir())
		viper.SetConfigFile(path)

		// 读取配置文件
		err = viper.ReadInConfig()
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
}

// findAvailAbleConfigFile 检查支持的配置文件是否存在
// 返回 第一个是全部路径(包含路径和名字) 第二个是返回这个配置文件的类型 第三个返回错误
// 扫描顺序 跟 GetSupportedConfigFileSuffix 的返回的数组顺序有关
func findAvailAbleConfigFile() (path string, Typ string, err error) {
	var (
		dir   = GetDefaultConfigFileDir()
		names = GetSupportedConfigFileName()
	)
	for i := 0; i < len(names); i++ {
		path = filepath.Join(dir, names[i])
		_, err = os.Stat(filepath.Join(dir, names[i]))
		if os.IsNotExist(err) {
			err = nil
			continue
		} else {
			return path, GetSupportedConfigFileSuffix()[i], nil
		}
	}
	return "", "", os.ErrNotExist
}

// GetDefaultConfigFileDir 这个用来获取默认配置文件的路径(文件夹)
func GetDefaultConfigFileDir() string {
	var (
		configPath = ""
	)
	// 这样做是为了符合 FHS 规范（如果是windows系统，获取当前工作目录作为配置路径。否则使用Linux下符合FHS规范的路径）
	if runtime.GOOS == "windows" {
		configPath, _ = os.Getwd()
	} else {
		configPath = "/usr/local/share/etc/waveform/"
	}
	return configPath
}

// GetDefaultConfigFileFull 这个用来获取默认的配置文件全名 包含路径和文件名
func GetDefaultConfigFileFull() string {
	return filepath.Join(GetDefaultConfigFileDir(), GetDefaultConfigFileName())
}

// GetDefaultConfigFileName 这个用来获取默认的配置文件名称
func GetDefaultConfigFileName() string {
	return strings.Join([]string{DefaultConfigFilePrefix, "yaml"}, ".")
}

// GetSupportedConfigFileName 这个可以获得支持的配置文件名称
func GetSupportedConfigFileName() []string {
	var (
		suffixs = GetSupportedConfigFileSuffix()
		res     []string
	)
	for i := 0; i < len(suffixs); i++ {
		res = append(res, strings.Join([]string{DefaultConfigFilePrefix, suffixs[i]}, "."))
	}
	return res
}

// GetSupportedConfigFileSuffix 支持的配置文件格式
func GetSupportedConfigFileSuffix() []string {
	return []string{"json", "yaml", "yml"}
}

// IsConfigFileSupport 这个用来检查一个配置文件是否属于支持的配置文件
func IsConfigFileSupport(name string) bool {
	name = filepath.Base(name)
	sp := strings.Split(name, ".")
	// 这里不复用写好的方法就是这里有优化
	if len(sp) != 2 || sp[0] != DefaultConfigFilePrefix {
		return false
	}
	a := GetSupportedConfigFileSuffix()

	return slices.Contains(a, sp[1])
}

func GetExampleConfig() string {
	return configFull
}
