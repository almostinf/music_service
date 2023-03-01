package entity

import (
	"gorm.io/gorm"
)

type SongNode struct {
	gorm.Model
	SongID     uint  `json:"song_id"`
	NextSongID uint  `json:"next_song_id"`
	PrevSongID uint  `json:"prev_song_id"`
	Song       *Song `json:"-" gorm:"foreignkey:SongID, NextSongID, PrevSongID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Playlist struct {
	gorm.Model
	Name      string   `json:"name" gorm:"type:varchar(64)"`
	Title     string   `json:"title" gorm:"type:text"`
	HeadSong  SongNode `json:"head_song" gorm:"-"`
	TailSong  SongNode `json:"tail_song" gorm:"-"`
	CreatorID uint      `json:"creator_id"`
	User      User     `json:"-" gorm:"foreignkey:CreatorID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
