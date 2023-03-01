package usecase

import (
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
)

type Playlist struct {
	repo repository.Playlist
}

func NewPlaylist(r repository.Playlist) *Playlist {
	return &Playlist{r}
}

func (u *Playlist) Get() []entity.Playlist {
	return u.repo.Get()
}

func (u *Playlist) Create(playlist *entity.Playlist) (*entity.Playlist, error) {
	return u.repo.Create(playlist)
}

func (u *Playlist) Update(id uint, playlist *entity.Playlist) (*entity.Playlist, error) {
	return u.repo.Update(id, playlist)
}

func (u *Playlist) Delete(id uint) error {
	return u.repo.Delete(id)
}

func (u *Playlist) Play(user_id uint, song_id uint, playlist_id uint) error {
	return u.repo.Play(user_id, song_id, playlist_id)
}

func (u *Playlist) Pause(user_id uint, song_id uint, playlist_id uint) error {
	return u.repo.Pause(user_id, song_id, playlist_id)
}

func (u *Playlist) AddSong(user_id uint, song_id uint, playlist_id uint) error {
	return u.repo.AddSong(user_id, song_id, playlist_id)
}

func (u *Playlist) Next(user_id uint, song_id uint, playlist_id uint) error {
	return u.repo.Next(user_id, song_id, playlist_id)
}

func (u *Playlist) Prev(user_id uint, song_id uint, playlist_id uint) error {
	return u.repo.Prev(user_id, song_id, playlist_id)
}
