package dao

import (
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/woshikedayaa/waveform-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
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

	driver := parseDriver(config.G().DB.Driver)
	if driver == UnknownDriver {
		return &OpErr{
			op:         "parse driver",
			err:        errors.New(fmt.Sprintf("unknown driver %s", config.G().DB.Driver)),
			suggestion: "see the document",
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
	// finish
	return nil
}

func parseDriver(s string) DriverType {
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

func openConnection(dt DriverType, dc config.DB) (*gorm.DB, error) {
	dsn := ""
	switch dt {
	case SqliteDriver:
		dsn = strings.Join([]string{dc.DbName, "db"}, ".")
		return gorm.Open(sqlite.Open(dsn))
	case MysqlDriver:
		// todo build mysql dsn
		return gorm.Open(mysql.Open(dsn))
	// todo full support current database type
	default:
		return nil, &OpErr{
			op:         "connect",
			suggestion: "unsupported database type",
		}
	}
}
