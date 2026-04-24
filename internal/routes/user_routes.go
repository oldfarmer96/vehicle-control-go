package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/middlewares"
)

func SetupUserRoutes(router fiber.Router, userCtrl *controllers.UserController) {
	users := router.Group("/users", middlewares.Protected())

	users.Post("/", userCtrl.Create)
	users.Get("/", middlewares.RequireRole("ADMINISTRADOR", "CONSULTOR"), userCtrl.List)
	users.Put("/:id", middlewares.RequireRole("ADMINISTRADOR"), userCtrl.Update)
	users.Patch("/:id/toggle-active", middlewares.RequireRole("ADMINISTRADOR"), userCtrl.ToggleActive)
}
