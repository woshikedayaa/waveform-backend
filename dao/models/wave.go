package models

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type Wave struct {
	// ID read-only
	ID    int       `gorm:"primary_key;auto_increment;->" json:"id"`
	Name  string    `gorm:"size:255;not null" json:"name"`
	CTime time.Time `gorm:"not null;autoCreateTime;column:create_time" json:"createTime"`
	DTime time.Time `gorm:"column:delete_time" json:"deleteTime"`
	Head  string    `gorm:"head" json:"head"`
}

func (w *Wave) Read() ([]byte, error) {
	if w.DTime.Unix() != 0 {
		return nil, errors.New("无法恢复波形图，该波形图已被删除")
	}
	bytes, err := os.ReadFile(w.getSavePath())
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (w *Wave) Write(p []byte) error {
	err := os.WriteFile(w.getSavePath(), p, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (w *Wave) TableName() string {
	return "wave"
}

// 非 windows: /var/lib/waveform/*
// windows: ./data/*
func (w *Wave) getSavePath() string {
	dir := "/var/lib/waveform/"
	if runtime.GOOS == "windows" {
		wd, _ := os.Getwd()
		dir = filepath.Join(wd, "data")
	}
	return filepath.Join(dir, strconv.Itoa(w.ID))
}
