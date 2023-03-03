// This code implements a Go package that provides methods to interact
// with entities.
package usecase

import (
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
)

// This code defines a use case struct called Song that provides
// methods to interact with songs in a music service application.
type Song struct {
	repo repository.Song
}

// Create a new song entity with DI
func NewSong(r repository.Song) *Song {
	return &Song{r}
}

func (u *Song) Get() []entity.Song {
	return u.repo.Get()
}

func (u *Song) Create(song *entity.Song) (*entity.Song, error) {
	return u.repo.Create(song)
}

func (u *Song) Update(id uint, song *entity.Song) (*entity.Song, error) {
	return u.repo.Update(id, song)
}

func (u *Song) Delete(id uint) error {
	return u.repo.Delete(id)
}
