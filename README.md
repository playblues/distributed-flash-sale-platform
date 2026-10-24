# Распределенная event-driven платформа экспресс-бронирования билетов (Flash-Sale Engine)

[![Java](https://img.shields.io/badge/Java-21-orange.svg?style=flat&logo=openjdk)](https://www.oracle.com/java/)
[![Spring Boot](https://img.shields.io/badge/Spring%20Boot-3.3.4-brightgreen.svg?style=flat&logo=springboot)](https://spring.io/projects/spring-boot)
[![Go](https://img.shields.io/badge/Go-1.23%2B-blue.svg?style=flat&logo=go)](https://golang.org/)
[![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-KRaft-red.svg?style=flat&logo=apachekafka)](https://kafka.apache.org/)
[![Redis](https://img.shields.io/badge/Redis-7.0-red.svg?style=flat&logo=redis)](https://redis.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue.svg?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED.svg?style=flat&logo=docker)](https://www.docker.com/)

Полиглотная событийно-ориентированная микросервисная платформа для обработки пиковых нагрузок при продаже билетов и flash-распродажах.

Система гарантирует **полное отсутствие овербукинга (0 Race Conditions)**, **согласованность в конечном счете (Eventual Consistency)**, **защиту от двойных списаний (Transactional Outbox)**, **двухфазную идемпотентность** и **рассылку статусов в реальном времени через WebSockets**.

---

## Архитектура системы

Высокопроизводительный слой приема трафика, лимитов и кэширования на **Go (Fiber)** вынесен на периметр и защищает транзакционное ядро на **Java (Spring Boot)**. Межсервисное взаимодействие между сервисами заказов, платежей, остатков и уведомлений организовано асинхронно через шину событий **Apache Kafka**.

```mermaid
flowchart TD
    Client([Web UI / Клиент]) -->|HTTP / WebSockets| Gateway[API Gateway Go / Fiber]
    
    subgraph Go Edge Services
        Gateway -->|1. Sliding Window & Идемпотентность| Redis[(Redis 7)]
        Gateway -->|2. Бронь билета| Inventory[Inventory Service Go]
        Inventory -->|3. Атомарный Lua-скрипт| Redis
    end
    
    subgraph Java Core Services
        Gateway -->|4. Создать заказ| Order[Order Service Spring Boot]
        Order -->|5. ACID транзакция: Order + Outbox| Postgres[PostgreSQL 16]
        Order -->|6. Outbox Publisher| Kafka{Apache Kafka}
        
        Kafka -->|7. Вычитать заказ| Payment[Payment Service Spring Boot]
        Payment -->|8. Провести платеж| Postgres
        Payment -->|9. Отправить статус оплаты| Kafka
    end
    
    subgraph Event-Driven Reactions
        Kafka -->|10. Обновить статус PAID / CANCELLED| Order
        Kafka -->|11. Откат брони при FAILED| Inventory
        Kafka -->|12. Стрим событий| Notification[Notification Service Go]
        Notification -->|13. Push в браузер| Client
    end
```

---

## Ключевые инженерные решения и паттерны

* **Защита от овербукинга (Redis Lua):** Бронирование остатков выполняется атомарным скриптом на языке Lua внутри памяти Redis. Это полностью исключает состояние гонки (Race Condition) при параллельных запросах сотен пользователей.
* **1 билет в одни руки:** Lua-скрипт проверяет наличие уже существующей активной брони у пользователя (`KEYS[2]`). Если пользователь уже забронировал билет, повторный запрос отклоняется (`409 already_reserved`).
* **Transactional Outbox (Защита от Dual-Write):** `Order Service` сохраняет заказ и событие для Kafka в рамках единой неделимой транзакции PostgreSQL. Фоновый паблишер гарантированно доставляет события в брокер с семантикой At-Least-Once даже при временном отказе Kafka.
* **Двухфазная идемпотентность (`X-Idempotency-Key`):** Шлюз атомарно блокирует ключ (`SETNX PROCESSING 30s`) на время обработки запроса, исключая дублирование заказов при одновременных кликах. При успехе результат кэшируется на 24 часа; при ошибке блокировка немедленно снимается для возможности повтора.
* **Компенсация на шлюзе (Gateway Saga Rollback):** Если бронирование в Redis прошло, но создание заказа в `Order Service` завершилось сбоем сети, шлюз автоматически вызывает `/inventory/rollback`, предотвращая утечку билетов.
* **Saga Orchestration & Compensation:** Если платеж отклоняется (`price > 100 000`), `Payment Service` публикует статус `FAILED`; `Order Service` переводит заказ в `CANCELLED`, а `Inventory Service` асинхронно возвращает билет на склад в Redis (`INCR`).
* **Dead Letter Queue (DLQ / DLT):** «Ядовитые» сообщения (`price == 666`), вызывающие критический сбой после 3 автоматических ретраев, изолируются в топик `orders-topic.DLT` с диагностическими заголовками.
* **Real-Time WebSockets:** Потокобезопасный Hub на Go рассылает мгновенные пуш-уведомления об изменении статусов оплаты прямо в браузер пользователя с поддержкой множественных вкладок одного пользователя.
* **Sliding Window Rate Limiter:** Защита шлюза от DDoS и спама на базе структур данных Redis Sorted Sets (`ZSET`).

---

## Честный запуск проекта (пошагово)

### 1. Что должно быть установлено
* **[Docker Desktop](https://www.docker.com/)** (рекомендуется выделить в настройках Docker от 4 ГБ оперативной памяти, так как поднимаются 10 контейнеров: Kafka KRaft, PostgreSQL, Redis, Jaeger, AKHQ и 5 микросервисов).
* **Node.js** (для быстрого запуска веб-интерфейса через `npx serve`) или **Python 3**.

---

### 2. Запуск всего бэкенда в Docker

В корне проекта выполните одну команду:

```bash
docker compose up --build -d
```

> [!TIP]
> Сборка образов займет 1-3 минуты при первом запуске (скачаются базовые образы Go и Maven, скомпилируются бинарники).

Проверить, что все 10 контейнеров успешно запустились и работают:

```bash
docker compose ps
```

Если нужно посмотреть логи конкретного сервиса:
```bash
docker compose logs -f order-service
# или
docker compose logs -f gateway-service
```

---

### 3. Запуск веб-дашборда (React Frontend)

Дашборд написан на **React + Vite + Tailwind CSS**:

```bash
cd frontend
npm install
npm run dev
```

👉 **Откройте дашборд в браузере:** [http://localhost:5500](http://localhost:5500)

---

### 4. Доступные сервисы и панели управления

| Сервис | Адрес | Что там смотреть |
| :--- | :--- | :--- |
| **Веб-дашборд платформы** | [http://localhost:5500](http://localhost:5500) | Интерактивный UI с живым трекером Saga, складом и логами |
| **Панель Kafka (AKHQ)** | [http://localhost:8080](http://localhost:8080) | Просмотр топиков `orders-topic`, `payments-topic`, `orders-topic.DLT` |
| **Трассировка (Jaeger)** | [http://localhost:16686](http://localhost:16686) | Сквозные спаны и задержки распределенных запросов |
| **API Gateway** | [http://localhost:3000](http://localhost:3000) | Точка входа REST API |
| **Inventory Service** | [http://localhost:3001](http://localhost:3001) | Управление остатками и Redis |
| **Notification WebSockets**| `ws://localhost:3002/ws/:userId` | Поток событий реального времени |

---

## Интерактивный тест-драйв (за 2 минуты)

В веб-интерфейсе реализованы быстрые кнопки для проверки всех распределенных сценариев:

1. **Завоз билетов:** В карточке события нажмите кнопку **`+10`** — остаток билетов в Redis станет равен 10.
2. **Сценарий 1 (Успешная покупка):** Нажмите **`🟢 Успех (5 000 ₽)`**.
   * Шлюз проверяет JWT и Rate Limit.
   * Inventory атомарно списывает билет в Redis.
   * `Order Service` сохраняет заказ в Postgres и записывает событие в Outbox.
   * `OutboxPublisher` доставляет событие в Kafka.
   * `Payment Service` проводит платеж и публикует `SUCCESS`.
   * Браузер мгновенно получает push по WebSocket, а статус заказа меняется на **PAID**.
3. **Сценарий 2 (Saga Rollback):** Нажмите **`🟡 Saga Rollback`** (сумма 150 000 ₽).
   * Платеж отклоняется банком (`price > 100 000`).
   * Заказ переходит в статус `CANCELLED`.
   * `Inventory Service` ловит событие отката и возвращает билет обратно на склад в Redis (`INCR`).
4. **Сценарий 3 (Poison Pill & DLT):** Нажмите **`🟣 Poison Pill (666 ₽)`**.
   * Имитируется критический сбой базы данных в `Payment Service`.
   * Сообщение автоматически ретраится 3 раза с интервалом в 1 секунду.
   * После исчерпания попыток сообщение изолируется в топик `orders-topic.DLT` (можно увидеть в [AKHQ](http://localhost:8080)).
5. **Сценарий 4 (Flash-Sale Spurt):** Нажмите **`⚡ Залп (5 покупок)`**.
   * Платформа отправляет 5 параллельных запросов от разных пользователей одновременно.
   * Счетчик билетов плавно уменьшается, исключая перепродажи.

---

## Запуск нагрузочных тестов (K6)

Для симуляции пикового наплыва покупателей:

```bash
k6 run load-tests/flash_sale_test.js
```

Тест запускает до 100 виртуальных пользователей, проверяет отсутствие 5xx ошибок, корректную отсечку распроданных билетов (`409 Conflict`) и стабильность очередей.

---

## Остановка проекта

Чтобы корректно остановить все контейнеры:

```bash
docker compose down
```

Если необходимо также удалить накопленные данные БД и Kafka:
```bash
docker compose down -v
```