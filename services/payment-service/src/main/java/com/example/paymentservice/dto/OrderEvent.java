package com.example.paymentservice.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@JsonIgnoreProperties(ignoreUnknown = true)
public record OrderEvent (
    Long id,
    @JsonProperty("userId")
    Long userId,
    @JsonProperty("eventId")
    Long eventId,
    BigDecimal price,
    String status,
    LocalDateTime createdAt
){}
