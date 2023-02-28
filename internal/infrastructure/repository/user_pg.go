package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

type User struct {
	*gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db}
}

func (r *User) Get() []entity.User {
	var users []entity.User
	r.Find(&users)
	return users
}

func (r *User) Create(user *entity.User) (*entity.User, error) {
	if err := r.DB.Migrator().HasTable(user); !err {
		if err := r.DB.Debug().Migrator().CreateTable(user); err != nil {
			return user, err
		}
	}
	if err := r.DB.First(user, "id = ?", user.ID).Error; err != nil {
		if err := r.DB.Create(user).Error; err != nil {
			return user, err
		}
		return user, nil
	} else {
		return user, err
	}
}

func (r *User) Update(id uint, user *entity.User) (*entity.User, error) {
	var finded_user entity.User
	if err := r.DB.First(&finded_user, "id = ?", id).Error; err != nil {
		return user, err
	}
	if err := r.DB.Save(user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *User) Delete(id uint) error {
	var user entity.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return err
	}
	if err := r.DB.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
