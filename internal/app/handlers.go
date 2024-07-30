package app

import (
	"html/template"
	"net/http"

	"github.com/AlexGithub777/safety-device-app/internal/models"
	"github.com/labstack/echo/v4"
)

func (a *App) HomeHandler(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	return tmpl.Execute(c.Response().Writer, nil)
}

func (a *App) DevicesHandler(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles("templates/devices.html"))
	return tmpl.Execute(c.Response().Writer, nil)
}

func (a *App) GetDevices(c echo.Context) error {
	devices := []models.Device{}
	rows, err := a.DB.Query("SELECT id, name, status FROM devices")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var device models.Device
		if err := rows.Scan(&device.ID, &device.Name, &device.Status); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		devices = append(devices, device)
	}
	return c.JSON(http.StatusOK, devices)
}
