package database

import (
	"database/sql"
	"fmt"
	"log"

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

func (db *DB) FetchAllDevices(buildingCode string) ([]struct {
    models.EmergencyDevice
}, error) {
    var query string
    var args []interface{}

    // Define the base query
    query = `
    SELECT 
    ed.emergencydeviceid, 
    edt.emergencydevicetypename,
    ed.extinguishertype,
    r.name as RoomName,
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
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Define the result slice
    emergencyDevices := []struct {
        models.EmergencyDevice
    }{}

    // Scan the results
    for rows.Next() {
        var emergencyDevice struct {
            models.EmergencyDevice
        }
        if err := rows.Scan(
            &emergencyDevice.EmergencyDeviceID,
            &emergencyDevice.EmergencyDeviceTypeName, 
            &emergencyDevice.ExtinguisherType,
            &emergencyDevice.RoomName,
            &emergencyDevice.SerialNumber,
            &emergencyDevice.ManufactureDate,
            &emergencyDevice.LastInspectionDate,
            &emergencyDevice.Description,
            &emergencyDevice.Size,
            &emergencyDevice.Status,
        ); err != nil {
            return nil, err
        }
        emergencyDevices = append(emergencyDevices, emergencyDevice)
    }

    return emergencyDevices, nil
}
