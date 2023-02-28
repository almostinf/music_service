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
	}
)
