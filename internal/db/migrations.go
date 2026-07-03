package db

import (
	"log"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/models"
	"gorm.io/gorm"
)

// MigrateDatabase ejecuta todas las migraciones
func MigrateDatabase(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Crear tablas
	if err := db.AutoMigrate(
		&models.User{},
		&models.Recaudo{},
		&models.Split{},
		&models.Notificacion{},
		&models.AuditLog{},
	); err != nil {
		log.Printf("Error running migrations: %v", err)
		return err
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// SeedDatabase inserta datos de prueba
func SeedDatabase(db *gorm.DB) error {
	log.Println("Seeding database...")

	// Verificar si ya existen datos
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	return nil
}
