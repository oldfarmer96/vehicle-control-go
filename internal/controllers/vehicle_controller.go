package controllers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/services"
	"github.com/oldfarmer96/vehicle-control-go/pkg/response"
)

var validate = validator.New()

type VehicleController struct {
	vehicleService *services.VehicleService
}

func NewVehiclecontroler(s *services.VehicleService) *VehicleController {
	return &VehicleController{vehicleService: s}
}

func (c *VehicleController) Create(ctx fiber.Ctx) error {
	var payload models.CreaateVehicleDTO

	if err := ctx.Bind().JSON(&payload); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Body inválido")
	}

	payload.Normalize()

	if err := validate.Struct(&payload); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Error de validación")
	}

	vehicle, err := c.vehicleService.CreateVehicle(ctx.Context(), payload)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, vehicle)
}

func (c *VehicleController) GetAll(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	placa := ctx.Query("placa", "")

	result, err := c.vehicleService.GetAllVehicles(ctx.Context(), page, limit, placa)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, result)
}

func (c *VehicleController) GetByPlaca(ctx fiber.Ctx) error {
	placa := ctx.Params("placa")
	if placa == "" {
		return response.Error(ctx, fiber.StatusBadRequest, "placa es requerida")
	}

	vehicle, err := c.vehicleService.GetVehicleByPlaca(ctx.Context(), placa)
	if err != nil {
		if err.Error() == "vehiculo no encontrado" {
			return response.Error(ctx, fiber.StatusNotFound, err.Error())
		}
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, vehicle)
}

func (c *VehicleController) AssignOwner(ctx fiber.Ctx) error {
	vehiculoID := ctx.Params("id")
	if vehiculoID == "" {
		return response.Error(ctx, fiber.StatusBadRequest, "id del vehiculo es requerido")
	}

	var payload models.AssignOwnerDTO
	if err := ctx.Bind().JSON(&payload); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, "Body inválido")
	}

	if payload.PersonaID == "" {
		return response.Error(ctx, fiber.StatusBadRequest, "personaId es requerido")
	}

	vehicle, err := c.vehicleService.AssignOwner(ctx.Context(), vehiculoID, payload.PersonaID)
	if err != nil {
		if err.Error() == "vehiculo no encontrado" {
			return response.Error(ctx, fiber.StatusNotFound, err.Error())
		}
		if err.Error() == "persona no encontrada" {
			return response.Error(ctx, fiber.StatusNotFound, err.Error())
		}
		if err.Error() == "esta persona ya esta asignada a este vehiculo" {
			return response.Error(ctx, fiber.StatusConflict, err.Error())
		}
		return response.Error(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(ctx, fiber.Map{
		"mensaje": "propietario asignado exitosamente",
		"vehiculo": vehicle,
	})
}
