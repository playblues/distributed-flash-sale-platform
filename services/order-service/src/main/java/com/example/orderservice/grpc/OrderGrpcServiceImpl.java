package com.example.orderservice.grpc;

import com.example.orderservice.dto.OrderRequest;
import com.example.orderservice.dto.OrderResponse;
import com.example.orderservice.grpc.order.CreateOrderRequest;
import com.example.orderservice.grpc.order.CreateOrderResponse;
import com.example.orderservice.grpc.order.OrderStatus;
import com.example.orderservice.grpc.order.OrderServiceGrpc;
import com.example.orderservice.service.OrderService;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

import java.math.BigDecimal;

@Slf4j
@Component
public class OrderGrpcServiceImpl extends OrderServiceGrpc.OrderServiceImplBase {

    private final OrderService orderService;

    public OrderGrpcServiceImpl(OrderService orderService) {
        this.orderService = orderService;
    }

    @Override
    public void createOrder(CreateOrderRequest request, StreamObserver<CreateOrderResponse> responseObserver) {
        try {
            if (request.getUserId() <= 0 || request.getEventId() <= 0 || request.getPrice() <= 0) {
                responseObserver.onError(Status.INVALID_ARGUMENT
                        .withDescription("user_id, event_id and price must be greater than zero")
                        .asRuntimeException());
                return;
            }

            log.info("gRPC CreateOrder вызван для user_id={}, event_id={}, price={}",
                    request.getUserId(), request.getEventId(), request.getPrice());

            OrderRequest orderRequest = new OrderRequest(
                    request.getUserId(),
                    request.getEventId(),
                    BigDecimal.valueOf(request.getPrice())
            );

            OrderResponse response = orderService.createOrder(orderRequest);

            OrderStatus grpcStatus = switch (response.status()) {
                case PENDING -> OrderStatus.ORDER_STATUS_PENDING;
                case PAID -> OrderStatus.ORDER_STATUS_PAID;
                case CANCELLED -> OrderStatus.ORDER_STATUS_CANCELLED;
            };

            CreateOrderResponse grpcResponse = CreateOrderResponse.newBuilder()
                    .setId(response.id())
                    .setUserId(response.userId())
                    .setEventId(response.eventId())
                    .setPrice(response.price().doubleValue())
                    .setStatus(grpcStatus)
                    .setCreatedAt(response.createdAt() != null ? response.createdAt().toString() : "")
                    .build();

            responseObserver.onNext(grpcResponse);
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("Ошибка при обработке gRPC CreateOrder: {}", e.getMessage(), e);
            responseObserver.onError(Status.INTERNAL
                    .withDescription("Internal server error: " + e.getMessage())
                    .asRuntimeException());
        }
    }
}
