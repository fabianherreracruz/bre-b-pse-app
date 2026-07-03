package services

import (
	"fmt"
	"log"
	"time"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/config"
	"github.com/fabianherreracruz/bre-b-pse-app/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required,oneof=recaudador beneficiario"`
}

type AuthResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Token     string    `json:"token"`
	ExpiresAt int64    `json:"expires_at"`
}

type TokenClaims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// Register crea un nuevo usuario
func (a *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Validar que el email no exista
	var existingUser models.User
	if err := a.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, fmt.Errorf("failed to process password")
	}

	// Crear usuario
	user := models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.db.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, fmt.Errorf("failed to create user")
	}

	// Generar token
	token, expiresAt, err := a.generateToken(user)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ User registered: %s", user.Email)

	return &AuthResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// Login autentica un usuario
func (a *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Buscar usuario
	var user models.User
	if err := a.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		log.Printf("User not found: %s", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Validar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Invalid password for user: %s", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Validar estado del usuario
	if user.Status != "active" {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Generar token
	token, expiresAt, err := a.generateToken(user)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ User logged in: %s", user.Email)

	return &AuthResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// generateToken genera un JWT
func (a *AuthService) generateToken(user models.User) (string, int64, error) {
	expiresAt := time.Now().Add(a.cfg.JWTExpiration).Unix()

	claims := TokenClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", 0, fmt.Errorf("failed to generate token")
	}

	return signedToken, expiresAt, nil
}

// ValidateToken valida un JWT
func (a *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetUserByID obtiene un usuario por ID
func (a *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := a.db.Where("id = ?", userID).First(&user).Error; err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

// UpdateUser actualiza datos del usuario
func (a *AuthService) UpdateUser(userID uuid.UUID, firstName, lastName, phone string) (*models.User, error) {
	var user models.User
	if err := a.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Phone = phone
	user.UpdatedAt = time.Now()

	if err := a.db.Save(&user).Error; err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, fmt.Errorf("failed to update user")
	}

	return &user, nil
}

// ChangePassword cambia la contraseña del usuario
func (a *AuthService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	var user models.User
	if err := a.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Validar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Hash nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to process password")
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := a.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password")
	}

	log.Printf("✅ Password changed for user: %s", user.Email)
	return nil
}
