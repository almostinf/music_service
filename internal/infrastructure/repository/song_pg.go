package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

type Song struct {
	*gorm.DB
}

func NewSongRepository(db *gorm.DB) *Song {
	return &Song{db}
}

func (r *Song) Get() []entity.Song {
	var songs []entity.Song
	r.Find(&songs)
	return songs
}

func (r *Song) Create(song *entity.Song) (*entity.Song, error) {
	if err := r.DB.Create(song).Error; err != nil {
		return song, err
	}
	return song, nil
}
