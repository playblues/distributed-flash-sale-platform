package handler

import (
	"context"
	"fmt"
	"strconv"

	"inventory-service/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type InventoryHandler struct {
	repo *repository.InventoryRepository
}

func NewInventoryHandler(repo *repository.InventoryRepository) *InventoryHandler {
	return &InventoryHandler{repo: repo}
}

func (h *InventoryHandler) InitStock(c *fiber.Ctx) error {
	eventID := c.Query("event_id", "1")
	amountStr := c.Query("amount", "5")

	amount, err := strconv.Atoi(amountStr)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверное количество",
		})
	}

	ctx := context.Background()

	if err := h.repo.InitializeStock(ctx, eventID, amount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка инициализации",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Склад успешно инициализирован",
	})
}

func (h *InventoryHandler) GetStock(c *fiber.Ctx) error {
	eventID := c.Query("event_id", "1")
	stock, err := h.repo.GetStock(context.Background(), eventID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить остаток склада",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"eventId": eventID,
		"stock":   stock,
	})
}

func (h *InventoryHandler) Reserve(c *fiber.Ctx) error {
	eventID := c.Query("event_id", "1")
	userID := c.Query("user_id")
	if userID == "" {
		userID = c.Get("user_id", "guest")
	}

	ctx := context.Background()
	result, err := h.repo.ReserveStock(ctx, eventID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка БД",
		})
	}
	if result == -1 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status": "already_reserved",
			"error":  "У вас уже есть активная бронь на это событие",
		})
	}
	if result == 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status": "sold_out",
			"error":  "Билеты закончились",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "reserved",
		"message": fmt.Sprintf("Билет забронирован на 10 минут для пользователя %s", userID),
	})
}

func (h *InventoryHandler) Rollback(c *fiber.Ctx) error {
	eventID := c.Query("event_id", "1")
	userID := c.Query("user_id")
	if userID == "" {
		userID = c.Get("user_id", "guest")
	}

	ctx := context.Background()
	if err := h.repo.RollBack(ctx, eventID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка отката брони",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "rolled_back",
		"message": fmt.Sprintf("Бронь успешно отменена для пользователя %s", userID),
	})
}
