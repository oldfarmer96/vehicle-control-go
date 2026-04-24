// Package controllers - logica http
package controllers

import (
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
	"github.com/oldfarmer96/vehicle-control-go/pkg/jwt"
	"github.com/oldfarmer96/vehicle-control-go/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	userStore *store.UserStore
}

func NewAuthController(store *store.UserStore) *AuthController {
	return &AuthController{userStore: store}
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind().JSON(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Datos inválidos")
	}

	user, err := ac.userStore.FindByUsername(c.Context(), payload.Username)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Credenciales incorrectas")
	}

	if !user.Activo {
		return response.Error(c, fiber.StatusForbidden, "Usuario desactivado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Credenciales incorrectas")
	}

	token, err := jwt.GenerateToken(user.ID, user.Rol, user.Username)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Error generando token")
	}

	err = createCookie(c, token)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, user)
}

func (ac *AuthController) Profile(c fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userRole := c.Locals("userRole").(string)
	userEmail := c.Locals("userEmail").(string)

	return c.JSON(fiber.Map{
		"id":   userID,
		"role": userRole,
		"sub":  userEmail,
	})
}

func (ac *AuthController) Logout(c fiber.Ctx) error {
	cookieName := os.Getenv("COOKIE_NAME")
	env := os.Getenv("APP_ENV")
	isProd := env == "production" || env == "prod"

	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   isProd,
		SameSite: getSameSite(isProd),
		Path:     "/",
	})

	return response.Success(c, fiber.Map{"message": "Sesión cerrada"})
}

// helpers
func createCookie(c fiber.Ctx, token string) error {
	env := os.Getenv("APP_ENV")
	cookieName := os.Getenv("COOKIE_NAME")

	if cookieName == "" {
		return errors.New("falta el nombre de la cookie")
	}

	isProd := env == "production" || env == "prod"

	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   isProd,
		SameSite: getSameSite(isProd),
		Path:     "/",
	})

	return nil
}

func getSameSite(isProduction bool) string {
	if isProduction {
		return "Strict"
	}
	return "Lax"
}
