package models

type SafetyDevice struct {
	SafetyDeviceID   int    `json:"safety_device_id"`
	SafetyDeviceType string `json:"safety_device_type"`
	RoomID           int    `json:"room_id"`
	Status           string `json:"status"`
}
