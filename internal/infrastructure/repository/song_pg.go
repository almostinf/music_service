package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

type SongRepository struct {
	*gorm.DB
}

func NewSongRep(db *gorm.DB) *SongRepository {
	return &SongRepository{db}
}

func (r *SongRepository) Get() []entity.Song {
	var songs []entity.Song
	r.Find(&songs)
	return songs
}

func (r *SongRepository) Add()

func GetSongs(db *gorm.DB) []entity.Song {
	var songs []entity.Song
	db.Find(&songs)
	return songs
}

func AddSong(db *gorm.DB, song *entity.Song) (*entity.Song, error) {
	if err := db.Save(&song).Error; err != nil {
		return song, err
	}
	return song, nil
}
