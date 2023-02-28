package interfaces

import (
	"github.com/almostinf/music_service/internal/entity"
)

type (
	Song interface {
		Get() []entity.Song
		Create(*entity.Song) (*entity.Song, error)
		Update(string, *entity.Song) (*entity.Song, error)
		Delete(string) error
	}
)
