-- +goose Up
CREATE TABLE UserT (
    UserID SERIAL PRIMARY KEY,
    Username VARCHAR(100) NOT NULL,
    Password VARCHAR(255) NOT NULL
);

CREATE TABLE SiteT (
    SiteID SERIAL PRIMARY KEY,
    SiteName VARCHAR(100) NOT NULL,
    SiteAddress VARCHAR(255)
);

CREATE TABLE BuildingT (
    BuildingID SERIAL PRIMARY KEY,
    SiteID INT REFERENCES SiteT(SiteID),
    BuildingCode VARCHAR(100) NOT NULL
);

CREATE TABLE RoomT (
    RoomID SERIAL PRIMARY KEY,
    BuildingID INT REFERENCES BuildingT(BuildingID),
    Name VARCHAR(100) NOT NULL
);

CREATE TABLE Emergency_Device_TypeT (
    EmergencyDeviceTypeID SERIAL PRIMARY KEY,
    EmergencyDeviceTypeName VARCHAR(50) NOT NULL
);

CREATE TABLE Emergency_DeviceT (
    EmergencyDeviceID SERIAL PRIMARY KEY,
    EmergencyDeviceTypeID  INT REFERENCES Emergency_Device_TypeT(EmergencyDeviceTypeID),
    RoomID INT REFERENCES RoomT(RoomID),
    ManufactureDate DATE NOT NULL,
    SerialNumber VARCHAR(100),
    Description VARCHAR(80) NULL,
    Size VARCHAR(50),
    LastInspectionDate DATE NULL,
    Status VARCHAR(50) NULL
);
 
CREATE TABLE Emergency_Device_InspectionT (
    EmergencyDeviceInspectionID SERIAL PRIMARY KEY,
    EmergencyDeviceID INT REFERENCES Emergency_DeviceT(EmergencyDeviceID),
    UserID INT REFERENCES UserT(UserID),
    InspectionDate DATE NOT NULL,
    Notes VARCHAR(255) NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Added missing comma
    IsConspicuous BOOLEAN NULL,
    IsAccessible BOOLEAN NULL,
    IsAssignedLocation BOOLEAN NULL,
    IsSignVisible BOOLEAN NULL,
    IsAntiTamperDeviceIntact BOOLEAN NULL,
    IsSupportBracketSecure BOOLEAN NULL,
    AreOperatingInstructionsClear BOOLEAN NULL,
    IsMaintenanceTagAttached BOOLEAN NULL,
    IsExternalDamagePresent BOOLEAN NULL,
    IsChargeGaugeNormal BOOLEAN NULL,
    IsReplaced BOOLEAN NULL,
    AreMaintenanceRecordsComplete BOOLEAN NULL,
    WorkOrderRequired BOOLEAN NULL,
    IsInspectionComplete BOOLEAN NULL
);

-- +goose Down
DROP TABLE IF EXISTS Emergency_Device_InspectionT;
DROP TABLE IF EXISTS Emergency_DeviceT;
DROP TABLE IF EXISTS Emergency_Device_TypeT;
DROP TABLE IF EXISTS RoomT;
DROP TABLE IF EXISTS BuildingT;
DROP TABLE IF EXISTS SiteT;
DROP TABLE IF EXISTS UserT;