package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleGetLogin serves the home page
func (a *App) HandleGetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

// HandleGetDashboard serves the dashboard page
func (a *App) HandleGetDashboard(c echo.Context) error {

	return c.Render(http.StatusOK, "dashboard.html", nil)
}

// HandleGetAdmin serves the admin page
func (a *App) HandleGetAdmin(c echo.Context) error {
	return c.Render(http.StatusOK, "admin.html", nil)
}

// HandleGetRegister serves the register page
func (a *App) HandleGetRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

// HandleGetForgotPassword serves the forgot password page
func (a *App) HandleGetForgotPassword(c echo.Context) error {
	return c.Render(http.StatusOK, "forgot_password.html", nil)
}

// HandlePostForgotPassword handles the forgot password form submission
func (a *App) HandlePostForgotPassword(c echo.Context) error {
	return nil
}

// HandlePostRegister handles the register form submission
func (a *App) HandlePostRegister(c echo.Context) error {
	return nil
}

// HandlePostLogin handles the login form submission
func (a *App) HandlePostLogin(c echo.Context) error {
	return nil
}

// HandleGetLogout logs the user out
func (a *App) HandleGetLogout(c echo.Context) error {
	return nil
}
