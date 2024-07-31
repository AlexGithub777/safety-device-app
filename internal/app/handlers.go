package app

import (
	"fmt"
	"log"
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

func (a *App) CreateFireExtinguisher(c echo.Context) error {
	// Parse form data
	room := c.FormValue("room")
	fireExtinguisherTypeID := c.FormValue("fire_extinguisher_type_id")
	serialNumber := c.FormValue("serial_number")
	dateOfManufacture := c.FormValue("date_of_manufacture")
	expireDate := c.FormValue("expire_date")
	size := c.FormValue("size")
	misc := c.FormValue("misc")
	status := c.FormValue("status")

	// Validate input
	if room == "" || fireExtinguisherTypeID == "" || serialNumber == "" ||
		dateOfManufacture == "" || expireDate == "" || size == "" || status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "All fields are required"})
	}

	// Insert new safety device
	var safetyDeviceID int
	err := a.DB.QueryRow(`
        INSERT INTO safety_devices (safety_device_type, room_id)
        VALUES ($1, $2) RETURNING safety_device_id`,
		"Fire Extinguisher", room).Scan(&safetyDeviceID)
	if err != nil {
		log.Println("Error inserting safety device:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating safety device"})
	}

	// Insert new fire extinguisher
	var fireExtinguisherID int
	err = a.DB.QueryRow(`
        INSERT INTO fire_extinguishers (
            safety_device_id, 
            fire_extinguisher_type_id, 
            serial_number, 
            date_of_manufacture, 
            expire_date, 
            size, 
            misc, 
            status
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        ) RETURNING fire_extinguisher_id
    `, safetyDeviceID, fireExtinguisherTypeID, serialNumber, dateOfManufacture, expireDate, size, misc, status).
		Scan(&fireExtinguisherID)
	if err != nil {
		log.Println("Error inserting fire extinguisher:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating fire extinguisher"})
	}

	// Fetch the newly created fire extinguisher details
	var newFireExtinguisher models.FireExtinguisher
	err = a.DB.QueryRow(`
        SELECT safety_device_id, fire_extinguisher_type_id, serial_number, date_of_manufacture, expire_date, size, misc, status
        FROM fire_extinguishers
        WHERE fire_extinguisher_id = $1`,
		fireExtinguisherID).Scan(
		&newFireExtinguisher.SafetyDeviceID,
		&newFireExtinguisher.FireExtinguisherTypeID,
		&newFireExtinguisher.SerialNumber,
		&newFireExtinguisher.DateOfManufacture,
		&newFireExtinguisher.ExpireDate,
		&newFireExtinguisher.Size,
		&newFireExtinguisher.Misc,
		&newFireExtinguisher.Status,
	)
	if err != nil {
		log.Println("Error fetching newly created fire extinguisher:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching new fire extinguisher"})
	}

	// Build HTML for the new row
	newRowHTML := fmt.Sprintf(`
        <tr>
            <td>%d</td>
            <td>%d</td>
            <td>%d</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
        </tr>`,
		newFireExtinguisher.FireExtinguisherID,
		newFireExtinguisher.SafetyDeviceID,
		newFireExtinguisher.FireExtinguisherTypeID,
		newFireExtinguisher.SerialNumber,
		newFireExtinguisher.DateOfManufacture,
		newFireExtinguisher.ExpireDate,
		newFireExtinguisher.Size,
		newFireExtinguisher.Misc,
		newFireExtinguisher.Status,
	)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Fire extinguisher created successfully.",
		"rowHTML": newRowHTML,
	})
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
