// Package entity defines the data structures used in the application.
package entity

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// The CurSongInfo entity stores the information about the current
// playing song for a user, such as the ID of the song, whether it is
// playing or not, and the ID of the playlist it belongs to.
type CurSongInfo struct {
	gorm.Model
	CurPlayingSongID uint `json:"cur_playing_song_id"`
	IsPlaying        bool `json:"is_playing"`
	PlaylistID       uint `json:"playlist_id"`
}

// The User entity stores user information such as first name,
// last name, email, password, current song and playlists.
type User struct {
	gorm.Model
	FirstName     string        `json:"first_name" gorm:"type:varchar(64)"`
	LastName      string        `json:"last_name" gorm:"type:varchar(64)"`
	Email         string        `json:"email" gorm:"unique;not null"`
	Password      string        `json:"password" gorm:"not null"`
	CurSongInfoID uint          `json:"cur_song_info_id"`
	CurSongInfo   CurSongInfo   `json:"-"`
	Playlists     pq.Int64Array `json:"playlists" gorm:"type:integer[]; foreignkey:CreatorID"`
}
