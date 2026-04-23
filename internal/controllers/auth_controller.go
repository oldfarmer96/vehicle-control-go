// Package controllers - logica http
package controllers

import (
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

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return response.Success(c, fiber.Map{"message": "Login exitoso"})
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
	c.ClearCookie("access_token")
	return response.Success(c, fiber.Map{"message": "Logout exitoso"})
}
