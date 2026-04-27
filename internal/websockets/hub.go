// Package websockets - hub websockets
package websockets

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
)

type SocketEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = true
	log.Printf("Cliente conectado. Total: %d\n", len(h.clients))
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[conn]; ok {
		delete(h.clients, conn)
		log.Printf("Cliente desconectado. Total: %d\n", len(h.clients))
	}
}

func (h *Hub) Broadcast(eventName string, payload any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	message := SocketEvent{Event: eventName, Data: payload}

	for conn := range h.clients {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error broadcast %v", err)
			// TODO: fijarce despues
			conn.Close()
			delete(h.clients, conn)
		}
	}
}
