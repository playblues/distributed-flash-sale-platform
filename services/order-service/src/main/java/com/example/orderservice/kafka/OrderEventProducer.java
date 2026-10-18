package com.example.orderservice.kafka;

import com.example.orderservice.dto.OrderResponse;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class OrderEventProducer {

    private final KafkaTemplate<String, Object> kafkaTemplate;

    private final ObjectMapper objectMapper;

    private static final String TOPIC = "orders-topic";

    public void sendOrderCreatedEvent(OrderResponse orderResponse){
        try{
            log.info("Отправка события о создании заказа №{} в Kafka", orderResponse.id());

            String jsonMessage = objectMapper.writeValueAsString(orderResponse);

            kafkaTemplate.send(TOPIC, String.valueOf(orderResponse.id()), jsonMessage);
        }
        catch (Exception e) {
            log.error("Ошибка при отправке сообщения в Kafka: {}", e.getMessage(), e);
        }
    }
}
