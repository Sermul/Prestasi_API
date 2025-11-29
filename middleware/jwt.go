package middleware

import (
	"strings"

	"prestasi_api/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Missing or invalid token"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := helper.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])  // ⬅️ simpan role_name

		return c.Next()
	}
}
