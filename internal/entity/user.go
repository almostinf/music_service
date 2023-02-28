package entity

import (
	"gorm.io/gorm"
)

type CurSongInfo struct {
	CurPlayingSongID uint      `json:"cur_playing_song_id"`
	Song             *Song     `gorm:"foreignkey:CurPlayingSongID"`
	IsPlaying        bool      `json:"is_playing"`
	PlaylistID       uint      `json:"playlist_id"`
	Playlist         *Playlist `gorm:"foreignkey:PlaylistID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type User struct {
	gorm.Model
	FirstName string       `json:"first_name" gorm:"type:varchar(64)"`
	LastName  string       `json:"last_name" gorm:"type:varchar(64)"`
	Email     string       `json:"email" gorm:"unique;not null"`
	Password  string       `json:"password" gorm:"not null"`
	CurSong   *CurSongInfo `json:"cur_song" gorm:"foreignkey:CurPlayingSongID"`
	Playlists []*Playlist  `json:"playlists" gorm:"many2many:user_playlists;"`
}
