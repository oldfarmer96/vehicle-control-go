package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/middlewares"
)

func SetupVehicleRoutes(router fiber.Router, vehicleCtrl *controllers.VehicleController) {
	vehicles := router.Group("/vehicle", middlewares.Auth())

	vehicles.Get("/", vehicleCtrl.GetAll)
	vehicles.Get("/:placa/placa", vehicleCtrl.GetByPlaca)
	vehicles.Post("/", vehicleCtrl.Create)
	vehicles.Post("/:id/assign-owner", vehicleCtrl.AssignOwner)
}
