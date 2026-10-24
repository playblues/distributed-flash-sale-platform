package com.example.orderservice.service;

import com.example.orderservice.domain.Enums.OrderStatus;
import com.example.orderservice.domain.Order;
import com.example.orderservice.dto.OrderRequest;
import com.example.orderservice.dto.OrderResponse;
import com.example.orderservice.kafka.OrderEventProducer;
import com.example.orderservice.domain.OutboxMessage;
import com.example.orderservice.repository.OrderRepository;
import com.example.orderservice.repository.OutboxRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Slf4j
@Service
public class OrderService {
    private final OrderRepository orderRepository;
    private final OutboxRepository outboxRepository;
    private final OrderEventProducer eventProducer;
    private final ObjectMapper objectMapper;

    public OrderService(OrderRepository orderRepository,
                        OutboxRepository outboxRepository,
                        OrderEventProducer eventProducer,
                        ObjectMapper objectMapper) {
        this.orderRepository = orderRepository;
        this.outboxRepository = outboxRepository;
        this.eventProducer = eventProducer;
        this.objectMapper = objectMapper;
    }

    @Transactional
    public OrderResponse createOrder(OrderRequest request){
        Order order = Order.builder()
                .userId(request.userId())
                .eventId(request.eventId())
                .price(request.price())
                .status(OrderStatus.PENDING)
                .createdAt(LocalDateTime.now())
                .build();

        Order savedOrder = orderRepository.save(order);

        OrderResponse response = mapToResponse(savedOrder);

        try {
            String payload = objectMapper.writeValueAsString(response);
            OutboxMessage outboxMessage = OutboxMessage.builder()
                    .aggregateType("ORDER")
                    .aggregateId(savedOrder.getId())
                    .eventType("ORDER_CREATED")
                    .payload(payload)
                    .processed(false)
                    .createdAt(LocalDateTime.now())
                    .build();
            outboxRepository.save(outboxMessage);
            log.info("Событие для заказа №{} успешно записано в Outbox", savedOrder.getId());
        } catch (Exception e) {
            log.error("Ошибка сохранения сообщения в Outbox: {}", e.getMessage(), e);
            throw new RuntimeException("Не удалось сохранить событие заказа в Outbox", e);
        }

        return response;
    }
    @Transactional
    public void updateOrderStatus(Long orderId, String paymentStatus){
        log.info("Поиск заказа по orderId={} в БД", orderId);
        orderRepository.findById(orderId).ifPresentOrElse(order ->{
            log.info("Заказ №{} найден со статусом {}", orderId, order.getStatus());
            if ("SUCCESS".equals(paymentStatus)){
                order.setStatus(OrderStatus.PAID);
            } else {
                order.setStatus(OrderStatus.CANCELLED);
            }
            orderRepository.save(order);
            log.info("Сохранен заказ №{} с новым статусом {}", orderId, order.getStatus());
        }, () -> {
            log.warn("Заказ №{} не найден в БД при обновлении статуса оплаты", orderId);
        });
    }
    private OrderResponse mapToResponse(Order order){
        return new OrderResponse(
                order.getId(),
                order.getUserId(),
                order.getEventId(),
                order.getPrice(),
                order.getStatus(),
                order.getCreatedAt()
        );
    }
}
