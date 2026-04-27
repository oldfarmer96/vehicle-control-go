package controllers

import (
	"log"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/oldfarmer96/vehicle-control-go/internal/websockets"
)

type WSController struct {
	hub *websockets.Hub
}

func NewWSController(hub *websockets.Hub) *WSController {
	return &WSController{hub: hub}
}

func (wc *WSController) Upgrade(ctx fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (wc *WSController) HandleConnection(conn *websocket.Conn) {
	wc.hub.Register(conn)
	defer wc.hub.Unregister(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("WS cierre inesperado: %v", err)
			}
			break
		}
	}
}
