// Package entity defines the data structures used in the application.
package entity

import (
	"gorm.io/gorm"
)

// This struct defines a song node entity that represents a single node 
// in a double-linked playlist.
type SongNode struct {
	gorm.Model
	SongID     uint  `json:"song_id"`
	NextSongID uint  `json:"next_song_id"`
	PrevSongID uint  `json:"prev_song_id"`
	Song       *Song `json:"-" gorm:"foreignkey:SongID"`
}

// This struct represents a double-linked playlist entity with additional 
// information about name, title and creator id.
type Playlist struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(64)"`
	Title     string `json:"title" gorm:"type:text"`
	HeadSong  uint   `json:"head_song"`
	TailSong  uint   `json:"tail_song"`
	CreatorID uint   `json:"creator_id"`
}
