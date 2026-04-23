package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Extraer el token de la cookie
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No autorizado, token ausente",
			})
		}

		// 2. Parsear y validar el token
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token inválido o expirado",
			})
		}

		// 3. Extraer Claims y guardarlos en los "Locals" (Contexto de la petición)
		claims := token.Claims.(jwt.MapClaims)

		// Guardamos los datos para usarlos en el controlador
		c.Locals("userID", claims["sub"])
		c.Locals("userRole", claims["role"])
		c.Locals("userEmail", claims["email"])

		// 4. Continuar al siguiente handler (el controlador)
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
