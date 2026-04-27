package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/services"
	"github.com/oldfarmer96/vehicle-control-go/internal/websockets"
)

type AccessEventController struct {
	hub     *websockets.Hub
	service *services.AccessEventService
}

func NewAccessEventController(hub *websockets.Hub, service *services.AccessEventService) *AccessEventController {
	return &AccessEventController{
		hub:     hub,
		service: service,
	}
}

func (c *AccessEventController) ReceiveEvent(ctx fiber.Ctx) error {
	var req models.AccessEventRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Body inválido",
		})
	}

	resp, err := c.service.ProcessAccessEvent(ctx, &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.hub.Broadcast("new-access-event", resp)

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

