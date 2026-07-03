package services

import (
	"testing"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/config"
	"github.com/fabianherreracruz/bre-b-pse-app/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto migrate
	db.AutoMigrate(&models.User{})
	return db
}

func TestRegisterUser(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: 86400,
	}
	authService := NewAuthService(db, cfg)

	req := RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "+573001234567",
		Password:  "password123",
		Role:      "recaudador",
	}

	resp, err := authService.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "john@example.com", resp.Email)
	assert.Equal(t, "John", resp.FirstName)
	assert.NotEmpty(t, resp.Token)
}

func TestRegisterDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: 86400,
	}
	authService := NewAuthService(db, cfg)

	req := RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "+573001234567",
		Password:  "password123",
		Role:      "recaudador",
	}

	// Primer registro
	_, err := authService.Register(req)
	assert.NoError(t, err)

	// Intento duplicado
	resp, err := authService.Register(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "email already registered", err.Error())
}

func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: 86400,
	}
	authService := NewAuthService(db, cfg)

	// Registrar usuario
	registerReq := RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "+573001234567",
		Password:  "password123",
		Role:      "recaudador",
	}
	_, err := authService.Register(registerReq)
	assert.NoError(t, err)

	// Login
	loginReq := LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	resp, err := authService.Login(loginReq)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "john@example.com", resp.Email)
	assert.NotEmpty(t, resp.Token)
}

func TestLoginInvalidPassword(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: 86400,
	}
	authService := NewAuthService(db, cfg)

	// Registrar usuario
	registerReq := RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "+573001234567",
		Password:  "password123",
		Role:      "recaudador",
	}
	_, err := authService.Register(registerReq)
	assert.NoError(t, err)

	// Login con contraseña incorrecta
	loginReq := LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}
	resp, err := authService.Login(loginReq)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestValidateToken(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: 86400,
	}
	authService := NewAuthService(db, cfg)

	// Registrar y obtener token
	registerReq := RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "+573001234567",
		Password:  "password123",
		Role:      "recaudador",
	}
	resp, err := authService.Register(registerReq)
	assert.NoError(t, err)

	// Validar token
	claims, err := authService.ValidateToken(resp.Token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "john@example.com", claims.Email)
}

func TestChangePassword(t *testing.T) {
	db := setupTestDB(t)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: 86400,
	}
	authService := NewAuthService(db, cfg)

	// Registrar usuario
	registerReq := RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "+573001234567",
		Password:  "password123",
		Role:      "recaudador",
	}
	resp, err := authService.Register(registerReq)
	assert.NoError(t, err)

	// Cambiar contraseña
	err = authService.ChangePassword(resp.ID, "password123", "newpassword123")
	assert.NoError(t, err)

	// Intentar login con nueva contraseña
	loginReq := LoginRequest{
		Email:    "john@example.com",
		Password: "newpassword123",
	}
	loginResp, err := authService.Login(loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, loginResp)
}
