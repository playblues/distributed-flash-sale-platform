package com.example.paymentservice.service;

import com.example.paymentservice.domain.Payment;
import com.example.paymentservice.dto.OrderEvent;
import com.example.paymentservice.dto.PaymentEvent;
import com.example.paymentservice.repository.PaymentRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Slf4j
@Service
public class PaymentService {

    private final PaymentRepository paymentRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;
    private final ObjectMapper objectMapper;

    private static final String PAYMENTS_TOPIC = "payments-topic";

    public PaymentService(PaymentRepository paymentRepository,
                          KafkaTemplate<String, Object> kafkaTemplate,
                          ObjectMapper objectMapper) {
        this.paymentRepository = paymentRepository;
        this.kafkaTemplate = kafkaTemplate;
        this.objectMapper = objectMapper;
    }

    @Transactional
    public void processPayment(OrderEvent orderEvent){
        log.info("Начало обработки платежа для заказа №{}", orderEvent.id());

        if (orderEvent.price().doubleValue() == 666.0) {
            log.error("Имитация сбоя БД при обработке заказа №{}", orderEvent.id());
            throw new RuntimeException("Симуляция падения базы данных (Poison Pill)");
        }

        String paymentStatus = "SUCCESS";
        if (orderEvent.price().doubleValue() > 100_000){
            paymentStatus = "FAILED";
        }

        Payment payment = Payment.builder()
                .orderId(orderEvent.id())
                .amount(orderEvent.price())
                .status(paymentStatus)
                .createdAt(LocalDateTime.now())
                .build();
        paymentRepository.save(payment);
        log.info("Платёж сохранен в БД со статусом {}", paymentStatus);

        try{
            PaymentEvent paymentEvent = new PaymentEvent(orderEvent.id(), orderEvent.userId(), orderEvent.eventId(), paymentStatus);
            String jsonMessage = objectMapper.writeValueAsString(paymentEvent);

            log.info("Отправка оплаты для заказа №{} пользователя {}", orderEvent.id(), orderEvent.userId());
            kafkaTemplate.send(PAYMENTS_TOPIC, String.valueOf(orderEvent.id()), jsonMessage);
            log.info("Событие оплаты отправлено в топик {}", PAYMENTS_TOPIC);
        } catch (Exception e) {
            log.error("Ошибка отправки в Kafka: {}", e.getMessage());
        }
    }
}
