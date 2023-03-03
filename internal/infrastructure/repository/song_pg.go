// This code implements a Go package that interacts with a database
// and provides CRUD operations for entities. The package uses
// the Gorm ORM to connect to the database.
package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

// The Song type wraps a *gorm.DB object and has four methods:
// Get, Create, Update, Delete which provide basic CRUD functionality
// for the Song entity.
type Song struct {
	*gorm.DB
}

// Create a new instance of song with DI.
func NewSong(db *gorm.DB) *Song {
	return &Song{db}
}

// The Get method retrieves all the songs from the database.
func (r *Song) Get() []entity.Song {
	var songs []entity.Song
	r.Find(&songs)
	return songs
}

// The Create method creates a new song in the database.
func (r *Song) Create(song *entity.Song) (*entity.Song, error) {
	// Check if the song table exist
	if err := r.DB.Migrator().HasTable(song); !err {
		if err := r.DB.Debug().Migrator().CreateTable(song); err != nil {
			return song, err
		}
	}

	// Try to find the song entity in database
	if err := r.DB.First(song, "id = ?", song.ID).Error; err != nil {
		if err := r.DB.Create(song).Error; err != nil {
			return song, err
		}
		return song, nil
	} else {
		return song, err
	}
}

// The Update method updates an existing song.
func (r *Song) Update(id uint, song *entity.Song) (*entity.Song, error) {
	// Update song entity in database
	if err := r.DB.Model(&entity.Song{}).Where("id = ?", id).Updates(song).Error; err != nil {
		return song, err
	}

	return song, nil
}

// The Delete method deletes a song from the database.
func (r *Song) Delete(id uint) error {
	// Try to find a song entity in database
	var song entity.Song
	if err := r.DB.First(&song, "id = ?", id).Error; err != nil {
		return err
	}

	// Delete finded song
	if err := r.DB.Delete(&song).Error; err != nil {
		return err
	}

	return nil
}
