package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// HandleGetLogin serves the home page
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (a *App) HandleGetLogin(c echo.Context) error {
	// Check if the user is already logged in
	cookie, err := c.Cookie("token")
	if err == nil && cookie.Value != "" {
		// Parse the JWT token
		token, err := jwt.ParseWithClaims(cookie.Value, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the token secret
			secret := os.Getenv("JWT_SECRET")
			return []byte(secret), nil
		})

		if err == nil && token.Valid {
			// Token is valid
			if claims, ok := token.Claims.(*CustomClaims); ok {
				// Put the claims data in the context
				c.Set("user", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)

				// User is already logged in, redirect to the dashboard
				return c.Redirect(http.StatusSeeOther, "/dashboard")
			}
		}
	}

	// If no valid token, render login page
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

func (a *App) HandlePostLogin(c echo.Context) error {
	// Get the form values
	username := c.FormValue("username")
	password := c.FormValue("password")
	remember := c.FormValue("remember") // Check if "remember" checkbox was checked

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

	// Determine the expiration time based on the "remember" checkbox
	var expiresAt time.Time
	if remember == "on" {
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
	} else {
		expiresAt = time.Now().Add(72 * time.Hour)
	}

	// Create custom claims
	claims := &CustomClaims{
		UserID:   strconv.Itoa(user.UserID),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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
	cookie.Path = "/"
	cookie.HttpOnly = true // Prevents JavaScript access to cookie
	cookie.Secure = true   // Requires HTTPS
	cookie.Expires = expiresAt

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
