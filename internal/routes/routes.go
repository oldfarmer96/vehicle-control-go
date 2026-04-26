// Package routes - set for all routes
package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
)

func Setup(app *fiber.App, authCtrl *controllers.AuthController, userCtrl *controllers.UserController, vehicleCtrl *controllers.VehicleController, personaCtrl *controllers.PersonaController) {
	api := app.Group("/api/v1")

	SetupAuthRoutes(api, authCtrl)
	SetupUserRoutes(api, userCtrl)
	SetupVehicleRoutes(api, vehicleCtrl)
	SetupPersonaRoutes(api, personaCtrl)
}
