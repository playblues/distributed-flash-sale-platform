package com.example.orderservice.controller;

import com.example.orderservice.domain.Enums.OrderStatus;
import com.example.orderservice.dto.OrderRequest;
import com.example.orderservice.dto.OrderResponse;
import com.example.orderservice.exeption.GlobalExceptionHandler;
import com.example.orderservice.service.OrderService;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.http.MediaType;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;
import org.springframework.validation.beanvalidation.LocalValidatorFactoryBean;

import java.math.BigDecimal;
import java.time.LocalDateTime;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

class OrderControllerTest {

    private MockMvc mockMvc;
    private OrderService orderService;
    private ObjectMapper objectMapper;

    @BeforeEach
    void setUp() {
        orderService = Mockito.mock(OrderService.class);
        objectMapper = new ObjectMapper().registerModule(new JavaTimeModule());

        LocalValidatorFactoryBean validatorFactoryBean = new LocalValidatorFactoryBean();
        validatorFactoryBean.afterPropertiesSet();

        mockMvc = MockMvcBuilders.standaloneSetup(new OrderController(orderService))
                .setControllerAdvice(new GlobalExceptionHandler())
                .setMessageConverters(new MappingJackson2HttpMessageConverter(objectMapper))
                .setValidator(validatorFactoryBean)
                .build();
    }

    @Test
    void deserialize_snakeCaseJson_shouldPopulateFields() throws Exception {
        String json = """
                {
                    "user_id": 42,
                    "event_id": 1,
                    "price": 5000.00
                }
                """;

        OrderRequest request = objectMapper.readValue(json, OrderRequest.class);

        assertNotNull(request);
        assertEquals(42L, request.userId());
        assertEquals(1L, request.eventId());
        assertEquals(new BigDecimal("5000.00"), request.price());
    }

    @Test
    void deserialize_camelCaseJson_shouldPopulateFields() throws Exception {
        String json = """
                {
                    "userID": 42,
                    "eventId": 1,
                    "price": 5000.00
                }
                """;

        OrderRequest request = objectMapper.readValue(json, OrderRequest.class);

        assertNotNull(request);
        assertEquals(42L, request.userId());
        assertEquals(1L, request.eventId());
        assertEquals(new BigDecimal("5000.00"), request.price());
    }

    @Test
    void createOrderApi_withSnakeCaseJson_shouldSucceed() throws Exception {
        OrderResponse response = new OrderResponse(1L, 42L, 1L, new BigDecimal("5000.00"), OrderStatus.PENDING, LocalDateTime.now());
        when(orderService.createOrder(any(OrderRequest.class))).thenReturn(response);

        String jsonPayload = """
                {
                    "user_id": 42,
                    "event_id": 1,
                    "price": 5000.00
                }
                """;

        mockMvc.perform(post("/api/v1/orders")
                .contentType(MediaType.APPLICATION_JSON)
                .content(jsonPayload))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.userId").value(42))
                .andExpect(jsonPath("$.eventId").value(1))
                .andExpect(jsonPath("$.price").value(5000.00));
    }

    @Test
    void createOrderApi_withMissingFields_shouldReturnBadRequest() throws Exception {
        String jsonPayload = """
                {
                    "price": 5000.00
                }
                """;

        mockMvc.perform(post("/api/v1/orders")
                .contentType(MediaType.APPLICATION_JSON)
                .content(jsonPayload))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.userId").value("User ID не может быть пустым"))
                .andExpect(jsonPath("$.eventId").value("Event ID не может быть пустым"));
    }
}
