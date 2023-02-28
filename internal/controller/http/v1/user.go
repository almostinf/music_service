package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
		h.POST("/", r.create)
		h.PUT("/", r.update)
		h.DELETE("/", r.delete)
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

func (r *userRoutes) create(c *gin.Context) {
	var request addUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := r.u.Create(&entity.User{
		Model:     gorm.Model{},
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
		CurSong:   nil,
		Playlists: make([]*entity.Playlist, 0),
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("user service problems: %s", err))
		return
	}

	c.JSON(http.StatusCreated, user)
}

type updateUserRequest struct {
	ID        string `json:"id" binding:"required" example:"0"`
	FirstName string `json:"first_name" binding:"required" example:"test"`
	LastName  string `json:"last_name" binding:"required" example:"test"`
	Email     string `json:"email" binding:"required" example:"test@test.ru"`
	Password  string `json:"password" binding:"required" example:"testtest"`
}

func (r *userRoutes) update(c *gin.Context) {
	var request updateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := r.u.Update(request.ID, &entity.User{
		Model:     gorm.Model{},
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("user service problems: %s", err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *userRoutes) delete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.u.Delete(id); err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("user service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}
