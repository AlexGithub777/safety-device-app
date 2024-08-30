package database

import (
	"database/sql"
	"log"
)

func SeedData(db *sql.DB) {
	// Insert into SiteT
	var siteID int
	err := db.QueryRow(`INSERT INTO SiteT (SiteName, SiteAddress) VALUES 
    ('EIT', '123 Campus Rd') RETURNING SiteID`).Scan(&siteID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert into BuildingT
	var buildingID int
	err = db.QueryRow(`INSERT INTO BuildingT (SiteID, BuildingCode) VALUES 
    ($1, 'A') RETURNING BuildingID`, siteID).Scan(&buildingID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert into RoomT
	var roomID int
	err = db.QueryRow(`INSERT INTO RoomT (BuildingID, Name) VALUES 
    ($1, 'Room 1') RETURNING RoomID`, buildingID).Scan(&roomID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert into Emergency_Device_TypeT
	var emergencyDeviceTypeID int
	err = db.QueryRow(`INSERT INTO Emergency_Device_TypeT (EmergencyDeviceTypeName) VALUES 
    ('Fire Extinguisher') RETURNING EmergencyDeviceTypeID`).Scan(&emergencyDeviceTypeID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert multiple emergency devices into Emergency_DeviceT
	for i := 1; i <= 3; i++ {
		_, err = db.Exec(`INSERT INTO Emergency_DeviceT (EmergencyDeviceTypeID, RoomID, ManufactureDate, SerialNumber, Description, Size, Status) 
        VALUES ($1, $2, '2024-01-01', 'SN0000' || $3, 'Test Fire Extinguisher ' || $3, '5kg', 'Active')`,
			emergencyDeviceTypeID, roomID, i)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Test data seeded successfully!")
}