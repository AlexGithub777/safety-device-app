package database

import (
	"database/sql"
	"log"
)

func SeedData(db *sql.DB) {
	// Insert into sites
	var siteID int
	err := db.QueryRow(`INSERT INTO sites (name, address) VALUES 
    ('Test Site', '123 Test Address') RETURNING site_id`).Scan(&siteID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert into buildings
	var buildingID int
	err = db.QueryRow(`INSERT INTO buildings (site_id, name) VALUES 
    ($1, 'Test Building') RETURNING building_id`, siteID).Scan(&buildingID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert into rooms
	var roomID int
	err = db.QueryRow(`INSERT INTO rooms (building_id, name) VALUES 
    ($1, 'Test Room') RETURNING room_id`, buildingID).Scan(&roomID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert into safety_devices
	var safetyDeviceID int
	err = db.QueryRow(`INSERT INTO safety_devices (safety_device_type, room_id, status) VALUES 
    ('Fire Extinguisher', $1, 'Active') RETURNING safety_device_id`, roomID).Scan(&safetyDeviceID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert into fire_extinguisher_types
	var fireExtinguisherTypeID int
	err = db.QueryRow(`INSERT INTO fire_extinguisher_types (type_name) VALUES 
    ('CO2') RETURNING fire_extinguisher_type_id`).Scan(&fireExtinguisherTypeID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert into fire_extinguishers
	_, err = db.Exec(`INSERT INTO fire_extinguishers (safety_device_id, fire_extinguisher_type_id, serial_number, date_of_manufacture, expire_date, size, misc, status)
    VALUES ($1, $2, 'B78866283', '2024-01-01', '2029-01-01', '5kg', 'Test Fire Extinguisher', 'Active')`,
		safetyDeviceID, fireExtinguisherTypeID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Test data seeded successfully!")
}
