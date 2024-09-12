package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// HandleGetLogin serves the home page
func (a *App) HandleGetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

// HandleGetDashboard serves the dashboard page
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	// Create a new user
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Check if the user already exists
	existingUser, err := a.DB.GetUserByUsername(username)
	if err == nil {
		return fmt.Errorf("user %s already exists", existingUser.Username)
	}

	// Create the user
	err = a.DB.CreateUser(&user)
	if err != nil {
		return err
	}

	// Redirect to the login page
	return c.Redirect(http.StatusSeeOther, "/")
}

// HandlePostLogin handles the login form submission
func (a *App) HandlePostLogin(c echo.Context) error {
	// Get the form values
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Validate the user's credentials
	user, err := a.DB.GetUserByUsername(username)
	if err != nil {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"error": "Invalid username or password",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"error": "Invalid username or password",
		})
	}

	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.UserID
	claims["name"] = user.Username
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Get the token secret from the environment
	secret := os.Getenv("JWT_SECRET")

	// Generate encoded token
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	// Set the token as a cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true // Prevents JavaScript access to cookie
	cookie.Secure = false  // Requires HTTPS (remove this line if not using HTTPS)
	c.SetCookie(cookie)

	// Redirect to the dashboard
	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

// HandleGetLogout logs the user out
func (a *App) HandleGetLogout(c echo.Context) error {
	// Clear JWT or session cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour) // Expire the cookie
	c.SetCookie(cookie)

	// Redirect the user to the login page
	return c.Redirect(http.StatusSeeOther, "/")
}
