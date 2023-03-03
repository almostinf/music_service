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
	playlist_rep := repository.NewPlaylist(db)

	// UseCase
	userUseCase := usecase.NewUser(*user_rep)
	songUseCase := usecase.NewSong(*song_rep)
	playlistUseCase := usecase.NewPlaylist(*playlist_rep)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, userUseCase, songUseCase, playlistUseCase)
	if err := handler.Run(":8080"); err != nil {
		log.Fatal("Can not run server on specified port")
	}

	// Migrations
	if err := db.Debug().AutoMigrate(&entity.User{}, &entity.Song{}, &entity.Playlist{}, &entity.SongNode{}, &entity.CurSongInfo{}); err != nil {
		log.Fatal("Can not automigrate")
	}
}
