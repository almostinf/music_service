package http

import (
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
)

type User struct {
	repo repository.User
}

func NewUserUseCase(r repository.User) *User {
	return &User{r}
}

func (u *User) Get() []entity.User {
	return u.repo.Get()
}

func (u *User) Add(user *entity.User) (*entity.User, error) {
	user, err := u.repo.Add(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
