package models

type FireExtinguisherRecord struct {
	FireExtinguisherRecordID int    `json:"fire_extinguisher_record_id"`
	FireExtinguisherID       int    `json:"fire_extinguisher_id"`
	Date                     string `json:"date"`
	Notes                    string `json:"notes"`
}
