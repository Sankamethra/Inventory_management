package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "Authorization header is required",
			})
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid authorization format, expected 'Bearer {token}'",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid or expired token",
				"error":   err.Error(),
			})
		}

		// Validate token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid token claims",
			})
		}

		// Check token expiration
		if exp, ok := claims["exp"].(float64); !ok || float64(time.Now().Unix()) > exp {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "Token has expired",
			})
		}

		// Validate required claims
		if _, ok := claims["user_id"].(float64); !ok {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid token: missing user_id",
			})
		}

		if _, ok := claims["role"].(string); !ok {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid token: missing role",
			})
		}

		// Set user in context
		c.Locals("user", claims)
		return c.Next()
	}
}
