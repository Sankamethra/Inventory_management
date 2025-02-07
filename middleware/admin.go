package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(jwt.MapClaims)
		role := user["role"].(string)

		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"message": "Admin access required",
			})
		}

		return c.Next()
	}
}