package dao

import (
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/woshikedayaa/waveform-backend/config"
	"github.com/woshikedayaa/waveform-backend/dao/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var db *gorm.DB

type DriverType uint8

const (
	UnknownDriver DriverType = math.MaxUint8
	SqliteDriver  DriverType = iota
	MysqlDriver
	H2Driver
	PostGreDriver
)

func InitDataBase() error {
	var (
		err error
	)

	driver := ParseDriver(config.G().DB.Driver)
	if driver == UnknownDriver {
		return &OpErr{
			op:         "parse driver",
			err:        errors.New(fmt.Sprintf("未知的驱动类型 %s", config.G().DB.Driver)),
			suggestion: "阅读文档",
		}
	}
	db, err = openConnection(driver, config.G().DB)
	if err != nil {
		ope := &OpErr{}
		if errors.As(err, &ope) {
			return err
		}
		ope.err = err
		ope.op = "connect"
		return ope
	}
	// 建表
	err = db.AutoMigrate(&models.Wave{})
	if err != nil {
		ope := &OpErr{}

		ope.err = err
		ope.op = "create"
		ope.suggestion = "反馈给开发者"
		return ope
	}
	// finish
	return nil
}

func ParseDriver(s string) DriverType {
	driver := strings.ToLower(s)
	m := map[string]DriverType{
		"sqlite":  SqliteDriver,
		"postgre": PostGreDriver,
		"mysql":   MysqlDriver,
		"h2":      H2Driver,
	}

	if dt, ok := m[driver]; ok {
		return dt
	}
	return UnknownDriver
}

func openConnection(dt DriverType, dc *config.DB) (*gorm.DB, error) {
	dsn := ""
	switch dt {
	case SqliteDriver:
		dsn = strings.Join([]string{filepath.Join(getDataBaseDir(), dc.DbName), "db"}, ".")
		return gorm.Open(sqlite.Open(dsn))
	case MysqlDriver:
		// todo mysql dsn
		return gorm.Open(mysql.Open(dsn))
	// todo full support current table type
	default:
		return nil, &OpErr{
			op:  "connect",
			err: errors.New("不支持的数据库类型"),
		}
	}
}

// Conn 返回当前数据库连接
func Conn() *gorm.DB {
	return db
}

func getDataBaseDir() string {
	var path string
	if runtime.GOOS != "windows" {
		path = "/var/lib/waveform/"
	} else {
		dir, _ := os.Getwd()
		path = filepath.Join(dir, "data")
	}
	return path
}
