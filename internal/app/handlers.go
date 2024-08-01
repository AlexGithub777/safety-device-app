package app

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/models"
	"github.com/labstack/echo/v4"
)

// HomeHandler serves the home page
func (a *App) HomeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func (a *App) FireExtinguisherHandler(c echo.Context) error {
	pageStr := c.QueryParam("page")
	sizeStr := c.QueryParam("size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 10
	}

	// Fetch total count of fire extinguishers for pagination
	var total int
	err = a.DB.QueryRow("SELECT COUNT(*) FROM fire_extinguishers").Scan(&total)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching count")
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(size)))

	// Render the root template with pagination data
	return c.Render(http.StatusOK, "fire_extinguishers.html", map[string]interface{}{
		"Page":       page,
		"Size":       size,
		"TotalPages": totalPages,
	})
}

func (a *App) CreateFireExtinguisher(c echo.Context) error {
	// Parse form data
	roomStr := c.FormValue("room")
	fireExtinguisherTypeIDStr := c.FormValue("fire_extinguisher_type_id")
	serialNumber := c.FormValue("serial_number")
	dateOfManufactureStr := c.FormValue("date_of_manufacture")
	expireDateStr := c.FormValue("expire_date")
	size := c.FormValue("size")
	misc := c.FormValue("misc")
	status := c.FormValue("status")

	// Validate input
	if roomStr == "" || fireExtinguisherTypeIDStr == "" || serialNumber == "" ||
		dateOfManufactureStr == "" || expireDateStr == "" || size == "" || status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "All fields are required"})
	}

	// Convert room ID to integer
	roomID, err := strconv.Atoi(roomStr)
	if err != nil {
		log.Printf("Error converting room to integer: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid room ID"})
	}

	// Convert fire extinguisher type ID to integer
	fireExtinguisherTypeID, err := strconv.Atoi(fireExtinguisherTypeIDStr)
	if err != nil {
		log.Printf("Error converting fireExtinguisherTypeID to integer: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid fire extinguisher type ID"})
	}

	// Convert date strings to time.Time using yyyy-mm-dd format
	dateOfManufacture, err := time.Parse("2006-01-02", dateOfManufactureStr)
	if err != nil {
		log.Printf("Error parsing date of manufacture: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date of manufacture format"})
	}

	expireDate, err := time.Parse("2006-01-02", expireDateStr)
	if err != nil {
		log.Printf("Error parsing expire date: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expire date format"})
	}

	// Insert new safety device
	var safetyDeviceID int
	err = a.DB.QueryRow(`
        INSERT INTO safety_devices (safety_device_type, room_id)
        VALUES ($1, $2) RETURNING safety_device_id`,
		"Fire Extinguisher", roomID).Scan(&safetyDeviceID)
	if err != nil {
		log.Println("Error inserting safety device:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating safety device"})
	}

	// Create new FireExtinguisher instance
	newFireExtinguisher := models.FireExtinguisher{
		SafetyDeviceID:         safetyDeviceID,
		FireExtinguisherTypeID: fireExtinguisherTypeID,
		SerialNumber:           serialNumber,
		DateOfManufacture:      dateOfManufacture,
		ExpireDate:             expireDate,
		Size:                   size,
		Misc:                   misc,
		Status:                 status,
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
    `, newFireExtinguisher.SafetyDeviceID, newFireExtinguisher.FireExtinguisherTypeID, newFireExtinguisher.SerialNumber, newFireExtinguisher.DateOfManufacture, newFireExtinguisher.ExpireDate, newFireExtinguisher.Size, newFireExtinguisher.Misc, newFireExtinguisher.Status).
		Scan(&fireExtinguisherID)
	if err != nil {
		log.Println("Error inserting fire extinguisher:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating fire extinguisher"})
	}

	// Update the model with the new ID
	newFireExtinguisher.FireExtinguisherID = fireExtinguisherID

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
		newFireExtinguisher.DateOfManufacture.Format("02-01-2006"),
		newFireExtinguisher.ExpireDate.Format("02-01-2006"),
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
	pageStr := c.QueryParam("page")
	sizeStr := c.QueryParam("size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 10
	}

	// Fetch total count of fire extinguishers for pagination
	var total int
	err = a.DB.QueryRow("SELECT COUNT(*) FROM fire_extinguishers").Scan(&total)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching count")
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(size)))

	offset := (page - 1) * size

	fireExtinguishers := []models.FireExtinguisher{}
	rows, err := a.DB.Query(`
        SELECT fire_extinguisher_id, safety_device_id, fire_extinguisher_type_id, serial_number, date_of_manufacture, expire_date, size, misc, status 
        FROM fire_extinguishers 
        LIMIT $1 OFFSET $2`, size, offset)
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

	data := map[string]interface{}{
		"FireExtinguishers": fireExtinguishers,
		"Page":              page,
		"Size":              size,
		"TotalPages":        totalPages,
	}

	// Check if this is an htmx request
	if c.Request().Header.Get("HX-Request") != "" {
		// Render only the table rows and pagination controls
		return c.Render(http.StatusOK, "fire_extinguishers_table.html", data)
	}

	// Render the root template with pagination data and fire extinguishers data
	return c.Render(http.StatusOK, "fire_extinguishers.html", data)
}
