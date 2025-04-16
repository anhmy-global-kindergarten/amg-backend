package database

import (
	"amg-backend/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// PostgreSQL DSN (Data Source Name) format: host=host port=port user=user password=password dbname=dbname sslmode=disable TimeZone=UTC
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Database connected successfully.")
	return db, nil
}
