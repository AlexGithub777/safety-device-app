package main

import (
	"log"

	"github.com/AlexGithub777/safety-device-app/internal/app"
	"github.com/AlexGithub777/safety-device-app/internal/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the app
	application := app.NewApp(cfg)

	// Start the app
	log.Fatal(application.Router.Start(":8080"))
}
