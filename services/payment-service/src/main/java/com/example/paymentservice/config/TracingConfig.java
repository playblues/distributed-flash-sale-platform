package com.example.paymentservice.config;

import io.opentelemetry.exporter.otlp.http.trace.OtlpHttpSpanExporter;
import io.opentelemetry.sdk.trace.export.SpanExporter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class TracingConfig {

    @Value("${management.otlp.tracing.endpoint:http://localhost:4318/v1/traces}")
    private String otlpEndpoint;

    @Bean
    public SpanExporter otlpHttpSpanExporter() {
        return OtlpHttpSpanExporter.builder()
                .setEndpoint(otlpEndpoint)
                .build();
    }
}