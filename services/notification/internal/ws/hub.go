package ws

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	clients   map[int64]map[*websocket.Conn]bool
	clientsMu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]map[*websocket.Conn]bool),
	}
}

func (h *Hub) Register(userID int64, conn *websocket.Conn) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()
	if _, exists := h.clients[userID]; !exists {
		h.clients[userID] = make(map[*websocket.Conn]bool)
	}
	h.clients[userID][conn] = true
	fmt.Printf("Пользователь %d подключился к WebSockets (активных вкладок: %d)\n", userID, len(h.clients[userID]))
}

func (h *Hub) Unregister(userID int64, conn *websocket.Conn) {
	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()
	if conns, exists := h.clients[userID]; exists {
		delete(conns, conn)
		_ = conn.Close()
		if len(conns) == 0 {
			delete(h.clients, userID)
			fmt.Printf("Все сессии пользователя %d отключены\n", userID)
		}
	}
}

func (h *Hub) SendNotification(userID int64, message string) bool {
	h.clientsMu.RLock()
	conns, exists := h.clients[userID]
	if !exists || len(conns) == 0 {
		h.clientsMu.RUnlock()
		return false
	}
	activeConns := make([]*websocket.Conn, 0, len(conns))
	for c := range conns {
		activeConns = append(activeConns, c)
	}
	h.clientsMu.RUnlock()

	sentAny := false
	for _, conn := range activeConns {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Printf("Ошибка отправки WS сессии пользователя %d: %v", userID, err)
			h.Unregister(userID, conn)
		} else {
			sentAny = true
		}
	}
	return sentAny
}
