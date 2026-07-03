package db

import (
	"fmt"
	"log"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("✅ Database connection established")
	return db, nil
}

func MigrateDatabase(db *gorm.DB) error {
	// Las migraciones se ejecutarán aquí
	return nil
}
