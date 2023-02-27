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
	if err := r.DB.First(&user, "id = ?", user.ID).Error; err != nil {
		if err := r.DB.Create(user).Error; err != nil {
			return user, err
		}
		return user, nil
	} else {
		return user, err
	}
}

func (r *User) Update(user *entity.User) (*entity.User, error) {
	var finded_user entity.User
	if err := r.DB.First(&finded_user, "id = ?", user.ID).Error; err != nil {
		return user, err
	}
	if err := r.DB.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *User) Delete(id string) error {
	var user entity.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return err
	}
	if err := r.DB.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
