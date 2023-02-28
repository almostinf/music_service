package app

import (
	"log"

	"github.com/almostinf/music_service/config"
	v1 "github.com/almostinf/music_service/internal/controller/http/v1"
	"github.com/almostinf/music_service/internal/entity"
	"github.com/almostinf/music_service/internal/infrastructure/repository"
	"github.com/almostinf/music_service/internal/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run(cfg *config.Config) {
	// Database
	db, err := gorm.Open(postgres.Open(cfg.PG.URL), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect database")
	}

	// Repository
	user_rep := repository.NewUser(db)
	song_rep := repository.NewSong(db)

	// UseCase
	userUseCase := usecase.NewUser(*user_rep)
	songUseCase := usecase.NewSong(*song_rep)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, userUseCase, songUseCase)
	handler.Run(":8080")

	// Migrations
	if err := db.Debug().AutoMigrate(&entity.User{}, &entity.Song{}, &entity.Playlist{}); err != nil {
		log.Fatal("Could not automigrate")
	}
}
