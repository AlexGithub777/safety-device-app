package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/models"
)

func seedData(db *sql.DB) {
	var siteID, buildingIDA, buildingIDB int
	var roomA1ID, roomB1ID int
	var co2TypeID, waterTypeID, dryTypeID int
	var emergencyDeviceTypeID int

	// Insert Site
	err := db.QueryRow(`
		INSERT INTO SiteT (SiteName, SiteAddress) 
		VALUES ('EIT', '501 Gloucester Street, Taradale, Napier 4112') RETURNING SiteID`).Scan(&siteID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert Buildings A and B
	err = db.QueryRow(`
		INSERT INTO BuildingT (SiteID, BuildingCode) 
		VALUES ($1, 'A') RETURNING BuildingID`, siteID).Scan(&buildingIDA)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow(`
		INSERT INTO BuildingT (SiteID, BuildingCode) 
		VALUES ($1, 'B') RETURNING BuildingID`, siteID).Scan(&buildingIDB)
	if err != nil {
		log.Fatal(err)
	}

	// Insert Rooms A1 and B1
	err = db.QueryRow(`
		INSERT INTO RoomT (BuildingID, RoomCode) 
		VALUES ($1, 'A1') RETURNING RoomID`, buildingIDA).Scan(&roomA1ID)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow(`
		INSERT INTO RoomT (BuildingID, RoomCode) 
		VALUES ($1, 'B1') RETURNING RoomID`, buildingIDB).Scan(&roomB1ID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert Extinguisher Types
	err = db.QueryRow(`
		INSERT INTO Extinguisher_TypeT (ExtinguisherTypeName) 
		VALUES ('CO2') RETURNING ExtinguisherTypeID`).Scan(&co2TypeID)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow(`
		INSERT INTO Extinguisher_TypeT (ExtinguisherTypeName) 
		VALUES ('Water') RETURNING ExtinguisherTypeID`).Scan(&waterTypeID)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow(`
		INSERT INTO Extinguisher_TypeT (ExtinguisherTypeName) 
		VALUES ('Dry') RETURNING ExtinguisherTypeID`).Scan(&dryTypeID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert Emergency Device Type
	err = db.QueryRow(`
		INSERT INTO Emergency_Device_TypeT (EmergencyDeviceTypeName) 
		VALUES ('Fire Extinguisher') RETURNING EmergencyDeviceTypeID`).Scan(&emergencyDeviceTypeID)
	if err != nil {
		log.Fatal(err)
	}

	// Create Emergency Devices using the models.EmergencyDevice struct
	devices := []models.EmergencyDevice{
		{
			EmergencyDeviceTypeName: "Fire Extinguisher",
			ExtinguisherTypeName:    sql.NullString{Valid: true, String: "CO2"},
			RoomCode:                "A1",
			SerialNumber:            sql.NullString{Valid: true, String: "SN00001"},
			ManufactureDate:         sql.NullTime{Valid: true, Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			LastInspectionDate:      sql.NullTime{Valid: true, Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			Description:             sql.NullString{Valid: true, String: "Test Fire Extinguisher 1"},
			Size:                    sql.NullString{Valid: true, String: "5kg"},
			Status:                  sql.NullString{Valid: true, String: "Active"},
		},
		{
			EmergencyDeviceTypeName: "Fire Extinguisher",
			ExtinguisherTypeName:    sql.NullString{Valid: true, String: "Water"},
			RoomCode:                "B1",
			SerialNumber:            sql.NullString{Valid: true, String: "SN00002"},
			ManufactureDate:         sql.NullTime{Valid: true, Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			LastInspectionDate:      sql.NullTime{Valid: false},
			Description:             sql.NullString{Valid: true, String: "Test Fire Extinguisher 2"},
			Size:                    sql.NullString{Valid: true, String: "5kg"},
			Status:                  sql.NullString{Valid: true, String: "Expired"},
		},
		{
			EmergencyDeviceTypeName: "Fire Extinguisher",
			ExtinguisherTypeName:    sql.NullString{Valid: true, String: "Dry"},
			RoomCode:                "A1",
			SerialNumber:            sql.NullString{Valid: true, String: "SN00003"},
			ManufactureDate:         sql.NullTime{Valid: true, Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			LastInspectionDate:      sql.NullTime{Valid: false},
			Description:             sql.NullString{Valid: true, String: "Test Fire Extinguisher 3"},
			Size:                    sql.NullString{Valid: true, String: "5kg"},
			Status:                  sql.NullString{Valid: true, String: "Inactive"},
		},
	}

	// Insert Emergency Devices into the database
	for _, device := range devices {
		var roomID int
		var extinguisherTypeID int

		// Map RoomCode to RoomID
		if device.RoomCode == "A1" {
			roomID = roomA1ID
		} else {
			roomID = roomB1ID
		}

		// Map ExtinguisherTypeName to ExtinguisherTypeID
		switch device.ExtinguisherTypeName.String {
		case "CO2":
			extinguisherTypeID = co2TypeID
		case "Water":
			extinguisherTypeID = waterTypeID
		case "Dry":
			extinguisherTypeID = dryTypeID
		}

		_, err := db.Exec(`
			INSERT INTO Emergency_DeviceT 
				(EmergencyDeviceTypeID, RoomID, ExtinguisherTypeID, SerialNumber, ManufactureDate, LastInspectionDate, Description, Size, Status) 
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			emergencyDeviceTypeID, roomID, extinguisherTypeID,
			device.SerialNumber, device.ManufactureDate, device.LastInspectionDate, device.Description, device.Size, device.Status,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding complete.")
}
