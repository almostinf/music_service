// This is a Go package for defining routes using the Gin web framework.
package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/interfaces"
)

// The playlistRoutes struct is defined to hold a playlist service that
// implements the interfaces.Playlist interface.
type playlistRoutes struct {
	u interfaces.Playlist
}

// This function creates a new instance of playlistRoutes with DI
func newplaylistRoutes(handler *gin.RouterGroup, u interfaces.Playlist) {
	r := &playlistRoutes{u}

	h := handler.Group("/playlist")
	{
		h.GET("/", r.get)
		h.POST("/", r.create)
		h.PUT("/", r.update)
		h.DELETE("/", r.delete)
		h.PUT("/play", r.play)
		h.PUT("/pause", r.pause)
		h.PUT("/add_song", r.addSong)
		h.PUT("/next", r.next)
		h.PUT("/prev", r.prev)
	}
}

// A structure for validating a json.
type playlistsResponse struct {
	Playlists []entity.Playlist `json:"playlists"`
}

// The get handler returns a JSON response containing a list of all
// playlists fetched from the playlist service.
func (r *playlistRoutes) get(c *gin.Context) {
	playlists := r.u.Get()
	c.JSON(http.StatusOK, playlistsResponse{playlists})
}

// A structure for validating a json.
type addPlaylistRequest struct {
	Name      string `json:"name" binding:"required" example:"hello"`
	Title     string `json:"title" binding:"required" example:"world"`
	CreatorID uint   `json:"creator_id" binding:"required" example:"0"`
}

// The create handler creates a new playlist using the data from the
// request JSON body and returns the newly created playlist as a JSON
// response.
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
		HeadSong:  0,
		TailSong:  0,
		CreatorID: request.CreatorID,
	})

	if err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.JSON(http.StatusCreated, playlist)
}

// A structure for validating a json.
type updatePlaylistRequest struct {
	ID    uint   `json:"id" binding:"required" example:"1"`
	Name  string `json:"name" binding:"required" example:"hello"`
	Title string `json:"title" binding:"required" example:"world"`
}

// The update handler updates an existing playlist with new data from the
// request JSON body and returns the updated playlist as a JSON response.
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

// The delete handler deletes an existing playlist identified by the ID
// in the request query parameter.
func (r *playlistRoutes) delete(c *gin.Context) {
	id := c.Query("id")

	conv_id, err := strconv.Atoi(id)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "failed uint conversation")
	}

	if err := r.u.Delete(uint(conv_id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}

// An utility function required to retrieve the user, song and playlist
// IDs from query.
func (r *playlistRoutes) getUserSongPlaylist(c *gin.Context) (uint, uint, uint, error) {
	user_id := c.Query("user_id")
	song_id := c.Query("song_id")
	playlist_id := c.Query("playlist_id")

	conv_user_id, err := strconv.Atoi(user_id)
	if err != nil {
		return 0, 0, 0, errors.New("failed uint conversation")
	}

	conv_song_id, err := strconv.Atoi(song_id)
	if err != nil {
		return 0, 0, 0, errors.New("failed uint conversation")
	}

	conv_playlist_id, err := strconv.Atoi(playlist_id)
	if err != nil {
		return 0, 0, 0, errors.New("failed uint conversation")
	}

	return uint(conv_user_id), uint(conv_song_id), uint(conv_playlist_id), nil
}

// Play the song with the specified user, songs and playlist IDs.
func (r *playlistRoutes) play(c *gin.Context) {
	user_id, song_id, playlist_id, err := r.getUserSongPlaylist(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := r.u.Play(uint(user_id), uint(song_id), uint(playlist_id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}

// Pause the song with the specified user, songs and playlist IDs.
func (r *playlistRoutes) pause(c *gin.Context) {
	user_id, song_id, playlist_id, err := r.getUserSongPlaylist(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := r.u.Pause(uint(user_id), uint(song_id), uint(playlist_id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}

// Add the song to the playlist with the specified user, songs and playlist IDs.
func (r *playlistRoutes) addSong(c *gin.Context) {
	user_id, song_id, playlist_id, err := r.getUserSongPlaylist(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := r.u.AddSong(uint(user_id), uint(song_id), uint(playlist_id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}

// Play the next song with the specified user, songs and playlist IDs.
func (r *playlistRoutes) next(c *gin.Context) {
	user_id, song_id, playlist_id, err := r.getUserSongPlaylist(c)

	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := r.u.Next(uint(user_id), uint(song_id), uint(playlist_id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}

// Pause the previous song with the specified user, songs and playlist IDs.
func (r *playlistRoutes) prev(c *gin.Context) {
	user_id, song_id, playlist_id, err := r.getUserSongPlaylist(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := r.u.Prev(uint(user_id), uint(song_id), uint(playlist_id)); err != nil {
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("playlist service problems: %s", err))
		return
	}

	c.Status(http.StatusOK)
}
