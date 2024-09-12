package app

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
)

func (a *App) initRoutes() {
	// Public routes
	a.Router.GET("/", a.HandleGetLogin)
	a.Router.GET("/register", a.HandleGetRegister)
	a.Router.GET("/forgot-password", a.HandleGetForgotPassword)
	a.Router.POST("/forgot-password", a.HandlePostForgotPassword)
	a.Router.POST("/register", a.HandlePostRegister)
	a.Router.POST("/login", a.HandlePostLogin)
	a.Router.GET("/logout", a.HandleGetLogout)

	// Get JWT secret from environment variable
	secret := os.Getenv("JWT_SECRET")

	// Protected routes
	protected := a.Router.Group("")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),
	}))

	protected.GET("/dashboard", a.HandleGetDashboard)
	protected.GET("/admin", a.HandleGetAdmin)
	// Add the rest of your protected routes here

	// API routes
	api := protected.Group("/api")

	api.GET("/emergency-device", a.HandleGetAllDevices)
	api.GET("/emergency-device-type", a.HandleGetAllDeviceTypes)
	api.GET("/extinguisher-type", a.HandleGetAllExtinguisherTypes)
	api.GET("/room", a.HandleGetAllRooms)
	api.GET("/building", a.HandleGetAllBuildings)
	api.GET("/site", a.HandleGetAllSites)
	//a.Router.POST("/api/emergency-device", a.HandleAddDevice)
	// Add the rest of your API routes here
}
