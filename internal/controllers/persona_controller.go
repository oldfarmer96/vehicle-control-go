package controllers

import (
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