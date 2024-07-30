package database

import (
	"database/sql"
	"fmt"

	"github.com/AlexGithub777/safety-device-app/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type DB struct {
	*sql.DB
}

func NewDB(cfg config.Config) (*DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
