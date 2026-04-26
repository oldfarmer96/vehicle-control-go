package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/middlewares"
)

func SetupPersonaRoutes(router fiber.Router, personaCtrl *controllers.PersonaController) {
	personas := router.Group("/persona", middlewares.Auth())

	personas.Get("/", personaCtrl.GetAll)
	personas.Post("/", personaCtrl.Create)
}