package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

type User struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) *User {
	return &User{db}
}

func (r *User) Get() []entity.User {
	var users []entity.User
	r.Find(&users)
	return users
}

func (r *User) Create(user *entity.User) (*entity.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return user, err
	}
	return user, nil
}
