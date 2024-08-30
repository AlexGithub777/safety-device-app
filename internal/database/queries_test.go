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

// TestFetchAllDevices with valid non-null fields
func TestFetchAllDevices(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create SQL mock: %v", err)
	}
	defer db.Close()

	dbInstance := &database.DB{DB: db}

	expectedDevices1 := []models.EmergencyDevice{
		{
			EmergencyDeviceID:       1,
			EmergencyDeviceTypeName: "TypeA",
			ExtinguisherTypeName:    sql.NullString{String: "ExtinguisherA", Valid: true},
			RoomName:                "Room101",
			SerialNumber:            sql.NullString{String: "SN123", Valid: true},
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

	for _, device := range expectedDevices1 {
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

	// Assert Expected NextInspectionDate and ExpireDate for the first device (all fields are not null)
	expectedDevices1[0].NextInspectionDate = sql.NullTime{Time: time.Date(2024, time.October, 15, 0, 0, 0, 0, time.UTC), Valid: true}
	expectedDevices1[0].ExpireDate = sql.NullTime{Time: time.Date(2029, time.August, 1, 0, 0, 0, 0, time.UTC), Valid: true}

	actualDevices, err := dbInstance.FetchAllDevices("your_building_code")
	if err != nil {
		t.Fatalf("Failed to fetch devices: %v", err)
	}

	assert.ElementsMatch(t, expectedDevices1, actualDevices)

	// Additional assertions for NextInspectionDate and ExpireDate
	for i, actualDevice := range actualDevices {
		assert.Equal(t, expectedDevices1[i].NextInspectionDate.Time, actualDevice.NextInspectionDate.Time, "NextInspectionDate does not match")
		assert.Equal(t, expectedDevices1[i].ExpireDate.Time, actualDevice.ExpireDate.Time, "ExpireDate does not match")
	}
}

// TestFetchAllDevices with valid null fields
func TestFetchAllDevices2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create SQL mock: %v", err)
	}
	defer db.Close()

	dbInstance := &database.DB{DB: db}

	expectedDevices2 := []models.EmergencyDevice{
		{
			EmergencyDeviceID:       2,
			EmergencyDeviceTypeName: "TypeB",
			ExtinguisherTypeName:    sql.NullString{Valid: true, String: "N/A"},
			RoomName:                "Room102",
			SerialNumber:            sql.NullString{Valid: true, String: "N/A"},
			ManufactureDate:         sql.NullTime{Time: time.Date(0001, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			LastInspectionDate:      sql.NullTime{Time: time.Date(0001, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			Description:             sql.NullString{Valid: true, String: "N/A"},
			Size:                    sql.NullString{Valid: true, String: "N/A"},
			Status:                  sql.NullString{Valid: true, String: "N/A"},
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

	for _, device := range expectedDevices2 {
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

	// Assert Expected NextInspectionDate and ExpireDate for the second device (all fields which are nullable fields are null)
	expectedDevices2[0].NextInspectionDate = sql.NullTime{Time: time.Date(0001, time.April, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	expectedDevices2[0].ExpireDate = sql.NullTime{Time: time.Date(0006, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true}

	actualDevices, err := dbInstance.FetchAllDevices("your_building_code")
	if err != nil {
		t.Fatalf("Failed to fetch devices: %v", err)
	}

	assert.ElementsMatch(t, expectedDevices2, actualDevices)

	// Additional assertions for NextInspectionDate and ExpireDate
	for i, actualDevice := range actualDevices {
		assert.Equal(t, expectedDevices2[i].NextInspectionDate.Time, actualDevice.NextInspectionDate.Time, "NextInspectionDate does not match")
		assert.Equal(t, expectedDevices2[i].ExpireDate.Time, actualDevice.ExpireDate.Time, "ExpireDate does not match")
	}
}
