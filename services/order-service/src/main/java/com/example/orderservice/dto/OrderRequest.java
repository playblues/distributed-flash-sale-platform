package com.example.orderservice.dto;

import com.fasterxml.jackson.annotation.JsonAlias;
import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;

import java.math.BigDecimal;

public record OrderRequest(
    @JsonProperty("user_id")
    @JsonAlias({"userId", "userID"})
    @NotNull(message = "User ID не может быть пустым")
    Long userId,

    @JsonProperty("event_id")
    @JsonAlias("eventId")
    @NotNull(message = "Event ID не может быть пустым")
    Long eventId,

    @JsonProperty("price")
    @NotNull(message = "Цена не может быть пустой")
    @Positive(message = "Цена должна быть больше нуля")
    BigDecimal price
) { }

