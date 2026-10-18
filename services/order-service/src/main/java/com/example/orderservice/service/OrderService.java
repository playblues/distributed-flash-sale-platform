package com.example.orderservice.service;

import com.example.orderservice.domain.Enums.OrderStatus;
import com.example.orderservice.domain.Order;
import com.example.orderservice.dto.OrderRequest;
import com.example.orderservice.dto.OrderResponse;
import com.example.orderservice.kafka.OrderEventProducer;
import com.example.orderservice.repository.OrderRepository;
import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Slf4j
@Service
@RequiredArgsConstructor
public class OrderService {
    private final OrderRepository orderRepository;
    private final OrderEventProducer eventProducer;

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

        eventProducer.sendOrderCreatedEvent(response);

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
