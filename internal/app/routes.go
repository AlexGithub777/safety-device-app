package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HomeHandler)
	// Routes
	a.Router.GET("/fire-extinguishers", a.FireExtinguisherHandler)
	a.Router.GET("/map", a.MapHandler)

	// API routes
	a.Router.GET("/fire-extinguishers/data", a.GetFireExtinguishersHTML)
	a.Router.POST("/api/fire-extinguishers", a.CreateFireExtinguisher)
	a.Router.GET("/fire-extinguishers/pagination", a.GetPaginationControls)
	a.Router.GET("/fire-extinguishers/:id", a.SingleFireExtinguisherHandler)
	a.Router.GET("/buildings/:id", a.SingleBuildingHandler)

}
