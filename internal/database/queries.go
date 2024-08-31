package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/models"
)

func (db *DB) GetAllDevices(buildingCode string) ([]models.EmergencyDevice, error) {
	var query string
	var args []interface{}

	// Define the base query
	query = `
    SELECT 
        ed.emergencydeviceid, 
        edt.emergencydevicetypename,
        et.extinguishertypename AS ExtinguisherTypeName,
        r.roomcode,
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
			&device.RoomCode,
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
			device.ExtinguisherTypeName.Valid = false
		}
		if !device.SerialNumber.Valid {
			device.SerialNumber.String = "N/A"
			device.SerialNumber.Valid = false
		}
		if !device.Description.Valid {
			device.Description.String = "N/A"
			device.Description.Valid = false
		}
		if !device.Size.Valid {
			device.Size.String = "N/A"
			device.Size.Valid = false
		}
		if !device.Status.Valid {
			device.Status.String = "N/A"
			device.Status.Valid = false
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

func (db *DB) GetAllDeviceTypes() ([]models.EmergencyDeviceType, error) {
	query := `
	SELECT emergencydevicetypeid, emergencydevicetypename
	FROM emergency_device_typeT
	ORDER BY emergencydevicetypename
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deviceTypes []models.EmergencyDeviceType

	// Scan the results
	for rows.Next() {
		var deviceType models.EmergencyDeviceType
		err := rows.Scan(
			&deviceType.EmergencyDeviceTypeID,
			&deviceType.EmergencyDeviceTypeName,
		)
		if err != nil {
			return nil, err
		}

		deviceTypes = append(deviceTypes, deviceType)
	}

	return deviceTypes, nil
}

func (db *DB) GetAllExtinguisherTypes() ([]models.ExtinguisherType, error) {
	query := `
	SELECT extinguishertypeid, extinguishertypename
	FROM Extinguisher_TypeT
	ORDER BY extinguishertypename
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var extinguisherTypes []models.ExtinguisherType

	// Scan the results
	for rows.Next() {
		var extinguisherType models.ExtinguisherType
		err := rows.Scan(
			&extinguisherType.ExtinguisherTypeID,
			&extinguisherType.ExtinguisherTypeName,
		)
		if err != nil {
			return nil, err
		}

		extinguisherTypes = append(extinguisherTypes, extinguisherType)
	}

	return extinguisherTypes, nil
}

func (db *DB) GetAllRooms(buildingId string) ([]models.Room, error) {
	var query string
	var args []interface{}

	// Define the base query
	query = ` SELECT r.roomid, r.buildingid, r.roomcode
			  FROM roomT r`

	// Add filtering by building code if provided
	if buildingId != "" {
		query += `
		JOIN buildingT b ON r.buildingid = b.buildingid
		WHERE b.buildingcode = $1
		`
		args = append(args, buildingId)
	}

	fmt.Println(query)

	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var rooms []models.Room

	// Scan the results
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.RoomID,
			&room.BuildingID,
			&room.RoomCode,
		)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (db *DB) GetAllBuildings() ([]models.Building, error) {
	query := `
	SELECT buildingid, buildingcode
	FROM buildingT
	ORDER BY buildingcode
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []models.Building

	// Scan the results
	for rows.Next() {
		var building models.Building
		err := rows.Scan(
			&building.BuildingID,
			&building.BuildingCode,
		)
		if err != nil {
			return nil, err
		}

		buildings = append(buildings, building)
	}

	return buildings, nil
}
