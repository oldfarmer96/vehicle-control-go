package middlewares

import (
	"github.com/gofiber/fiber/v3"
)

func ApiKeyMiddleware(apiKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		provided := c.Get("x-api-key")
		if provided == "" || provided != apiKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key inválida o missing",
			})
		}
		return c.Next()
	}
}