package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HomeHandler)
	// Routes
	a.Router.GET("/dashboard", a.HandleDashboard)
	a.Router.GET("/dashboard/pagination", a.GetPaginationControls)

	// API routes
	a.Router.GET("/api/safety-device", a.GetDeviceTableHTML)
	a.Router.POST("/api/safety-device", a.HandleAddDevice)
	a.Router.GET("/api/safety-device/:id", a.SingleFireExtinguisherHandler)
	a.Router.GET("/api/buildings/:id", a.SingleBuildingHandler)

}
