package http

import (
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
	"github.com/gin-gonic/gin"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) *UserUseCase {
	return &UserUseCase{r}
}

func (u *UserUseCase) GetUsers(c *gin.Context) {
	c.JSON(200, u.repo.Get())
}

func (u *UserUseCase) CreateUser(c *gin.Context) {
	user := new(entity.User)

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	user, err = u.repo.Add(user)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, user)
}
