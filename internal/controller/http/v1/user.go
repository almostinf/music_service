// This is a Go package for defining routes using the Gin web framework.
package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/interfaces"
)

// The userRoutes struct is defined to hold a user service that
// implements the interfaces.User interface.
type userRoutes struct {
	u interfaces.User
}

// This function creates a new instance of userRoutes with DI
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

// A structure for validating a json.
type usersResponse struct {
	Users []entity.User `json:"users"`
}

// The get handler returns a JSON response containing a list of all
// users fetched from the user service.
func (r *userRoutes) get(c *gin.Context) {
	users := r.u.Get()
	c.JSON(http.StatusOK, usersResponse{users})
}

// A structure for validating a json.
type addUserRequest struct {
	FirstName string `json:"first_name" binding:"required" example:"test"`
	LastName  string `json:"last_name" binding:"required" example:"test"`
	Email     string `json:"email" binding:"required" example:"test@test.ru"`
	Password  string `json:"password" binding:"required" example:"testtest"`
}

// The create handler creates a new user using the data from the
// request JSON body and returns the newly created user as a JSON
// response.
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
		CurSongInfoID:   0,
		Playlists: make([]int64, 0),
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("user service problems: %s", err))
		return
	}

	c.JSON(http.StatusCreated, user)
}

// A structure for validating a json.
type updateUserRequest struct {
	ID        uint   `json:"id" binding:"required" example:"0"`
	FirstName string `json:"first_name" binding:"required" example:"test"`
	LastName  string `json:"last_name" binding:"required" example:"test"`
	Email     string `json:"email" binding:"required" example:"test@test.ru"`
	Password  string `json:"password" binding:"required" example:"testtest"`
}

// The update handler updates an existing user with new data from the
// request JSON body and returns the updated song as a JSON response.
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

// The delete handler deletes an existing user identified by the ID
// in the request query parameter.
func (r *userRoutes) delete(c *gin.Context) {
	id := c.Query("id")

	conv_id, err := strconv.Atoi(id)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "failed uint conversation")
	}

	if id == "" {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.u.Delete(uint(conv_id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("user service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}
