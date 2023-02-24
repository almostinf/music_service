package entity

import (
	"container/list"
	"time"

	"gorm.io/gorm"
)

type Playlist struct {
	gorm.Model
	ID        uint64        `json:"id" gorm:"primaryKey"`
	Title     string        `json:"title" gorm:"type:text"`
	User      User          `gorm:"foreignkey:CreatorID"`
	CreatorID uint64        `json:"creator_id" gorm:"not null"`
	Songs     *list.List    `json:"songs" gorm:"many2many:playlist_songs;"`
	Song      Song          `gorm:"foreignkey:CurSongID"`
	CurSongID *list.Element `json:"current_song"`
	CreatedAt time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
