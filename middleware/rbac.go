package middleware

import (
	"strings"

	"prestasi_api/app/repository"

	"github.com/gofiber/fiber/v2"
)

// RoleGuard menggunakan role_id dari token lalu cek ke database
func RoleGuard(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Ambil role_id dari JWTMiddleware
		roleID := c.Locals("role_id")
		if roleID == nil {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Query role name dari database
		roleRepo := repository.NewRolePostgresRepository()
		role, err := roleRepo.GetByID(roleID.(string))
		if err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "Invalid role"})
		}

		roleName := role.Name

		// Case-insensitive compare
		for _, allowed := range allowedRoles {
			if strings.EqualFold(allowed, roleName) {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
	}
}
