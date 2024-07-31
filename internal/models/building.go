package models

type Building struct {
	BuildingID int    `json:"building_id"`
	SiteID     int    `json:"site_id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
}
