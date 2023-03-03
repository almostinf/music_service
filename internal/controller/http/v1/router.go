package v1

import (
	"github.com/almostinf/music_service/internal/infrastructure/interfaces"
	"github.com/gin-gonic/gin"
)

// Creates a new Router of user, song and playlist routes
func NewRouter(handler *gin.Engine, u interfaces.User, s interfaces.Song, p interfaces.Playlist) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Routers
	h := handler.Group("/v1")
	{
		newUserRoutes(h, u)
		newSongRoutes(h, s)
		newplaylistRoutes(h, p)
	}
}
