package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/services"
	"github.com/oldfarmer96/vehicle-control-go/pkg/response"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(s *services.UserService) *UserController {
	return &UserController{userService: s}
}

func (c *UserController) Create(ctx fiber.Ctx) error {
	var payload models.CreateUserDTO
	if err := ctx.Bind().JSON(&payload); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Body inválido")
	}

	user, err := c.userService.CreateUser(ctx.Context(), payload)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, user)
}

func (c *UserController) List(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	search := ctx.Query("search", "")

	result, err := c.userService.ListUsers(ctx.Context(), page, limit, search)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, result)
}

func (c *UserController) Update(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var payload models.UpdateUserDTO
	if err := ctx.Bind().JSON(&payload); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Body inválido")
	}

	user, err := c.userService.UpdateUser(ctx.Context(), id, payload)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, user)
}

func (c *UserController) ToggleActive(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.userService.ToggleUserActive(ctx.Context(), id)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, user)
}

func (c *UserController) Profile(ctx fiber.Ctx) error {
	userID, okID := ctx.Locals("userID").(string)

	if !okID {
		return response.Error(ctx, fiber.StatusUnauthorized, "Error al leer los datos de la sesión")
	}

	user, err := c.userService.GetUserByID(ctx.Context(), userID)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			return response.Error(ctx, fiber.StatusNotFound, err.Error())
		}
		return response.Error(ctx, fiber.StatusInternalServerError, "Error interno al obtener el usuario")
	}

	return response.Success(ctx, user)
}
