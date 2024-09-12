package app

import (
	"fmt"
	"net/http"

	"github.com/AlexGithub777/safety-device-app/internal/models"
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
	// Get the form values
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	fmt.Println(username, email, password)

	// Create the user
	user := models.User{
		Username: username,
		Password: password,
		Email:    email,
	}

	// Insert the user into the database
	err := a.DB.CreateUser(&user)
	if err != nil {
		return err
	}

	// Redirect to the login page
	return c.Redirect(http.StatusSeeOther, "/")
}

// HandlePostLogin handles the login form submission
func (a *App) HandlePostLogin(c echo.Context) error {
	return nil
}

// HandleGetLogout logs the user out
func (a *App) HandleGetLogout(c echo.Context) error {
	return nil
}
