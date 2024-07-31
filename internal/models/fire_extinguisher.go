package models

import (
	"time"
)

type FireExtinguisher struct {
	FireExtinguisherID     int       `json:"fire_extinguisher_id"`
	SafetyDeviceID         int       `json:"safety_device_id"`
	FireExtinguisherTypeID int       `json:"fire_extinguisher_type_id"`
	SerialNumber           string    `json:"serial_number"`
	DateOfManufacture      time.Time `json:"date_of_manufacture"`
	ExpireDate             time.Time `json:"expire_date"`
	Size                   string    `json:"size"`
	Misc                   string    `json:"misc"`
	Status                 string    `json:"status"`
}
