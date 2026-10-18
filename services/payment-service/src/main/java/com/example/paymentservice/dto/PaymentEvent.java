package com.example.paymentservice.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record PaymentEvent (
    Long orderId,
    @JsonProperty("userId")
    Long userId,
    Long eventId,
    String paymentStatus
){}
