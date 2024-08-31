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
	testCases := []struct {
		name            string
		expectedDevices []models.EmergencyDevice
		nextInspection  time.Time
		expireDate      time.Time
	}{
		{
			name: "TestFetchAllDevices with valid non-null fields",
			expectedDevices: []models.EmergencyDevice{
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
			},
			nextInspection: time.Date(2024, time.October, 15, 0, 0, 0, 0, time.UTC),
			expireDate:     time.Date(2029, time.August, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "TestFetchAllDevices with valid null fields",
			expectedDevices: []models.EmergencyDevice{
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
			},
			nextInspection: time.Date(0001, time.April, 1, 0, 0, 0, 0, time.UTC),
			expireDate:     time.Date(0006, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "TestFetchAllDevices with invalid null fields",
			expectedDevices: []models.EmergencyDevice{
				{
					EmergencyDeviceID:       3,
					EmergencyDeviceTypeName: "TypeC",
					ExtinguisherTypeName:    sql.NullString{Valid: false, String: "N/A"}, // Null extinguisher type name
					RoomName:                "Room103",
					SerialNumber:            sql.NullString{Valid: false, String: "N/A"}, // Null serial number
					ManufactureDate:         sql.NullTime{Time: time.Date(0001, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
					LastInspectionDate:      sql.NullTime{Time: time.Date(0001, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
					Description:             sql.NullString{Valid: false, String: "N/A"}, // Null description
					Size:                    sql.NullString{Valid: false, String: "N/A"}, // Null size
					Status:                  sql.NullString{Valid: false, String: "N/A"}, // Null status
				},
			},
			nextInspection: time.Date(0001, time.April, 1, 0, 0, 0, 0, time.UTC),
			expireDate:     time.Date(0006, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create SQL mock: %v", err)
			}
			defer db.Close()

			dbInstance := &database.DB{DB: db}

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

			for _, device := range tc.expectedDevices {
				if i == 2 { // If it's the third test case (index starts from 0)
					rows.AddRow(
						device.EmergencyDeviceID,
						device.EmergencyDeviceTypeName,
						sql.NullString{Valid: false}, // Mock the database to return Valid: false
						device.RoomName,
						sql.NullString{Valid: false}, // Mock the database to return Valid: false
						device.ManufactureDate.Time,
						device.LastInspectionDate.Time,
						sql.NullString{Valid: false}, // Mock the database to return Valid: false
						sql.NullString{Valid: false}, // Mock the database to return Valid: false
						sql.NullString{Valid: false}, // Mock the database to return Valid: false
					)
				} else {
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
			}

			mock.ExpectQuery("^SELECT (.+) FROM emergency_deviceT").WillReturnRows(rows)

			// Assert Expected NextInspectionDate and ExpireDate for the first device (all fields are not null)
			tc.expectedDevices[0].NextInspectionDate = sql.NullTime{Time: tc.nextInspection, Valid: true}
			tc.expectedDevices[0].ExpireDate = sql.NullTime{Time: tc.expireDate, Valid: true}

			actualDevices, err := dbInstance.FetchAllDevices("your_building_code")
			if err != nil {
				t.Fatalf("Failed to fetch devices: %v", err)
			}

			// Assert the expected devices and actual devices are equal

			assert.ElementsMatch(t, tc.expectedDevices, actualDevices)

			// Additional assertions for NextInspectionDate and ExpireDate
			for i, actualDevice := range actualDevices {
				assert.Equal(t, tc.expectedDevices[i].NextInspectionDate.Time, actualDevice.NextInspectionDate.Time, "NextInspectionDate does not match")
				assert.Equal(t, tc.expectedDevices[i].ExpireDate.Time, actualDevice.ExpireDate.Time, "ExpireDate does not match")
			}
		})
	}
}
