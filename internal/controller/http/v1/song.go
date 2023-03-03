// This is a Go package for defining routes using the Gin web framework.
package v1

import (
	"fmt"
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/interfaces"
)

// The songRoutes struct is defined to hold a song service that
// implements the interfaces.Song interface.
type songRoutes struct {
	u interfaces.Song
}

// This function creates a new instance of songRoutes with DI
func newSongRoutes(handler *gin.RouterGroup, u interfaces.Song) {
	r := &songRoutes{u}

	h := handler.Group("/song")
	{
		h.GET("/", r.get)
		h.POST("/", r.create)
		h.PUT("/", r.update)
		h.DELETE("/", r.delete)
	}
}

// A structure for validating a json.
type songsResponse struct {
	Songs []entity.Song `json:"songs"`
}

// The get handler returns a JSON response containing a list of all
// songs fetched from the song service.
func (r *songRoutes) get(c *gin.Context) {
	songs := r.u.Get()
	c.JSON(http.StatusOK, songsResponse{songs})
}

// A structure for validating a json
type addSongRequest struct {
	Title    string        `json:"title" binding:"required" example:"world"`
	Artist   string        `json:"artist" binding:"required" example:"hello"`
	Duration time.Duration `json:"duration" binding:"required" example:"01:02:05"`
}

// The create handler creates a new song using the data from the
// request JSON body and returns the newly created song as a JSON
// response.
func (r *songRoutes) create(c *gin.Context) {
	var request addSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	song, err := r.u.Create(&entity.Song{
		Model:     gorm.Model{},
		Title:     request.Title,
		Artist:    request.Artist,
		Duration:  request.Duration,
		Playlists: make([]entity.Playlist, 0),
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("song service problems: %s", err))
		return
	}

	c.JSON(http.StatusCreated, song)
}

// A structure for validating a json.
type updateSongRequest struct {
	ID       uint          `json:"id" binding:"required" example:"0"`
	Title    string        `json:"title" binding:"required" example:"world"`
	Artist   string        `json:"artist" binding:"required" example:"hello"`
	Duration time.Duration `json:"duration" binding:"required" example:"01:02:05"`
}

// The update handler updates an existing song with new data from the
// request JSON body and returns the updated song as a JSON response.
func (r *songRoutes) update(c *gin.Context) {
	var request updateSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	song, err := r.u.Update(request.ID, &entity.Song{
		Model:    gorm.Model{},
		Title:    request.Title,
		Artist:   request.Artist,
		Duration: request.Duration,
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("song service problems: %s", err))
		return
	}

	c.JSON(http.StatusOK, song)
}

// The delete handler deletes an existing song identified by the ID
// in the request query parameter.
func (r *songRoutes) delete(c *gin.Context) {
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
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("song service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}
