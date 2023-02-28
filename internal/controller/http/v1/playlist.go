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

type playlistRoutes struct {
	u interfaces.Playlist
}

func newplaylistRoutes(handler *gin.RouterGroup, u interfaces.Playlist) {
	r := &playlistRoutes{u}

	h := handler.Group("/playlist")
	{
		h.GET("/", r.get)
		h.POST("/", r.create)
		h.PUT("/", r.update)
		h.DELETE("/", r.delete)
	}
}

type playlistsResponse struct {
	Playlists []entity.Playlist `json:"playlists"`
}

func (r *playlistRoutes) get(c *gin.Context) {
	playlists := r.u.Get()
	c.JSON(http.StatusOK, playlistsResponse{playlists})
}

type addPlaylistRequest struct {
	Name      string `json:"name" binding:"required" example:"hello"`
	Title     string `json:"title" binding:"required" example:"world"`
	CreatorID uint   `json:"creator_id" binding:"required" example:"0"`
}

func (r *playlistRoutes) create(c *gin.Context) {
	var request addPlaylistRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	playlist, err := r.u.Create(&entity.Playlist{
		Model:     gorm.Model{},
		Name:      request.Name,
		Title:     request.Title,
		HeadSong:  nil,
		TailSong:  nil,
		CreatorID: request.CreatorID,
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.JSON(http.StatusCreated, playlist)
}

type updatePlaylistRequest struct {
	ID    uint   `json:"id" binding:"required" example:"1"`
	Name  string `json:"name" binding:"required" example:"hello"`
	Title string `json:"title" binding:"required" example:"world"`
}

func (r *playlistRoutes) update(c *gin.Context) {
	var request updatePlaylistRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	playlist, err := r.u.Update(request.ID, &entity.Playlist{
		Model: gorm.Model{},
		Name:  request.Name,
		Title: request.Title,
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.JSON(http.StatusOK, playlist)
}

func (r *playlistRoutes) delete(c *gin.Context) {
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
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}
