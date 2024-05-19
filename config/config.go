package config

import (
	_ "embed"
	"errors"
	"github.com/spf13/viper"
	"github.com/woshikedayaa/waveform-backend/api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// 使用 //go:embed 注解将 config_full.yaml 嵌入到程序中作为字符串变量configFull
//
//go:embed config_full.yaml
var configFull string

// DB 数据库实例
var db *gorm.DB

// 全部配置
type Config struct {
	Server Server `json:"server" yaml:"server"`
	Mysql  Mysql  `json:"mysql" yaml:"mysql"`
	Log    Log    `json:"log" yaml:"log"`
	Kcp    Kcp    `json:"kcp" yaml:"kcp"`
}

// 服务器监听的地址和端口
type Server struct {
	Port int    `json:"port" yaml:"port"`
	Addr string `json:"addr" yaml:"addr"`
}

// 数据库配置
type Mysql struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	//高级配置，例如 charset
	Config   string `yaml:"config"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	//日志等级：debug（输出全部sql），dev（开发环境，只输出error），release（生产环境）
	LogLevel string `yaml:"log_level"`
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

// Dsn是定义在Mysql结构体上的方法，构建并返回一个符合MySQL连接字符串。

func (m *Mysql) Dsn() string {
	dsn := m.User + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.DB + "?" + m.Config
	return dsn
}

// InitGorm 初始化GORM数据库连接
func InitGorm() error {
	var err error
	dsn := G().Mysql.Dsn()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置数据库连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大存活时间
	sqlDB.SetConnMaxLifetime(time.Hour * 4)

	// 自动迁移数据库模型
	if err := db.AutoMigrate(&models.SaveWave{}); err != nil {
		return err
	}

	return nil
}

// DB 获取数据库实例
func DB() *gorm.DB {
	return db
}
