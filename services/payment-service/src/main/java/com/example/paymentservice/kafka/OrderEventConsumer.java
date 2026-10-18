package com.example.paymentservice.kafka;

import com.example.paymentservice.dto.OrderEvent;
import com.example.paymentservice.service.PaymentService;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@RequiredArgsConstructor
public class OrderEventConsumer {

    private final ObjectMapper objectMapper;
    private final PaymentService paymentService;

    @KafkaListener(topics = "orders-topic", groupId = "payment-group")
    public void consumeOrderEvent(String message){
        log.info("Поймано новое сообщение из Kafka: {}", message);
        try{
            OrderEvent orderEvent = objectMapper.readValue(message, OrderEvent.class);
            paymentService.processPayment(orderEvent);
        } catch (Exception e){
            log.error("Ошибка при обработки сообщения из Kafka: {}", e.getMessage());
            throw new RuntimeException("Сбой обработки заказа в Kafka", e);
        }
    }
}
