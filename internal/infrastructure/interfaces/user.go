// This package is needed to declare interfaces
package interfaces

import (
	"github.com/almostinf/music_service/internal/entity"
)

// This code defines an interface named User that specifies a set
// of methods that should be implemented by types that represent a user.
type (
	User interface {
		Get() []entity.User
		Create(*entity.User) (*entity.User, error)
		Update(uint, *entity.User) (*entity.User, error)
		Delete(uint) error
	}
)
