package com.example.orderservice.grpc;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Slf4j
@Component
public class GrpcServerRunner {

    @Value("${grpc.server.port:9091}")
    private int port;

    private final OrderGrpcServiceImpl orderService;
    private Server server;

    public GrpcServerRunner(OrderGrpcServiceImpl orderService) {
        this.orderService = orderService;
    }

    @PostConstruct
    public void start() throws IOException {
        server = ServerBuilder.forPort(port)
                .addService(orderService)
                .build()
                .start();

        log.info("Order Service gRPC Server успешно запущен на порту :{}", port);

        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            log.info("Остановка gRPC сервера через JVM shutdown hook");
            stop();
        }));
    }

    @PreDestroy
    public void stop() {
        if (server != null) {
            server.shutdown();
            log.info("Order Service gRPC Server остановлен");
        }
    }
}
