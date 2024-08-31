package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HandleHome)
	a.Router.GET("/dashboard", a.HandleDashboard)

	// API routes
	a.Router.GET("/api/emergency-device", a.HandleGetAllDevices)
	//a.Router.POST("/api/emergency-device", a.HandleAddDevice)

}
