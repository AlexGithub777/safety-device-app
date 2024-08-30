package models

// RoomT represents rooms in each building
type Room struct {
	RoomID     int    `json:"room_id"`
	BuildingID int    `json:"building_id"`
	Name       string `json:"name"`
}