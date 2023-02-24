package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/interfaces"
	user_usecase "github.com/almostinf/music_service/internal/usecase"
)

type userRoutes struct {
	u interfaces.User
}

func newUserRoutes(handler *gin.RouterGroup, u interfaces.User) {
		
}