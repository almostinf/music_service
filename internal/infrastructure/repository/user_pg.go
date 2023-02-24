package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Get() []entity.User {
	var users []entity.User
	r.Find(&users)
	return users
}

func (r *UserRepository) Add(user *entity.User) (*entity.User, error) {
	if err := r.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
