package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type GatewayHandler struct {
	jwtSecret           string
	inventoryServiceURL string
	orderServiceURL     string
}

func NewGatewayHandler(jwtSecret string, inventoryServiceURL string, orderServiceURL string) *GatewayHandler {
	return &GatewayHandler{
		jwtSecret:           jwtSecret,
		inventoryServiceURL: inventoryServiceURL,
		orderServiceURL:     orderServiceURL,
	}
}

func (h *GatewayHandler) Login(c *fiber.Ctx) error {
	claims := jwt.MapClaims{
		"user_id":  42,
		"username": "admin",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось создать токен",
		})
	}
	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}

type BuyRequest struct {
	EventID int64   `json:"eventId"`
	Price   float64 `json:"price"`
}

func (h *GatewayHandler) Buy(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	var req BuyRequest

	if err := c.BodyParser(&req); err != nil || req.Price == 0 {
		req.EventID = 1
		req.Price = 5000.0
	}
	if req.EventID == 0 {
		req.EventID = 1
	}

	inventoryURL := fmt.Sprintf("%s/api/v1/inventory/reserve?user_id=%v&event_id=%v", h.inventoryServiceURL, userID, req.EventID)
	resp, err := http.Post(inventoryURL, "application/json", nil)

	if err != nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			fmt.Printf("Ошибка соединения с Inventory: %v\n", err)
		} else {
			fmt.Printf("Inventory Service вернул статус: %d\n", resp.StatusCode)
		}
		if resp != nil && resp.StatusCode == http.StatusConflict {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status": "sold_out",
				"error":  "Билеты закончились",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Сервис инвентаризации недоступен",
		})
	}
	_ = resp.Body.Close()

	orderPayload := map[string]interface{}{
		"userId":  userID,
		"eventId": req.EventID,
		"price":   req.Price,
	}
	payloadBytes, _ := json.Marshal(orderPayload)

	orderURL := fmt.Sprintf("%s/api/v1/orders", h.orderServiceURL)
	orderResp, err := http.Post(orderURL, "application/json", bytes.NewBuffer(payloadBytes))

	if err != nil || orderResp.StatusCode != http.StatusCreated {
		if err != nil {
			fmt.Printf("Ошибка соединения с Order Service: %v\n", err)
		} else {
			errBody, _ := io.ReadAll(orderResp.Body)
			fmt.Printf("Order Service вернул статус: %d, Ошибка: %s", orderResp.StatusCode, errBody)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось создать заказ в системе",
		})
	}
	defer orderResp.Body.Close()

	body, _ := io.ReadAll(orderResp.Body)
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusCreated).Send(body)
}
