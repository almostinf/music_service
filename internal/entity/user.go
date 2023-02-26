package entity

import (
	"container/list"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string     `json:"id" gorm:"primaryKey; type:uuid;default:gen_random_uuid()"`
	FirstName string     `json:"first_name" gorm:"type:varchar(64)"`
	LastName  string     `json:"last_name" gorm:"type:varchar(64)"`
	Email     string     `json:"email" gorm:"unique;not null"`
	Password  string     `json:"password" gorm:"not null"`
	Playlist  Playlist   `gorm:"foreignkey:Playlists"`
	Playlists *list.List `json:"playlists"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
