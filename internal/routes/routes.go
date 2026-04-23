package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
)

// Setup orquesta todas las rutas de la aplicación
func Setup(app *fiber.App, authCtrl *controllers.AuthController, userCtrl *controllers.UserController) {
	// 1. Creamos el grupo base
	api := app.Group("/api/v1")

	// 2. Delegamos a cada archivo específico pasándole el grupo "api"
	SetupAuthRoutes(api, authCtrl)
	SetupUserRoutes(api, userCtrl)

	// Cuando tengas el de vehículos, lo agregarás aquí:
	// SetupVehicleRoutes(api, vehicleCtrl)
}
