package interfaces

import (
	"github.com/almostinf/music_service/internal/entity"
)

type (
	Song interface {
		Get() []entity.Song
		Add(entity.Song) (entity.Song, error)
	}

	SongRepository interface {
		Song
	}
)
