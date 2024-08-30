package app

import (
	"database/sql"
	"fmt"
	"log"
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


func (a *App) HandleDashboard(c echo.Context) error {
	// Render the root template without pagination data
	return c.Render(http.StatusOK, "fire_extinguishers.html", nil)
}


// GetAllDevices fetches all emergency devices from the database with optional filtering by building code
func (a *App) GetAllDevices(c echo.Context) error {
    buildingCode := c.QueryParam("building_code")
    var query string
    var args []interface{}

    // Define the base query
    query = `
	SELECT 
	ed.emergencydeviceid, 
	edt.emergencydevicetypename,
	ed.extinguishertype,
	r.name as Room,
	ed.serialnumber,
	ed.manufacturedate,
	ed.lastinspectiondate,
	ed.description,
	ed.size,
	ed.status 
	FROM emergency_deviceT ed
	JOIN roomT r ON ed.roomid = r.roomid
	JOIN emergency_device_typeT edt ON ed.emergencydevicetypeid = edt.emergencydevicetypeid
    `

    // Add filtering by building code if provided
    if buildingCode != "" {
        query += `
            JOIN buildingT b ON r.buildingid = b.buildingid
            WHERE b.buildingcode = $1
        `
        args = append(args, buildingCode)
    }

    // Prepare and execute the query
    rows, err := a.DB.Query(query, args...)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error fetching data: %v", err))
    }
    defer rows.Close()

    // Define the result slice
	emergencyDevices := []struct {
		models.EmergencyDevice
		EmergencyDeviceTypeName string `json:"emergency_device_type_name"`
		Room           string `json:"room"`
		ExtinguisherType string `json:"extinguisher_type"`
	}{}

	// Scan the results
	for rows.Next() {
		var emergencyDevice struct {
			models.EmergencyDevice
			EmergencyDeviceTypeName string `json:"emergency_device_type_name"`
			Room            string `json:"room"`
			ExtinguisherType string `json:"extinguisher_type"`
		}
		if err := rows.Scan(
			&emergencyDevice.EmergencyDeviceID,
			&emergencyDevice.EmergencyDeviceTypeName, // Scan into the new field
			&emergencyDevice.ExtinguisherType,
			&emergencyDevice.Room,
			&emergencyDevice.SerialNumber,
			&emergencyDevice.ManufactureDate,
			&emergencyDevice.LastInspectionDate,
			&emergencyDevice.Description,
			&emergencyDevice.Size,
			&emergencyDevice.Status,
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error scanning data: %v", err))
		}
		emergencyDevices = append(emergencyDevices, emergencyDevice)
	}

    // Return the results as JSON
    return c.JSON(http.StatusOK, emergencyDevices)
}

func (a *App) HandleAddDevice(c echo.Context) error {
    // Parse form data
    roomStr := c.FormValue("room")
    emergencyDeviceTypeIDStr := c.FormValue("emergency_device_type_id")
    serialNumber := c.FormValue("serial_number")
    manufactureDateStr := c.FormValue("manufacture_date")
    lastInspectionDateStr := c.FormValue("last_inspection_date")
    size := c.FormValue("size")
    description := c.FormValue("description")
    status := c.FormValue("status")

    // Validate input
    if roomStr == "" || emergencyDeviceTypeIDStr == "" || serialNumber == "" ||
        manufactureDateStr == "" || size == "" || status == "" {
        return a.handleError(c, http.StatusBadRequest, "All fields are required", nil)
    }

    // Convert room ID and emergency device type ID to integers
    roomID, err := strconv.Atoi(roomStr)
    if err != nil {
        log.Printf("Error converting room to integer: %v", err)
        return a.handleError(c, http.StatusBadRequest, "Invalid room ID", err)
    }

    emergencyDeviceTypeID, err := strconv.Atoi(emergencyDeviceTypeIDStr)
    if err != nil {
        log.Printf("Error converting emergency device type ID to integer: %v", err)
        return a.handleError(c, http.StatusBadRequest, "Invalid emergency device type ID", err)
    }

    // Parse date strings into time.Time format
    manufactureDate, err := time.Parse("2006-01-02", manufactureDateStr)
    if err != nil {
        log.Printf("Error parsing manufacture date: %v", err)
        return a.handleError(c, http.StatusBadRequest, "Invalid manufacture date format", err)
    }

    // Optional: Parse last inspection date if provided
    var lastInspectionDate sql.NullTime
    if lastInspectionDateStr != "" {
        parsedDate, err := time.Parse("2006-01-02", lastInspectionDateStr)
        if err != nil {
            return a.handleError(c, http.StatusBadRequest, "Invalid last inspection date format", err)
        }
        lastInspectionDate = sql.NullTime{Time: parsedDate, Valid: true}
    }

    // Insert new emergency device
    var emergencyDeviceID int
    err = a.DB.QueryRow(`
        INSERT INTO emergency_devices (
            emergency_device_type_id, 
            room_id, 
            manufacture_date, 
            serial_number, 
            description, 
            size, 
            last_inspection_date, 
            status
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        ) RETURNING emergency_device_id
    `,
        emergencyDeviceTypeID,
        roomID,
        manufactureDate,
        serialNumber,
        description,
        size,
        lastInspectionDate,
        status).Scan(&emergencyDeviceID)
    if err != nil {
        return a.handleError(c, http.StatusInternalServerError, "Error creating emergency device", err)
    }

    // Create the new EmergencyDevice model
    newDevice := models.EmergencyDevice{
        EmergencyDeviceID:    emergencyDeviceID,
        EmergencyDeviceTypeID: emergencyDeviceTypeID,
        RoomID:               roomID,
        ManufactureDate:      manufactureDate,
        SerialNumber:         serialNumber,
        Description:          description,
        Size:                 size,
        LastInspectionDate:   &lastInspectionDate.Time, // only set if valid
        Status:               status,
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
        newDevice.EmergencyDeviceID,
        newDevice.EmergencyDeviceTypeID,
        newDevice.RoomID,
        newDevice.SerialNumber,
        newDevice.ManufactureDate.Format("02-01-2006"),
        newDevice.Size,
        newDevice.Description,
        newDevice.LastInspectionDate.Format("02-01-2006"), // ensure this is set correctly
        newDevice.Status,
    )

    // Return success message and the new row HTML
    return c.JSON(http.StatusOK, map[string]string{
        "message": "Emergency device created successfully.",
        "rowHTML": newRowHTML,
    })
}

