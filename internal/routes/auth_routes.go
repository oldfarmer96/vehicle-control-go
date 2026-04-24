package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/middlewares"
)

func SetupAuthRoutes(router fiber.Router, authCtrl *controllers.AuthController) {
	auth := router.Group("/auth")
	auth.Post("/login", authCtrl.Login)
	auth.Post("/logout", authCtrl.Logout, middlewares.Auth())
	router.Get("/profile", authCtrl.Profile, middlewares.Auth())
}
