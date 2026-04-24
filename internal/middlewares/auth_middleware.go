// Package middlewares  protector de autenticacion
package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Extraer el token de la cookie
		cookieName := os.Getenv("COOKIE_NAME")
		tokenString := c.Cookies(cookieName)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No autorizado, token ausente",
			})
		}

		// 2. Parsear y validar el token
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token inválido o expirado",
			})
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Locals("userID", claims["sub"])
		c.Locals("userRole", claims["role"])
		c.Locals("userEmail", claims["email"])

		return c.Next()
	}
}

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole := c.Locals("userRole")
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "No tienes permisos suficientes",
		})
	}
}
