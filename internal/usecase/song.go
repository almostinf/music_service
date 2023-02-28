package usecase

import (
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
)

type Song struct {
	repo repository.Song
}

func NewSong(r repository.Song) *Song {
	return &Song{r}
}

func (u *Song) Get() []entity.Song {
	return u.repo.Get()
}

func (u *Song) Create(song *entity.Song) (*entity.Song, error) {
	return u.repo.Create(song)
}

func (u *Song) Update(id string, song *entity.Song) (*entity.Song, error) {
	return u.repo.Update(id, song)
}

func (u *Song) Delete(id string) error {
	return u.repo.Delete(id)
}
