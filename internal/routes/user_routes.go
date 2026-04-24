package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/middlewares"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
)

func SetupUserRoutes(router fiber.Router, userCtrl *controllers.UserController) {
	users := router.Group("/users", middlewares.Auth())

	users.Post("/", userCtrl.Create)
	users.Get("/", middlewares.UserRole(models.RoleAdmin, models.RoleConsultant), userCtrl.List)
	users.Get("/profile", userCtrl.Profile)
	users.Put("/:id", middlewares.UserRole(models.RoleAdmin), userCtrl.Update)
	users.Patch("/:id/toggle-active", middlewares.UserRole(models.RoleAdmin, models.RoleConsultant), userCtrl.ToggleActive)
}
