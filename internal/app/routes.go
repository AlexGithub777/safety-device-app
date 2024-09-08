package app

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

}
