package com.example.orderservice.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record PaymentEvent (
    @JsonProperty("orderId")
    Long orderId,
    @JsonProperty("userId")
    Long userId,
    @JsonProperty("paymentStatus")
    String paymentStatus
){ }
