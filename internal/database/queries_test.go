package database_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/AlexGithub777/safety-device-app/internal/database"
	"github.com/AlexGithub777/safety-device-app/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllDevices(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create SQL mock: %v", err)
	}
	defer db.Close()

	dbInstance := &database.DB{DB: db}

	expectedDevices := []models.EmergencyDevice{
		{
			EmergencyDeviceID:       1,
			EmergencyDeviceTypeName: "TypeA",
			ExtinguisherTypeName:    sql.NullString{String: "ExtinguisherA", Valid: true},
			RoomName:                "Room101",
			SerialNumber:            sql.NullString{String: "SN0001", Valid: true},
			ManufactureDate:         sql.NullTime{Time: time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			LastInspectionDate:      sql.NullTime{Time: time.Date(2024, time.July, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			Description:             sql.NullString{String: "Description of device", Valid: true},
			Size:                    sql.NullString{String: "10kg", Valid: true},
			Status:                  sql.NullString{String: "Active", Valid: true},
		},
	}

	rows := sqlmock.NewRows([]string{
		"emergencydeviceid",
		"emergencydevicetypename",
		"extinguishertypename",
		"roomname",
		"serialnumber",
		"manufacturedate",
		"lastinspectiondate",
		"description",
		"size",
		"status",
	})

	for _, device := range expectedDevices {
		rows.AddRow(
			device.EmergencyDeviceID,
			device.EmergencyDeviceTypeName,
			device.ExtinguisherTypeName.String,
			device.RoomName,
			device.SerialNumber.String,
			device.ManufactureDate.Time,
			device.LastInspectionDate.Time,
			device.Description.String,
			device.Size.String,
			device.Status.String,
		)
	}
	mock.ExpectQuery("^SELECT (.+) FROM emergency_deviceT").WillReturnRows(rows)

	actualDevices, err := dbInstance.FetchAllDevices("your_building_code")
	if err != nil {
		t.Fatalf("Failed to fetch devices: %v", err)
	}

	// Manually calculate the two additional fields for the expected devices
	for i := range expectedDevices {
		expectedDevices[i].ExpireDate = sql.NullTime{Time: expectedDevices[i].ManufactureDate.Time.AddDate(5, 0, 0), Valid: true}
		expectedDevices[i].NextInspectionDate = sql.NullTime{Time: expectedDevices[i].LastInspectionDate.Time.AddDate(0, 3, 0), Valid: true}
	}

	assert.ElementsMatch(t, expectedDevices, actualDevices)
}
