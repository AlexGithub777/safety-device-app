package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/config"
	"github.com/AlexGithub777/safety-device-app/internal/models"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type DB struct {
	*sql.DB
}

func NewDB(cfg config.Config) (*DB, error) {
	log.Println("Connecting to database...")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure connection is established
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")

	return &DB{db}, nil
}

func (db *DB) FetchAllDevices(buildingCode string) ([]models.EmergencyDevice, error) {
	var query string
	var args []interface{}

	// Define the base query
	query = `
    SELECT 
        ed.emergencydeviceid, 
        edt.emergencydevicetypename,
        et.extinguishertypename AS ExtinguisherTypeName,
        r.name AS RoomName,
        ed.serialnumber,
        ed.manufacturedate,
        ed.lastinspectiondate,
        ed.description,
        ed.size,
        ed.status 
    FROM emergency_deviceT ed
    JOIN roomT r ON ed.roomid = r.roomid
    LEFT JOIN emergency_device_typeT edt ON ed.emergencydevicetypeid = edt.emergencydevicetypeid
    LEFT JOIN Extinguisher_TypeT et ON ed.extinguishertypeid = et.extinguishertypeid
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
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var emergencyDevices []models.EmergencyDevice

	// Scan the results
	for rows.Next() {
		var device models.EmergencyDevice
		err := rows.Scan(
			&device.EmergencyDeviceID,
			&device.EmergencyDeviceTypeName,
			&device.ExtinguisherTypeName,
			&device.RoomName,
			&device.SerialNumber,
			&device.ManufactureDate,
			&device.LastInspectionDate,
			&device.Description,
			&device.Size,
			&device.Status,
		)
		if err != nil {
			return nil, err
		}

		// If any of the following fields are null, replace them with a default value
		if !device.ExtinguisherTypeName.Valid {
			device.ExtinguisherTypeName.String = "N/A"
		}
		if !device.SerialNumber.Valid {
			device.SerialNumber.String = "N/A"
		}
		if !device.Description.Valid {
			device.Description.String = "N/A"
		}
		if !device.Size.Valid {
			device.Size.String = "N/A"
		}
		if !device.Status.Valid {
			device.Status.String = "N/A"
		}

		// Handle dates and calculate expiry and next inspection dates
		if device.ManufactureDate.Valid {
			expiryDate := device.ManufactureDate.Time.AddDate(5, 0, 0)
			device.ExpireDate = sql.NullTime{
				Time:  expiryDate,
				Valid: true,
			}
		} else {
			device.ManufactureDate = sql.NullTime{
				Time:  time.Time{}, // Zero value of time.Time
				Valid: false,
			}
			device.ExpireDate = sql.NullTime{
				Time:  time.Time{}, // Zero value of time.Time
				Valid: false,
			}
		}

		if device.LastInspectionDate.Valid {
			nextInspectionDate := device.LastInspectionDate.Time.AddDate(0, 3, 0)
			device.NextInspectionDate = sql.NullTime{
				Time:  nextInspectionDate,
				Valid: true,
			}
		} else {
			device.NextInspectionDate = sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			}
		}



		emergencyDevices = append(emergencyDevices, device)
	}

	return emergencyDevices, nil
}