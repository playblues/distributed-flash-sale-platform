package handler

import (
	"context"
	"fmt"
	"time"

	inventoryv1 "gateway-service/pkg/api/inventory/v1"
	orderv1 "gateway-service/pkg/api/order/v1"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GatewayHandler struct {
	jwtSecret       string
	inventoryClient inventoryv1.InventoryServiceClient
	orderClient     orderv1.OrderServiceClient
}

func NewGatewayHandler(jwtSecret string, inventoryClient inventoryv1.InventoryServiceClient, orderClient orderv1.OrderServiceClient) *GatewayHandler {
	return &GatewayHandler{
		jwtSecret:       jwtSecret,
		inventoryClient: inventoryClient,
		orderClient:     orderClient,
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
	var uid int64 = 42
	switch v := c.Locals("user_id").(type) {
	case float64:
		uid = int64(v)
	case int64:
		uid = v
	case int:
		uid = int64(v)
	}

	var req BuyRequest
	if err := c.BodyParser(&req); err != nil || req.Price == 0 {
		req.EventID = 1
		req.Price = 5000.0
	}
	if req.EventID == 0 {
		req.EventID = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reserveResp, err := h.inventoryClient.ReserveStock(ctx, &inventoryv1.ReserveStockRequest{
		EventId: req.EventID,
		UserId:  uid,
	})

	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("gRPC ошибка соединения с Inventory Service: %v (code: %s)\n", st.Message(), st.Code())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Сервис инвентаризации недоступен по gRPC",
		})
	}

	if reserveResp.Status != inventoryv1.ReservationStatus_RESERVATION_STATUS_RESERVED {
		if reserveResp.Status == inventoryv1.ReservationStatus_RESERVATION_STATUS_ALREADY_RESERVED {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status": "already_reserved",
				"error":  "У вас уже есть активная бронь на это событие",
			})
		}
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status": "sold_out",
			"error":  "Билеты закончились",
		})
	}

	orderResp, err := h.orderClient.CreateOrder(ctx, &orderv1.CreateOrderRequest{
		UserId:  uid,
		EventId: req.EventID,
		Price:   req.Price,
	})

	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("gRPC ошибка создания заказа в Order Service: %v (code: %s)\n", st.Message(), st.Code())

		rbCtx, rbCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer rbCancel()

		_, rbErr := h.inventoryClient.RollbackStock(rbCtx, &inventoryv1.RollbackStockRequest{
			EventId: req.EventID,
			UserId:  uid,
		})
		if rbErr != nil {
			fmt.Printf("Ошибка отката брони через gRPC: %v\n", rbErr)
		} else {
			fmt.Printf("Успешный откат брони через gRPC для user %v на event %v после сбоя Order Service\n", uid, req.EventID)
		}

		if st.Code() == codes.InvalidArgument {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": st.Message(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось создать заказ в системе (бронь автоматически снята)",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        orderResp.Id,
		"userId":    orderResp.UserId,
		"eventId":   orderResp.EventId,
		"price":     orderResp.Price,
		"status":    orderResp.Status.String(),
		"createdAt": orderResp.CreatedAt,
	})
}
