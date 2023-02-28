package usecase

import (
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
)

type User struct {
	repo repository.User
}

func NewUser(r repository.User) *User {
	return &User{r}
}

func (u *User) Get() []entity.User {
	return u.repo.Get()
}

func (u *User) Create(user *entity.User) (*entity.User, error) {
	return u.repo.Create(user)
}

func (u *User) Update(id string, user *entity.User) (*entity.User, error) {
	return u.repo.Update(id, user)
}

func (u *User) Delete(id string) error {
	return u.repo.Delete(id)
}
