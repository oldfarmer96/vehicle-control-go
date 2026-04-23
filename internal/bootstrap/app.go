package bootstrap

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
	"github.com/oldfarmer96/vehicle-control-go/internal/routes"
	"github.com/oldfarmer96/vehicle-control-go/internal/services"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
)

// NewApp inicializa Fiber y enlaza todas las capas de la arquitectura
func NewApp(db *pgxpool.Pool) *fiber.App {
	// Inicializamos Fiber v3
	app := fiber.New(fiber.Config{
		AppName: "Vehicle Control API v1.0",
	})

	// --- 1. Capa Store (Base de datos) ---
	userStore := store.NewUserStore(db)

	// --- 2. Capa Services (Lógica de Negocio) ---
	userService := services.NewUserService(userStore)

	// --- 3. Capa Controllers (HTTP y Validaciones) ---
	authController := controllers.NewAuthController(userStore) // Auth suele usar el store de usuarios directamente o un AuthService
	userController := controllers.NewUserController(userService)

	// --- 4. Registrar Rutas ---
	routes.Setup(app, authController, userController)

	return app
}
