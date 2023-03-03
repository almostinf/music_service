// This code implements a Go package that interacts with a database
// and provides CRUD operations for entities. The package uses
// the Gorm ORM to connect to the database.
package repository

import (
	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

// The User type wraps a *gorm.DB object and has four methods:
// Get, Create, Update, Delete which provide basic CRUD functionality
// for the User entity.
type User struct {
	*gorm.DB
}

// Create a new instance of user entity with DI.
func NewUser(db *gorm.DB) *User {
	return &User{db}
}

// The Get method retrieves all the users from the database.
func (r *User) Get() []entity.User {
	var users []entity.User
	r.Find(&users)
	return users
}

// The Create method creates a new user in the database.
func (r *User) Create(user *entity.User) (*entity.User, error) {
	// Check if the current song info table exist
	var curSong entity.CurSongInfo
	if err := r.DB.Migrator().HasTable(&curSong); !err {
		if err := r.DB.Debug().Migrator().CreateTable(&curSong); err != nil {
			return user, err
		}
	}

	// Create a new entity of CurSongIngo
	if err := r.DB.Model(&entity.CurSongInfo{}).Create(&curSong).Error; err != nil {
		return user, err
	}

	user.CurSongInfoID = curSong.ID

	// Check if the user table exist
	if err := r.DB.Migrator().HasTable(user); !err {
		if err := r.DB.Debug().Migrator().CreateTable(user); err != nil {
			return user, err
		}
	}

	// Try to find the user entity in database
	if err := r.DB.First(user, "id = ?", user.ID).Error; err != nil {
		if err := r.DB.Create(user).Error; err != nil {
			return user, err
		}
		return user, nil
	} else {
		return user, err
	}
}

// The Update method updates an existing user.
func (r *User) Update(id uint, user *entity.User) (*entity.User, error) {
	// Update song entity in database
	if err := r.DB.Model(&entity.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// The Delete method deletes a user from the database.
func (r *User) Delete(id uint) error {
	// Try to find the user entity in database
	var user entity.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return err
	}

	// Delete finded user
	if err := r.DB.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
