-- +goose Up
CREATE TABLE sites (
    site_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255)
);

CREATE TABLE buildings (
    building_id SERIAL PRIMARY KEY,
    site_id INT REFERENCES sites(site_id),
    name VARCHAR(100) NOT NULL
);

CREATE TABLE rooms (
    room_id SERIAL PRIMARY KEY,
    building_id INT REFERENCES buildings(building_id),
    name VARCHAR(100) NOT NULL
);

CREATE TABLE safety_devices (
    safety_device_id SERIAL PRIMARY KEY,
    safety_device_type VARCHAR(50) NOT NULL,
    room_id INT REFERENCES rooms(room_id),
    status VARCHAR(50)
);

CREATE TABLE fire_extinguisher_types (
    fire_extinguisher_type_id SERIAL PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL
);

CREATE TABLE fire_extinguishers (
    fire_extinguisher_id SERIAL PRIMARY KEY,
    safety_device_id INT REFERENCES safety_devices(safety_device_id) UNIQUE,
    fire_extinguisher_type_id INT REFERENCES fire_extinguisher_types(fire_extinguisher_type_id),
    serial_number VARCHAR(100),
    date_of_manufacture DATE NOT NULL,
    expire_date DATE,
    size VARCHAR(50),
    misc TEXT,
    status VARCHAR(200)
);

CREATE TABLE fire_extinguisher_records (
    fire_extinguisher_record_id SERIAL PRIMARY KEY,
    fire_extinguisher_id INT REFERENCES fire_extinguishers(fire_extinguisher_id),
    date DATE NOT NULL,
    notes TEXT
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    role VARCHAR(50) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS fire_extinguisher_records;
DROP TABLE IF EXISTS fire_extinguishers;
DROP TABLE IF EXISTS fire_extinguisher_types;
DROP TABLE IF EXISTS safety_devices;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS buildings;
DROP TABLE IF EXISTS sites;
DROP TABLE IF EXISTS users;