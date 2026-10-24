package com.example.orderservice.kafka;

import com.example.orderservice.domain.OutboxMessage;
import com.example.orderservice.repository.OutboxRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Slf4j
@Component
@RequiredArgsConstructor
public class OutboxPublisher {

    private final OutboxRepository outboxRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;

    private static final String TOPIC = "orders-topic";

    @Scheduled(fixedDelay = 500)
    @Transactional
    public void publishOutboxMessages() {
        List<OutboxMessage> messages = outboxRepository.findTop50ByProcessedFalseOrderByCreatedAtAsc();
        if (messages.isEmpty()) {
            return;
        }

        for (OutboxMessage message : messages) {
            try {
                log.info("Outbox: отправка события {} для aggregateId={} в Kafka", message.getEventType(), message.getAggregateId());
                kafkaTemplate.send(TOPIC, String.valueOf(message.getAggregateId()), message.getPayload()).get();
                message.setProcessed(true);
                outboxRepository.save(message);
                log.info("Outbox: сообщение ID={} успешно отправлено и помечено processed=true", message.getId());
            } catch (Exception e) {
                log.error("Outbox: ошибка отправки сообщения ID={} в Kafka: {}. Будет повторено в следующем цикле.", message.getId(), e.getMessage());
                break;
            }
        }
    }
}
