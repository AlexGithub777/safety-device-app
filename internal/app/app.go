package app

import (
	"log"
	"os"

	"github.com/AlexGithub777/safety-device-app/internal/config"
	"github.com/AlexGithub777/safety-device-app/internal/database"
	"github.com/AlexGithub777/safety-device-app/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App holds the application state including database and router
type App struct {
    DB     *database.DB
    Router *echo.Echo
    Logger *log.Logger
}

// handleError is a method of App for handling errors
func (a *App) handleError(c echo.Context, statusCode int, message string, err error) error {
    a.Logger.Printf("Error: %v", err)  // Use the logger in the App struct
    return c.JSON(statusCode, map[string]string{"error": message})
}

// NewApp creates a new instance of App
func NewApp(cfg config.Config) *App {
    // Initialize Echo
    router := echo.New()

    // Set up renderer
    router.Renderer = utils.NewTemplateRenderer()

    // Middleware
    router.Use(middleware.Logger())  // Log requests
    router.Use(middleware.Recover()) // Recover from panics
    router.Use(middleware.CORS())    // Enable CORS

    // Initialize Database
    db, err := database.NewDB(cfg)
    if err != nil {
        panic(err)
    }

    // Initialize Logger
    logger := log.New(os.Stdout, "APP: ", log.LstdFlags)

    app := &App{
        DB:     db,
        Router: router,
        Logger: logger,
    }

    // Initialize routes
    app.initRoutes()

    return app
}
