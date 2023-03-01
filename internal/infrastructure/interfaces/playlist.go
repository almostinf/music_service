package interfaces

import (
	"github.com/almostinf/music_service/internal/entity"
)

type (
	Playlist interface {
		Get() []entity.Playlist
		Create(*entity.Playlist) (*entity.Playlist, error)
		Update(uint, *entity.Playlist) (*entity.Playlist, error)
		Delete(uint) error
		Play(uint, uint, uint) error
		Pause(uint, uint, uint) error
		AddSong(uint, uint, uint) error
		Next(uint, uint, uint) error
		Prev(uint, uint, uint) error
	}
)
