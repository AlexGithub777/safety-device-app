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
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)

// CustomClaims represents JWT custom claims
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// HandlePostForgotPassword handles the forgot password form submission
func (a *App) HandlePostForgotPassword(c echo.Context) error {
	email := c.FormValue("email")

	// Check if the email exists in the database
	user, err := a.DB.GetUserByEmail(email)
	if err != nil {
		return c.Render(http.StatusOK, "forgot_password.html", map[string]interface{}{
			"error": "Email not found",
		})
	}

	// Generate a new password
	newPassword, err := password.Generate(15, 10, 5, false, false)
	if err != nil {
		return c.Render(http.StatusOK, "forgot_password.html", map[string]interface{}{
			"error": "Could not generate password",
		})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password
	if err := a.DB.UpdatePassword(user.UserID, string(hashedPassword)); err != nil {
		return err
	}

	// Send the new password to the user's email
	if err := sendPasswordResetEmail(email, user.Username, newPassword); err != nil {
		return err
	}

	// Render the forgot password page with a success message
	return c.Render(http.StatusOK, "forgot_password.html", map[string]interface{}{
		"message": fmt.Sprintf("Password reset successful. Check your %s for the new password.", email),
	})
}

// HandlePostRegister handles the register form submission
func (a *App) HandlePostRegister(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Check if the user or email already exists
	if _, err := a.DB.GetUserByUsername(username); err == nil {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Username already exists",
		})
	}

	if _, err := a.DB.GetUserByEmail(email); err == nil {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Email already exists",
		})
	}

	// Create a new user
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := a.DB.CreateUser(&user); err != nil {
		return err
	}

	// Generate a success message
	message := fmt.Sprintf("Registration successful. Please login with your username: %s", username)
	return c.Render(http.StatusOK, "register.html", map[string]interface{}{
		"message": message,
	})
}

// HandlePostLogin handles user login
func (a *App) HandlePostLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	remember := c.FormValue("remember")

	// Validate the user's credentials
	user, err := a.DB.GetUserByUsername(username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"error": "Invalid username or password",
		})
	}

	// Determine expiration time based on "remember" checkbox
	expiresAt := time.Now().Add(72 * time.Hour)
	if remember == "on" {
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
	}

	// Generate token
	token, err := GenerateToken(user, expiresAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	// Set the token as a cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
	}
	c.SetCookie(cookie)

	// Return the token in the response as well
	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

// HandleGetLogout logs the user out
func (a *App) HandleGetLogout(c echo.Context) error {
	// Clear JWT or session cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	// Redirect the user to the login page
	return c.Redirect(http.StatusSeeOther, "/")
}

// HandleGetLogin serves the home page
func (a *App) HandleGetLogin(c echo.Context) error {
	// Check if the user is already logged in
	cookie, err := c.Cookie("token")
	if err == nil && cookie.Value != "" {
		// Parse the JWT token
		token, err := parseToken(cookie.Value)
		if err == nil && token.Valid {
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

// GenerateToken generates a JWT token
func GenerateToken(user *models.User, expiresAt time.Time) (string, error) {
	claims := &CustomClaims{
		UserID:   strconv.Itoa(user.UserID),
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

// parseToken parses and validates the JWT token
func parseToken(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// sendPasswordResetEmail sends a password reset email
func sendPasswordResetEmail(email, username, newPassword string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "alexscott200020@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "EDMS PASSWORD RESET")
	m.SetBody("text/plain", "Your Username is "+username+", Your new password is: "+newPassword)
	m.AddAlternative("text/html", `<html><body style="font-family: Arial, sans-serif; padding: 20px;">
		<h2 style="color: #333;">EDMS PASSWORD RESET</h2>
		<p style="margin-top: 20px;">Your Username is <strong>`+username+`</strong></p>
		<p>Your new password is: <strong>`+newPassword+`</strong></p>
	</body></html>`)

	d := gomail.NewDialer("smtp.gmail.com", 587, "alexscott200020@gmail.com", "rmua arvp tedv rvlr")
	return d.DialAndSend(m)
}
