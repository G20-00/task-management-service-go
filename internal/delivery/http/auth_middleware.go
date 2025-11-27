package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware verifica el JWT en el header Authorization.
func JWTMiddleware(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	token, err := ParseJWT(tokenStr)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}
	_, ok := GetUserIDFromToken(token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}
	return c.Next()
}
