package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HomeHandler)
	// Routes
	a.Router.GET("/dashboard", a.HandleDashboard)

	// API routes
	a.Router.GET("/api/emergency-device", a.GetAllDevices)
	a.Router.POST("/api/emergency-device", a.HandleAddDevice)

}
