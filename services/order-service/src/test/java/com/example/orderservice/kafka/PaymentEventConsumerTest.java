package com.example.orderservice.kafka;

import com.example.orderservice.dto.PaymentEvent;
import com.example.orderservice.service.OrderService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import static org.mockito.Mockito.*;

class PaymentEventConsumerTest {

    @Mock
    private ObjectMapper objectMapper;

    @Mock
    private OrderService orderService;

    @InjectMocks
    private PaymentEventConsumer paymentEventConsumer;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
    }

    @Test
    void consumePaymentEvent_Success() throws Exception {
        String message = "{\"orderId\":1,\"userId\":42,\"paymentStatus\":\"SUCCESS\"}";
        PaymentEvent event = new PaymentEvent(1L, 42L, "SUCCESS");
        
        when(objectMapper.readValue(message, PaymentEvent.class)).thenReturn(event);

        paymentEventConsumer.consumePaymentEvent(message);

        verify(objectMapper, times(1)).readValue(message, PaymentEvent.class);
        verify(orderService, times(1)).updateOrderStatus(1L, "SUCCESS");
    }

    @Test
    void consumePaymentEvent_DeserializationError_ShouldNotThrowException() throws Exception {
        String message = "invalid-json";
        
        when(objectMapper.readValue(message, PaymentEvent.class)).thenThrow(new RuntimeException("Deserialization error"));

        // Should not throw an exception because it's caught in the method
        paymentEventConsumer.consumePaymentEvent(message);

        verify(objectMapper, times(1)).readValue(message, PaymentEvent.class);
        verify(orderService, never()).updateOrderStatus(anyLong(), anyString());
    }
}
