package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HandleGetLogin)
	a.Router.GET("/dashboard", a.HandleGetDashboard)
	a.Router.GET("/admin", a.HandleGetAdmin)
	a.Router.GET("/register", a.HandleGetRegister)
	a.Router.GET("/login", a.HandleGetLogin)
	a.Router.GET("/forgot-password", a.HandleGetForgotPassword)

	// API routes
	a.Router.GET("/api/emergency-device", a.HandleGetAllDevices)
	a.Router.GET("/api/emergency-device-type", a.HandleGetAllDeviceTypes)
	a.Router.GET("/api/extinguisher-type", a.HandleGetAllExtinguisherTypes)
	a.Router.GET("/api/room", a.HandleGetAllRooms)
	a.Router.GET("/api/building", a.HandleGetAllBuildings)
	a.Router.GET("/api/site", a.HandleGetAllSites)
	//a.Router.POST("/api/emergency-device", a.HandleAddDevice)
}
