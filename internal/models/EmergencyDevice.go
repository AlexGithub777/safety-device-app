package models

import "time"

// Emergency_DeviceT represents the emergency devices
type EmergencyDevice struct {
	EmergencyDeviceID     int       `json:"emergency_device_id"`
	EmergencyDeviceTypeID int       `json:"emergency_device_type_id"`
	RoomID                int       `json:"room_id"`
	ManufactureDate       time.Time `json:"manufacture_date"`
	SerialNumber          string    `json:"serial_number"`
	Description           string    `json:"description"`
	Size                  string    `json:"size"`
	LastInspectionDate    *time.Time `json:"last_inspection_date,omitempty"`
	Status                string    `json:"status"`
}