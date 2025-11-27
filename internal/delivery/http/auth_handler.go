package http

import (
	"github.com/gofiber/fiber/v2"
)

// AuthHandler provee un endpoint para login básico y generación de JWT.
func AuthHandler(c *fiber.Ctx) error {
	type loginRequest struct {
		UserID string `json:"user_id"`
	}
	var req loginRequest
	if err := c.BodyParser(&req); err != nil || req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	token, err := GenerateJWT(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}
	return c.JSON(fiber.Map{"token": token})
}
