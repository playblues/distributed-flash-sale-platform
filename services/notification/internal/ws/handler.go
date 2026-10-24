package ws

import (
	"log"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *Handler) HandleConnection(c *websocket.Conn) {
	userParam := c.Params("userId")
	userID, err := strconv.ParseInt(userParam, 10, 64)
	if err != nil {
		log.Printf("Неверный формат userId: %v", err)
		return
	}

	h.hub.Register(userID, c)
	defer h.hub.Unregister(userID, c)

	for {
		messageType, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("WS соединение закрыто для пользователя %d: %v", userID, err)
			break
		}
		log.Printf("Получено сообщение от пользователя %d (type %d): %s", userID, messageType, string(msg))
	}
}
