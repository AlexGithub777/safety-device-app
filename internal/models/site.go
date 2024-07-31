package models

type Site struct {
	SiteID  int    `json:"site_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
