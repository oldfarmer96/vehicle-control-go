// Package bootstrap - para inicializar la app
package bootstrap

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/routes"
	"github.com/oldfarmer96/vehicle-control-go/internal/services"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
	"github.com/oldfarmer96/vehicle-control-go/pkg/response"
)

func NewApp(db *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Vehicle Control API v1.0",
	})

	app.Use(logger.New())

	app.Get("/health", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
		defer cancel()

		dbStatus := "connected"
		if err := db.Ping(ctx); err != nil {
			dbStatus = "disconnected"
		}

		return response.Success(c, fiber.Map{
			"status":    "ok",
			"database":  dbStatus,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	userStore := store.NewUserStore(db)
	userService := services.NewUserService(userStore)
	authController := controllers.NewAuthController(userStore)
	userController := controllers.NewUserController(userService)

	routes.Setup(app, authController, userController)

	return app
}
