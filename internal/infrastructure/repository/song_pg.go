package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

type Song struct {
	*gorm.DB
}

func NewSong(db *gorm.DB) *Song {
	return &Song{db}
}

func (r *Song) Get() []entity.Song {
	var songs []entity.Song
	r.Find(&songs)
	return songs
}

func (r *Song) Create(song *entity.Song) (*entity.Song, error) {
	if err := r.DB.Migrator().HasTable(song); !err {
		if err := r.DB.Debug().Migrator().CreateTable(song); err != nil {
			return song, err
		}
	}
	if err := r.DB.First(song, "id = ?", song.ID).Error; err != nil {
		if err := r.DB.Create(song).Error; err != nil {
			return song, err
		}
		return song, nil
	} else {
		return song, err
	}
}

func (r *Song) Update(id string, song *entity.Song) (*entity.Song, error) {
	var finded_song entity.Song
	if err := r.DB.First(&finded_song, "id = ?", id).Error; err != nil {
		return song, err
	}
	if err := r.DB.Save(&song).Error; err != nil {
		return song, err
	}
	return song, nil
}

func (r *Song) Delete(id string) error {
	var song entity.Song
	if err := r.DB.First(&song, "id = ?", id).Error; err != nil {
		return err
	}
	if err := r.DB.Delete(&song).Error; err != nil {
		return err
	}
	return nil
}
