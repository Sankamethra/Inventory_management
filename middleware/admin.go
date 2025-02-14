package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by AuthMiddleware)
		user, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "User context not found",
			})
		}

		// Check if role exists
		role, ok := user["role"].(string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "User role not found",
			})
		}

		// Verify admin role
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "Admin access required",
			})
		}

		return c.Next()
	}
}
