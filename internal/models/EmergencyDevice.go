// models/emergency_device.go
package models

import (
	"database/sql"
	"time"
)

type EmergencyDevice struct {
    EmergencyDeviceID        int       `json:"emergency_device_id"` // From emergency_deviceT table
    EmergencyDeviceTypeName  string    `json:"emergency_device_type_name"` // From emergency_device_typeT table
    ExtinguisherType         string    `json:"extinguisher_type"` // From emergency_deviceT table
    RoomName                 string    `json:"room_name"` // From roomT table
    SerialNumber             string    `json:"serial_number"` // From emergency_deviceT table
    ManufactureDate          time.Time `json:"manufacture_date"` // From emergency_deviceT table
    LastInspectionDate       sql.NullTime `json:"last_inspection_date"` // From emergency_deviceT table
    Description              string    `json:"description"` // From emergency_deviceT table
    Size                     string   `json:"size"` // From emergency_deviceT table
    Status                   string    `json:"status"` // From emergency_deviceT table
}