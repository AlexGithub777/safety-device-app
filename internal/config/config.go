package config

import (
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> d3f1aef86552b4414e16af4df61e0e15859fe0b5
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     int
}

func LoadConfig() Config {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the DB_PORT environment variable
	dbPortStr := os.Getenv("DB_PORT")

<<<<<<< HEAD
=======
	fmt.Println(dbPortStr)

>>>>>>> d3f1aef86552b4414e16af4df61e0e15859fe0b5
	// Convert the string to an integer and handle any potential errors
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	return Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
	}
}
