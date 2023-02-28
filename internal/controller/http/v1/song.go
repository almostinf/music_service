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

type songRoutes struct {
	u interfaces.Song
}

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

type songsResponse struct {
	Songs []entity.Song `json:"songs"`
}

func (r *songRoutes) get(c *gin.Context) {
	songs := r.u.Get()
	c.JSON(http.StatusOK, songsResponse{songs})
}

type addSongRequest struct {
	Title    string        `json:"title" binding:"required" example:"world"`
	Artist   string        `json:"artist" binding:"required" example:"hello"`
	Duration time.Duration `json:"duration" binding:"required" example:"01:02:05"`
}

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
		Playlists: make([]*entity.Playlist, 0),
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("song service problems: %s", err))
		return
	}

	c.JSON(http.StatusCreated, song)
}

type updateSongRequest struct {
	ID       uint          `json:"id" binding:"required" example:"0"`
	Title    string        `json:"title" binding:"required" example:"world"`
	Artist   string        `json:"artist" binding:"required" example:"hello"`
	Duration time.Duration `json:"duration" binding:"required" example:"01:02:05"`
}

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
