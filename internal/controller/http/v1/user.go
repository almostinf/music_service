package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/interfaces"
)

type userRoutes struct {
	u interfaces.User
}

func newUserRoutes(handler *gin.RouterGroup, u interfaces.User) {
	r := &userRoutes{u}

	h := handler.Group("/user")
	{
		h.GET("/", r.get)
		h.POST("/", r.add)
	}
}

type usersResponse struct {
	Users []entity.User `json:"users"`
}

func (r *userRoutes) get(c *gin.Context) {
	users := r.u.Get()
	c.JSON(http.StatusOK, usersResponse{users})
}

type addUserRequest struct {
	FirstName string `json:"first_name" binding:"required" example:"test"`
	LastName  string `json:"last_name" binding:"required" example:"test"`
	Email     string `json:"email" binding:"required" example:"test@test.ru"`
	Password  string `json:"password" binding:"required" example:"testtest"`
}

func (r *userRoutes) add(c *gin.Context) {
	var request addUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := r.u.Create(&entity.User{
		Model:     gorm.Model{},
		ID:        uuid.New().String(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
		Playlists: nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("user service problems: %s", err))
		fmt.Println(user)
		return
	}

	c.JSON(http.StatusOK, user)
}
