package ws

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	clients   map[int64]*websocket.Conn
	clientsMu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]*websocket.Conn),
	}
}

func (h *Hub) Register(userID int64, conn *websocket.Conn) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()
	h.clients[userID] = conn
	fmt.Printf("Пользователь %d подключился к WebSockets\n", userID)
}

func (h *Hub) Unregister(userID int64) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()
	if conn, exists := h.clients[userID]; exists {
		_ = conn.Close()
		delete(h.clients, userID)
		fmt.Printf("Пользователь %d отключился\n", userID)
	}
}

func (h *Hub) SendNotification(userID int64, message string) bool {
	h.clientsMu.RLock()
	conn, exists := h.clients[userID]
	h.clientsMu.RUnlock()

	if !exists {
		return false
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Printf("Ошибка отправки WS пользователю %d: %v", userID, err)
		h.Unregister(userID)
		return false
	}
	return true
}
