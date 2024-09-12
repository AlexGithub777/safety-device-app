package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (a *App) initRoutes() {
	// Public routes
	a.Router.GET("/", a.HandleGetLogin)
	a.Router.GET("/register", a.HandleGetRegister)
	a.Router.GET("/forgot-password", a.HandleGetForgotPassword)
	a.Router.POST("/forgot-password", a.HandlePostForgotPassword)
	a.Router.POST("/register", a.HandlePostRegister)
	a.Router.POST("/login", a.HandlePostLogin)
	a.Router.GET("/logout", a.HandleGetLogout)

	// Protected routes
	protected := a.Router.Group("")
	protected.Use(JWTMiddleware)

	protected.GET("/dashboard", a.HandleGetDashboard)
	protected.GET("/admin", a.HandleGetAdmin)
	// Add the rest of your protected routes here

	// API routes
	api := protected.Group("/api")

	api.GET("/emergency-device", a.HandleGetAllDevices)
	api.GET("/emergency-device-type", a.HandleGetAllDeviceTypes)
	api.GET("/extinguisher-type", a.HandleGetAllExtinguisherTypes)
	api.GET("/room", a.HandleGetAllRooms)
	api.GET("/building", a.HandleGetAllBuildings)
	api.GET("/site", a.HandleGetAllSites)
	//a.Router.POST("/api/emergency-device", a.HandleAddDevice)
	// Add the rest of your API routes here
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing jwt token")
			}
			return echo.NewHTTPError(http.StatusBadRequest, "invalid jwt token")
		}
		tokenString := cookie.Value

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Token is valid, you can access claims here if needed
			// For example: username := claims["name"].(string)
			c.Set("user", claims)
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt token")
	}
}
