package models

type SaveWave struct {
	ID        uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name      string `gorm:"size:255;not null" json:"name"`
	Timestamp int64  `gorm:"not null" json:"timestamp"`
	Data      []byte `json:"data"`
}
