package interfaces

import (
	"github.com/almostinf/music_service/internal/entity"
)

type (
	User interface {
		Get() []entity.User
		Add(*entity.User) (*entity.User, error)
	}
)
