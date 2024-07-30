package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HomeHandler)
	a.Router.GET("/devices", a.DevicesHandler)

	// API routes
	a.Router.GET("/api/devices", a.GetDevices)
}
