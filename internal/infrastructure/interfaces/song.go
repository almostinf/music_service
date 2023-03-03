// This package is needed to declare interfaces
package interfaces

import (
	"github.com/almostinf/music_service/internal/entity"
)

// This code defines an interface named Song that specifies a set 
// of methods that should be implemented by types that represent a song.
type (
	Song interface {
		Get() []entity.Song
		Create(*entity.Song) (*entity.Song, error)
		Update(uint, *entity.Song) (*entity.Song, error)
		Delete(uint) error
	}
)
