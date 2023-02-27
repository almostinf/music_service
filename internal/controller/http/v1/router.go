package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/almostinf/music_service/internal/infrastructure/interfaces"
)

func NewRouter(handler *gin.Engine, u interfaces.User, s interfaces.Song) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Routers
	h := handler.Group("/v1")
	{
		newUserRoutes(h, u)
		newSongRoutes(h, s)
	}
}
