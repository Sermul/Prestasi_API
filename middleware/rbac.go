package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RoleGuard(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Ambil role langsung dari JWT (bukan DB)
		role := c.Locals("role")
		if role == nil {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Compare case-insensitive
		for _, allowed := range allowedRoles {
			if strings.EqualFold(allowed, role.(string)) {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
	}
}
