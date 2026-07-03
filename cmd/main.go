package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/app"
	"github.com/fabianherreracruz/bre-b-pse-app/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run the application
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Starting BRE-B PSE Recaudos App on port %s\n", port)
	if err := application.Run(port); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
