package app

import (
	"github.com/AlexGithub777/safety-device-app/internal/config"
	"github.com/AlexGithub777/safety-device-app/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	DB     *database.DB
	Router *echo.Echo
}

func NewApp(cfg config.Config) *App {
	// Initialize Echo
	router := echo.New()

	// Middleware
	router.Use(middleware.Logger())  // Log requests
	router.Use(middleware.Recover()) // Recover from panics
	router.Use(middleware.CORS())    // Enable CORS

	// Initialize Database
	db, err := database.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	app := &App{
		DB:     db,
		Router: router,
	}

	// Initialize routes
	app.initRoutes()

	return app
}
