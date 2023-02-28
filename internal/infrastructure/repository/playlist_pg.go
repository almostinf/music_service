package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

type Playlist struct {
	*gorm.DB
}

func NewPlaylist(db *gorm.DB) *Playlist {
	return &Playlist{db}
}

func (r *Playlist) Get() []entity.Playlist {
	var playlists []entity.Playlist
	r.Find(&playlists)
	return playlists
}

func (r *Playlist) Create(playlist *entity.Playlist) (*entity.Playlist, error) {
	if err := r.DB.Migrator().HasTable(playlist); !err {
		if err := r.DB.Debug().Migrator().CreateTable(playlist); err != nil {
			return playlist, err
		}
	}
	if err := r.DB.First(playlist, "id = ?", playlist.ID).Error; err != nil {
		if err := r.DB.Create(playlist).Error; err != nil {
			return playlist, err
		}
		return playlist, nil
	} else {
		return playlist, err
	}
}

func (r *Playlist) Update(id uint, playlist *entity.Playlist) (*entity.Playlist, error) {
	var finded_playlist entity.Playlist
	if err := r.DB.First(&finded_playlist, "id = ?", id).Error; err != nil {
		return playlist, err
	}
	if err := r.DB.Save(playlist).Error; err != nil {
		return playlist, err
	}
	return playlist, nil
}

func (r *Playlist) Delete(id uint) error {
	var playlist entity.Playlist
	if err := r.DB.First(&playlist, "id = ?", id).Error; err != nil {
		return err
	}
	if err := r.DB.Delete(&playlist).Error; err != nil {
		return err
	}
	return nil
}
