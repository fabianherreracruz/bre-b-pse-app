package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// ePayco
	EPaycoClientID     string
	EPaycoClientSecret string
	EPaycoAPIKey       string
	EPaycoPrivateKey   string

	// JWT
	JWTSecret      string
	JWTExpiration  time.Duration

	// Twilio
	TwilioAccountSID    string
	TwilioAuthToken     string
	TwilioPhoneNumber   string
	TwilioWhatsAppNumber string

	// SendGrid
	SendGridAPIKey  string
	SendGridFromEmail string

	// App
	AppEnv  string
	AppPort string
	AppURL  string

	// Webhook
	WebhookSecret string
}

func LoadConfig() *Config {
	jwtExpiration := 86400 // 24 hours
	if exp := os.Getenv("JWT_EXPIRATION"); exp != "" {
		if parsedExp, err := strconv.Atoi(exp); err == nil {
			jwtExpiration = parsedExp
		}
	}

	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "bre_b_pse_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// ePayco
		EPaycoClientID:     getEnv("EPAYCO_CLIENT_ID", ""),
		EPaycoClientSecret: getEnv("EPAYCO_CLIENT_SECRET", ""),
		EPaycoAPIKey:       getEnv("EPAYCO_API_KEY", ""),
		EPaycoPrivateKey:   getEnv("EPAYCO_PRIVATE_KEY", ""),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiration: time.Duration(jwtExpiration) * time.Second,

		// Twilio
		TwilioAccountSID:     getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:      getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhoneNumber:    getEnv("TWILIO_PHONE_NUMBER", ""),
		TwilioWhatsAppNumber: getEnv("TWILIO_WHATSAPP_NUMBER", ""),

		// SendGrid
		SendGridAPIKey:    getEnv("SENDGRID_API_KEY", ""),
		SendGridFromEmail: getEnv("SENDGRID_FROM_EMAIL", ""),

		// App
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppURL:  getEnv("APP_URL", "http://localhost:8080"),

		// Webhook
		WebhookSecret: getEnv("WEBHOOK_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
