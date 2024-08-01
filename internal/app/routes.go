package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HomeHandler)
	// Routes
	a.Router.GET("/fire-extinguishers", a.FireExtinguisherHandler)

	// API routes
	a.Router.GET("/fire-extinguishers/data", a.GetFireExtinguishersHTML)
	a.Router.POST("/api/fire-extinguishers", a.CreateFireExtinguisher)
	a.Router.GET("/fire-extinguishers/pagination", a.GetPaginationControls)

}
