package interfaces

import (
	"github.com/almostinf/music_service/internal/entity"
)

type (
	User interface {
		Get() []entity.User
		Create(*entity.User) (*entity.User, error)
		Update(uint, *entity.User) (*entity.User, error)
		Delete(uint) error
	}
)
