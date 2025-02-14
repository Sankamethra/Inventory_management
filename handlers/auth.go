package handlers

import (
	"order-inventory/config"
	"order-inventory/models"
	"os"
	"time"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// Custom claims struct
type Claims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func Login(c *fiber.Ctx) error {
	var request LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}

	var user models.User
	result := config.DB.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
		})
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
		})
	}

	// Create claims with custom struct
	claims := Claims{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "order-inventory-system",
			Subject:   user.Email,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not generate token",
			"error":   err.Error(),
		})
	}

	// Return token and user info
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": signedToken,
			"user": fiber.Map{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"role":       user.Role,
			},
		},
	})
}

func Signup(c *fiber.Ctx) error {
	var request SignupRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}

	// Check if user already exists
	var existingUser models.User
	if result := config.DB.Where("email = ?", request.Email).First(&existingUser); result.Error == nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not hash password",
			"error":   err.Error(),
		})
	}

	// Create user
	user := models.User{
		Email:     request.Email,
		Password:  string(hashedPassword),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Role:      "user", // Default role
	}

	if result := config.DB.Create(&user); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user",
			"error":   result.Error.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user": fiber.Map{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"role":       user.Role,
			},
		},
	})
}

// SetupAdmin creates admin accounts (requires setup secret)
func SetupAdmin(c *fiber.Ctx) error {
	// Debug: Print environment variables
	fmt.Printf("SETUP_SECRET from env: %s\n", os.Getenv("SETUP_SECRET"))
	fmt.Printf("Setup-Secret from request: %s\n", c.Get("Setup-Secret"))

	// Verify setup secret
	setupSecret := c.Get("Setup-Secret")
	if setupSecret != os.Getenv("SETUP_SECRET") {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid setup secret",
		})
	}

	var request SignupRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}

	// Check if email is already registered
	var existingUser models.User
	if result := config.DB.Where("email = ?", request.Email).First(&existingUser); result.Error == nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email already registered",
		})
	}

	// Create admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not hash password",
			"error":   err.Error(),
		})
	}

	user := models.User{
		Email:     request.Email,
		Password:  string(hashedPassword),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Role:      "admin",
	}

	if result := config.DB.Create(&user); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create admin user",
			"error":   result.Error.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Admin account created successfully",
		"data": fiber.Map{
			"user": fiber.Map{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"role":       user.Role,
			},
		},
	})
}
