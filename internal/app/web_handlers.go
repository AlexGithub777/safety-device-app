package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HomeHandler serves the home page
func (a *App) HandleHome(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

// HandleDashboard serves the dashboard page
func (a *App) HandleDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard.html", nil)
}

// HandleAdmin serves the admin page
func (a *App) HandleAdmin(c echo.Context) error {
	return c.Render(http.StatusOK, "admin.html", nil)
}
