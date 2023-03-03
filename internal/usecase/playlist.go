// This code implements a Go package that provides methods to interact 
// with entities.
package usecase

import (
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
)

// This code defines a use case struct called Playlist that provides 
// methods to interact with playlists in a music service application.
type Playlist struct {
	repo repository.Playlist
}

// Create a new playlist entity with DI
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
