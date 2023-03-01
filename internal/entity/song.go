package entity

import (
	"time"

	"gorm.io/gorm"
)

type Song struct {
	gorm.Model
	Title     string        `json:"title" gorm:"type:text"`
	Artist    string        `json:"artist" gorm:"type:varchar(64)"`
	Duration  time.Duration `json:"duration" gorm:"not null"`
	Playlists []Playlist    `gorm:"many2many:playlist_songs;"`
}
