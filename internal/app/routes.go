package app

<<<<<<< HEAD
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
=======
func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HandleHome)
	a.Router.GET("/dashboard", a.HandleDashboard)
	a.Router.GET("/admin", a.HandleAdmin)
	a.Router.GET("/register", a.HandleRegister)
	a.Router.GET("/login", a.HandleHome)
	a.Router.GET("/forgot-password", a.HandleForgotPassword)

	// API routes
	a.Router.GET("/api/emergency-device", a.HandleGetAllDevices)
	a.Router.GET("/api/emergency-device-type", a.HandleGetAllDeviceTypes)
	a.Router.GET("/api/extinguisher-type", a.HandleGetAllExtinguisherTypes)
	a.Router.GET("/api/room", a.HandleGetAllRooms)
	a.Router.GET("/api/building", a.HandleGetAllBuildings)
	a.Router.GET("/api/site", a.HandleGetAllSites)
	//a.Router.POST("/api/emergency-device", a.HandleAddDevice)

>>>>>>> d3f1aef86552b4414e16af4df61e0e15859fe0b5
}
