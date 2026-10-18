package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"notification-service/internal/ws"

	"github.com/segmentio/kafka-go"
)

type PaymentEvent struct {
	OrderID       int64  `json:"orderId"`
	UserID        int64  `json:"userId"`
	PaymentStatus string `json:"paymentStatus"`
}

func StartConsumer(broker string, hub *ws.Hub) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "payments-topic",
		GroupID:  "notification-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()

	fmt.Println("Notification Service слушает payments-topic...")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка чтения Kafka: %v", err)
			continue
		}

		var event PaymentEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		fmt.Printf("Получено событие об оплате: Заказ %d, Статус %s (Пользователь %d)\n", event.OrderID, event.PaymentStatus, event.UserID)

		msg := fmt.Sprintf("Статус вашего заказа №%d изменился на: %s", event.OrderID, event.PaymentStatus)
		if ok := hub.SendNotification(event.UserID, msg); ok {
			fmt.Printf("Push-уведомление успешно отправлено пользователю %d\n", event.UserID)
		} else {
			fmt.Printf("Пользователь %d сейчас не онлайн (Push не отправлен)\n", event.UserID)
		}
	}
}
