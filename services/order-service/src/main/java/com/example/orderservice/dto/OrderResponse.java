package com.example.orderservice.dto;

import com.example.orderservice.domain.Enums.OrderStatus;

import java.math.BigDecimal;
import java.time.LocalDateTime;

import com.fasterxml.jackson.annotation.JsonProperty;

public record OrderResponse (
    @JsonProperty("id")
    Long id,
    @JsonProperty("userId")
    Long userId,
    @JsonProperty("eventId")
    Long eventId,
    @JsonProperty("price")
    BigDecimal price,
    @JsonProperty("status")
    OrderStatus status,
    @JsonProperty("createdAt")
    LocalDateTime createdAt
) {}
