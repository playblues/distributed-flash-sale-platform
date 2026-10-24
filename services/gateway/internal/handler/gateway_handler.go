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
			defer resp.Body.Close()
			invBody, _ := io.ReadAll(resp.Body)
			var invErr map[string]interface{}
			_ = json.Unmarshal(invBody, &invErr)
			return c.Status(fiber.StatusConflict).JSON(invErr)
		}
		if resp != nil {
			_ = resp.Body.Close()
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
			fmt.Printf("Order Service вернул статус: %d, Ошибка: %s\n", orderResp.StatusCode, errBody)
			_ = orderResp.Body.Close()
		}

		rollbackURL := fmt.Sprintf("%s/api/v1/inventory/rollback?user_id=%v&event_id=%v", h.inventoryServiceURL, userID, req.EventID)
		rbResp, rbErr := http.Post(rollbackURL, "application/json", nil)
		if rbErr != nil {
			fmt.Printf("Ошибка выполнения отката брони в Inventory: %v\n", rbErr)
		} else {
			_ = rbResp.Body.Close()
			fmt.Printf("Успешный откат брони для user %v на event %v после сбоя Order Service\n", userID, req.EventID)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось создать заказ в системе (бронь автоматически снята)",
		})
	}
	defer orderResp.Body.Close()

	body, _ := io.ReadAll(orderResp.Body)
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusCreated).Send(body)
}
