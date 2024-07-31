package app

import (
	"net/http"

	"github.com/AlexGithub777/safety-device-app/internal/models"
	"github.com/labstack/echo/v4"
)

// HomeHandler serves the home page
func (a *App) HomeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

// FireExtinguisherHandler serves the fire extinguishers page
func (a *App) FireExtinguisherHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "fire_extinguishers.html", nil)
}

// GetFireExtinguishersHTML returns fire extinguishers data as HTML
func (a *App) GetFireExtinguishersHTML(c echo.Context) error {
	fireExtinguishers := []models.FireExtinguisher{}
	rows, err := a.DB.Query("SELECT * FROM fire_extinguishers")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching data")
	}
	defer rows.Close()

	for rows.Next() {
		var fireExtinguisher models.FireExtinguisher
		if err := rows.Scan(
			&fireExtinguisher.FireExtinguisherID,
			&fireExtinguisher.SafetyDeviceID,
			&fireExtinguisher.FireExtinguisherTypeID,
			&fireExtinguisher.SerialNumber,
			&fireExtinguisher.DateOfManufacture,
			&fireExtinguisher.ExpireDate,
			&fireExtinguisher.Size,
			&fireExtinguisher.Misc,
			&fireExtinguisher.Status); err != nil {
			return c.String(http.StatusInternalServerError, "Error scanning data")
		}
		fireExtinguishers = append(fireExtinguishers, fireExtinguisher)
	}

	return c.Render(http.StatusOK, "fire_extinguishers_table.html", map[string]interface{}{
		"FireExtinguishers": fireExtinguishers,
	})
}
