// This code implements a Go package that provides methods to interact
// with entities.
package usecase

import (
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
)

// This code defines a use case struct called User that provides
// methods to interact with users in a music service application.
type User struct {
	repo repository.User
}

// Create a new user entity with DI
func NewUser(r repository.User) *User {
	return &User{r}
}

func (u *User) Get() []entity.User {
	return u.repo.Get()
}

func (u *User) Create(user *entity.User) (*entity.User, error) {
	return u.repo.Create(user)
}

func (u *User) Update(id uint, user *entity.User) (*entity.User, error) {
	return u.repo.Update(id, user)
}

func (u *User) Delete(id uint) error {
	return u.repo.Delete(id)
}
