package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/middlewares"
)

func SetupAccessEventsRoutes(api fiber.Router, ctrl *controllers.AccessEventController, apiKey string) {
	accessControl := api.Group("/access-control")

	accessControl.Post("/access-events",
		middlewares.ApiKeyMiddleware(apiKey),
		ctrl.ReceiveEvent,
	)
}