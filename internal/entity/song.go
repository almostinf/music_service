// Package entity defines the data structures used in the application.
package entity

import (
	"time"

	"gorm.io/gorm"
)

// The Song struct defines the properties of a music song entity in the system.
type Song struct {
	gorm.Model
	Title     string        `json:"title" gorm:"type:text"`
	Artist    string        `json:"artist" gorm:"type:varchar(64)"`
	Duration  time.Duration `json:"duration" gorm:"not null"`
	Playlists []Playlist    `gorm:"many2many:playlist_songs;"`
}
