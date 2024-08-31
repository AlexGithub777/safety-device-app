package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HandleHome)
	a.Router.GET("/dashboard", a.HandleDashboard)

	// API routes
	a.Router.GET("/api/emergency-device", a.HandleGetAllDevices)
	a.Router.GET("/api/emergency-device-type", a.HandleGetAllDeviceTypes)
	a.Router.GET("/api/extinguisher-type", a.HandleGetAllExtinguisherTypes)
	a.Router.GET("/api/room", a.HandleGetAllRooms)
	a.Router.GET("/api/building", a.HandleGetAllBuildings)
	//a.Router.POST("/api/emergency-device", a.HandleAddDevice)

}
