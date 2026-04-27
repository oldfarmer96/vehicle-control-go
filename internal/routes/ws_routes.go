package routes

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/controllers"
)

func SetupWSRoutes(router fiber.Router, wsCtrl *controllers.WSController) {
	ws := router.Group("/access-control")
	ws.Use(wsCtrl.Upgrade)
	ws.Get("/", websocket.New(wsCtrl.HandleConnection))
}
