package entity

import (
	"time"

	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	ID        string        `json:"id" gorm:"primaryKey"`
	Title     string        `json:"title" gorm:"type:text"`
	Artist    string        `json:"artist" gorm:"type:varchar(64)"`
	Duration  time.Duration `json:"duration" gorm:"not null"`
	CreatedAt time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
