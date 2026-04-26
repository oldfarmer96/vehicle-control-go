package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/services"
	"github.com/oldfarmer96/vehicle-control-go/pkg/response"
)

type PersonaController struct {
	personaService *services.PersonaService
}

func NewPersonaController(s *services.PersonaService) *PersonaController {
	return &PersonaController{personaService: s}
}

func (c *PersonaController) Create(ctx fiber.Ctx) error {
	var payload models.CreatePersonaDTO

	if err := ctx.Bind().JSON(&payload); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Body inválido")
	}

	payload.Normalize()

	persona, err := c.personaService.CreatePersona(ctx.Context(), payload)
	if err != nil {
		if err.Error() == "el dni ya esta registrado" {
			return response.Error(ctx, fiber.StatusConflict, err.Error())
		}
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, persona)
}

func (c *PersonaController) GetAll(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	search := ctx.Query("search", "")

	result, err := c.personaService.GetAllPersonas(ctx.Context(), page, limit, search)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, result)
}

func (c *PersonaController) ToggleAccessStatus(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return response.Error(ctx, fiber.StatusBadRequest, "id es requerido")
	}

	persona, err := c.personaService.ToggleAccessStatus(ctx.Context(), id)
	if err != nil {
		if err.Error() == "persona no encontrada" {
			return response.Error(ctx, fiber.StatusNotFound, err.Error())
		}
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, persona)
}