package app

func (a *App) initRoutes() {
	// Web routes
	a.Router.GET("/", a.HomeHandler)
	// Routes
	a.Router.GET("/fire-extinguishers", a.FireExtinguisherHandler)

	// API routes
	a.Router.GET("/api/fire-extinguishers-html", a.GetFireExtinguishersHTML)
}
