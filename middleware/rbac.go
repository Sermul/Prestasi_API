package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RoleGuard(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		role := c.Locals("role")
		if role == nil {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		roleName := role.(string)

		// Case-insensitive compare
		for _, allowed := range allowedRoles {
			if strings.EqualFold(allowed, roleName) {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
	}
}
