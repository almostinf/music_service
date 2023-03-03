package main

import (
	"log"

	"github.com/almostinf/music_service/config"
	"github.com/almostinf/music_service/internal/app"
)

func main() {
	// Configuration
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
		return
	}

	// Run
	app.Run(config)
}
