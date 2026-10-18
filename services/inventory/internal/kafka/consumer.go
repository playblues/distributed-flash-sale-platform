package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"inventory-service/internal/repository"

	"github.com/segmentio/kafka-go"
)

type PaymentEvent struct {
	OrderID       int64  `json:"orderId"`
	UserID        int64  `json:"userId"`
	EventID       int64  `json:"eventId"`
	PaymentStatus string `json:"paymentStatus"`
}

func StartPaymentConsumer(repo *repository.InventoryRepository) {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "payments-topic",
		GroupID:  "inventory-rollback-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()

	fmt.Println("Inventory Service слушает payments-topic на случай откатов")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка чтения Kafka: %v", err)
			continue
		}

		var event PaymentEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			continue
		}

		if event.PaymentStatus == "FAILED" {
			log.Printf("Оплата отклонена (Заказ %d). Откат брони...", event.OrderID)

			eventIDStr := fmt.Sprintf("%d", event.OrderID)
			userIDStr := fmt.Sprintf("%d", event.UserID)

			err := repo.RollBack(context.Background(), eventIDStr, userIDStr)
			if err != nil {
				log.Printf("Ошибка отката в Redis: %v", err)
			} else {
				log.Printf("Билет на Event %s успешно возвращен на склад", eventIDStr)
			}
		}
	}
}
