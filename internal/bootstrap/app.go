// Package bootstrap - para inicializar la app
package bootstrap

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/routes"
	"github.com/oldfarmer96/vehicle-control-go/internal/services"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
	"github.com/oldfarmer96/vehicle-control-go/pkg/env"
	"github.com/oldfarmer96/vehicle-control-go/pkg/response"
)

func NewApp(cfg env.Config, db *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Vehicle Control API v1.0",
	})

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.CORSURLs, ","),
		AllowCredentials: true,
	}))

	app.Get("/", func(c fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"app": "vehicle control go",
		})
	})

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

	// auth and user modules
	userStore := store.NewUserStore(db)
	userService := services.NewUserService(userStore)
	authController := controllers.NewAuthController(userStore)
	userController := controllers.NewUserController(userService)

	// persona module
	personaStore := store.NewPersonaStore(db)
	personaService := services.NewPersonaService(personaStore)
	personaController := controllers.NewPersonaController(personaService)

	// vehicle module
	vehicleStore := store.NewVehicleStore(db)
	vehicleService := services.NewVehicleService(vehicleStore, personaStore)
	vehicleController := controllers.NewVehiclecontroler(vehicleService)

	routes.Setup(app, authController, userController, vehicleController, personaController)

	return app
}
