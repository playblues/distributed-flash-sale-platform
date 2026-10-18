package com.example.orderservice.kafka;

import com.example.orderservice.dto.PaymentEvent;
import com.example.orderservice.service.OrderService;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@RequiredArgsConstructor
public class PaymentEventConsumer {

    private final ObjectMapper objectMapper;
    private final OrderService orderService;

    @KafkaListener(topics = "payments-topic", groupId = "order-group-v2")
    public void consumePaymentEvent(String message){
        try {
            log.info("Получено сообщение из Kafka (payments-topic): {}", message);
            PaymentEvent event = objectMapper.readValue(message, PaymentEvent.class);
            log.info("Распарсено событие оплаты: orderId={}, paymentStatus={}", event.orderId(), event.paymentStatus());

            log.info("Вызов orderService.updateOrderStatus для orderId={}", event.orderId());
            orderService.updateOrderStatus(event.orderId(), event.paymentStatus());
            log.info("Статус заказа №{} успешно обновлен на {}", event.orderId(), event.paymentStatus());
        } catch (Exception e) {
            log.error("Ошибка при обработке подтверждения оплаты: {}", e.getMessage(), e);
        }
    }
}
