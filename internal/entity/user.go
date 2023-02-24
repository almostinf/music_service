package entity

import (
	"container/list"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	ID        uint64    `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name" gorm:"type:varchar(64)"`
	LastName  string    `json:"last_name" gorm:"type:varchar(64)"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Playlists list.List `json:"playlists" gorm:"foreignkey:CreatorID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
