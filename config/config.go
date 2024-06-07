package config

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/woshikedayaa/waveform-backend/pkg/utils"
	"go.uber.org/zap/zapcore"
	"math"
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
	Server *Servers `json:"server" yaml:"server"`
	DB     *DB      `json:"db" yaml:"db"`
	Log    *Log     `json:"log" yaml:"log"`
}

type Servers struct {
	Http *HttpServer `json:"http" yaml:"http"`
	Kcp  *KcpServer  `json:"kcp" yaml:"kcp"`
}

type HttpServer struct {
	Addr string `json:"addr" yaml:"addr"`
	Port int    `json:"port" yaml:"port"`
	Cors *Cors  `json:"cors" yaml:"cors"`
}

type Cors struct {
	Enabled bool     `json:"enabled" yaml:"enabled"`
	Origins []string `json:"origins" yaml:"origins"`
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

// check 用来检查 config 是否合法
func (c *Config) check() error {
	assertZero := func(name string, val int) error {
		if val == 0 {
			return errors.New(fmt.Sprintf("config: %s 为 空或者值不能为 0 ", name))
		}
		return nil
	}

	assertNil := func(name string, val any) error {
		e := errors.New(fmt.Sprintf("config: %s 是 nil", name))
		if utils.CheckIsNil(val) {
			return e
		}
		return nil
	}

	assertRange := func(name string, val, min, max int) error {
		if val > max || val < min {
			return errors.New(fmt.Sprintf("config: %s=%d 不在范围 %d-%d 内", name, val, min, max))
		}
		return nil
	}
	var err error
	join := func(e error) {
		if e == nil {
			return
		}
		err = errors.Join(err, e)
	}
	// server
	if c.Server != nil {
		if c.Server.Http != nil {
			join(assertRange("config.server.http.port", c.Server.Http.Port, 1, math.MaxUint16))
			join(assertZero("config.server.http.addr", len(c.Server.Http.Addr)))
		} else {
			join(assertNil("config.server.http", c.Server.Http))
		}
		// todo kcp 的完整校验
		if c.Server.Kcp != nil {
			join(assertRange("config.server.kcp.port", c.Server.Kcp.Port, 1, math.MaxUint16))
			join(assertZero("config.server.kcp.addr", len(c.Server.Kcp.Addr)))
		} else {
			join(assertNil("config.server.kcp", c.Server.Kcp))
		}
	}
	// log
	if c.Log != nil {
		join(assertZero("config.log.format", len(c.Log.Format)))
		join(assertZero("config.log.level", len(c.Log.Level)))

		if !slices.Contains([]string{
			"json", "console", "console_with_color",
		}, c.Log.Format) {
			join(errors.New(fmt.Sprintf("config: %s=%s 不支持的输出格式", "config.log.format", c.Log.Format)))
		}
		if _, e2 := zapcore.ParseLevel(c.Log.Level); e2 != nil {
			join(errors.New(fmt.Sprintf("config: %s=%s 不能作为 日志等级", "config.log.level", c.Log.Level)))
		}
	}

	// db
	if c.DB != nil {
		join(assertZero("config.db.driver", len(c.DB.Driver)))
		join(assertZero("config.db.dbName", len(c.DB.Driver)))
	}
	// 因为这里不能直接调用dao 就没法判断是否是支持的 交给dao去判断了

	// 返回错误
	return err
}

// InitConfig 配置初始化
func InitConfig() error {
	defaultVipe := viper.New()
	userViper := viper.New()
	// 先读取默认的配置文件
	defaultVipe.SetConfigType("yaml")
	err := defaultVipe.ReadConfig(strings.NewReader(GetExampleConfig()))
	if err != nil {
		return errors.New("config: 读取默认配置错误 err: " + err.Error())
	}
	// 这里读取一遍配置文件
	path, typ, err := findAvailAbleConfigFile()
	if err == nil {
		// 这里是找到可用的配置文件了 就再配置文件读
		// 配置 Viper
		userViper.SetConfigType(typ)
		userViper.AddConfigPath(GetDefaultConfigFileDir())
		userViper.SetConfigFile(path)

		// 读取配置文件
		err = userViper.ReadInConfig()
		if err != nil {
			return errors.New("config: " + err.Error())
		}
	}
	// 合并两个
	// 先判断用户是否设置了数组 如果设置了数组就用用户的
	// viper的默认合并数组的策略是 合并 不是 直接修改
	if userViper.IsSet("log.output") {
		defaultVipe.Set("log.output", nil)
	}

	if userViper.IsSet("log.errOutput") {
		defaultVipe.Set("log.errOutput", nil)
	}
	if userViper.IsSet("server.http.cors.origins") {
		defaultVipe.Set("server.http.cors.origins", nil)
	}

	// 合并
	_ = viper.MergeConfigMap(defaultVipe.AllSettings())
	_ = viper.MergeConfigMap(userViper.AllSettings())
	// 将配置内容反序列化到config变量中
	err = viper.Unmarshal(config)
	if err != nil {
		return errors.New("config: " + err.Error())
	}
	// 检查 config
	return config.check()
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
		configPath = "/usr/local/etc/waveform/"
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
