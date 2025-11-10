# ChatServer

Микросервис управления чатами и сообщениями для мессенджера. Обеспечивает создание чатов, отправку сообщений в реальном времени через gRPC streaming, управление участниками и их ролями.

## Функциональность

- **CreateChat** - создание нового чата (требует JWT)
- **DeleteChat** - удаление чата (только владелец, требует JWT)
- **GetChat** - получение информации о чате и его сообщений (требует JWT)
- **ConnectChat** - подключение к чату через gRPC streaming для получения сообщений в реальном времени (требует JWT)
- **AddUser** - добавление пользователя в чат (требует JWT, роль admin/owner)
- **RemoveUser** - удаление пользователя из чата (требует JWT, роль admin/owner)
- **UpdateUserRole** - изменение роли участника (требует JWT, роль owner)
- **SendMessage** - отправка сообщения в чат (требует JWT)
- **EditMessage** - редактирование своего сообщения (требует JWT)
- **DeleteMessage** - удаление своего сообщения (требует JWT)
- **PinMessage** - закрепление/открепление сообщения (требует JWT, роль admin/owner)

## Технологический стек

- **Go:** 1.24+
- **gRPC + gRPC-Gateway** - API с автоматическим HTTP/JSON маппингом
- **gRPC Streaming** - real-time сообщения через server-side streaming
- **PostgreSQL** - хранение чатов, участников и сообщений
- **Redis** - кэш активных пользователей (синхронизация с UserServer)
- **Kafka** - event-driven обработка событий пользователей
- **OpenTelemetry (OTEL):**
  - **Metrics:** Prometheus → Grafana
  - **Logging:** uber-go/zap → Elasticsearch → Kibana
  - **Tracing:** Jaeger
- **JWT:** Аутентификация через токены из AuthService
- **Envoy Proxy** - маршрутизация и load balancing

## Архитектура

Сервис построен на основе **трёхслойной архитектуры**:

```
┌─────────────────────────────────────┐
│         Handler Layer               │  ← gRPC handlers + HTTP (gRPC-Gateway)
├─────────────────────────────────────┤
│         Service Layer               │  ← Бизнес-логика + Stream Manager
├─────────────────────────────────────┤
│       Repository Layer              │  ← PostgreSQL (chats, messages)
│                                     │    Redis (users cache)
└─────────────────────────────────────┘
```

**Ключевые архитектурные решения:**

- **Dependency Injection (DI) контейнер** - управление зависимостями между слоями
- **Graceful Shutdown** - корректное закрытие активных streaming соединений и БД
- **gRPC-Gateway** - автоматический маппинг HTTP/JSON → gRPC
- **Event-Driven Architecture** - реакция на события пользователей через Kafka

## Структура проекта

```
.
├── cmd/
│   └── server/
│       └── main.go                           # Entry point
├── internal/
│   ├── app/
│   │   ├── app.go                            # Application setup
│   │   └── service_provider.go               # DI container
│   ├── handler/
│   │   └── chat/                             # gRPC handlers
│   ├── service/
│   │   ├── chat/                             # Бизнес-логика
│   │   │   └── stream/                       # Stream manager
│   │   └── consumer/user/                    # Kafka consumers
│   ├── repository/
│   │   ├── chat/                             # Chat repository (PostgreSQL)
│   │   ├── chat_participant/                 # Participant repository
│   │   ├── message/                          # Message repository
│   │   └── converter/                        # Model converters
│   ├── client/
│   │   └── cache/                            # Redis client
│   ├── converter/
│   │   ├── chat.go                           # DTO converters
│   │   └── kafka/decoder/                    # Kafka event decoders
│   ├── middleware/
│   │   └── tokenVerifier.go                  # JWT validation middleware
│   ├── config/                               # Конфигурация из .env
│   ├── errors/domain/                        # Доменные ошибки
│   ├── model/                                # Domain models
│   └── validator/                            # Валидация запросов
├── api/chat/v1/
│   └── chat.proto                            # Protocol Buffers
├── pkg/chat/v1/                              # Сгенерированный Go код
├── migrations/                               # Database migrations
├── docker-compose.yaml
├── Dockerfile
└── Makefile
```

## Конфигурация

### Переменные окружения

```
# ====== Service ======
SERVICE_NAME=chat-server
SERVICE_VERSION=1.0.0

# ====== gRPC ======
GRPC_HOST=0.0.0.0
GRPC_PORT=50052
STREAMING_BUFFER_SIZE=5

# ====== HTTP ======
HTTP_HOST=0.0.0.0
HTTP_PORT=8081

# ====== PostgreSQL ======
PG_HOST=pg-chat
PG_PORT_INNER=5432
PG_PORT_OUTER=5432
PG_DATABASE_NAME=chat
PG_USER=chat-user
PG_PASSWORD=chat-password
PG_TIMEOUT=5s
MIGRATION_DIR=./migrations

# ====== Logger ======
LOGGER_LEVEL=DEBUG
LOGGER_AS_JSON=false
LOGGER_ENABLE_OLTP=true

# ====== OTEL ======
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
OTEL_ENVIRONMENT=development
OTEL_METRICS_PUSH_TIMEOUT=1s

# ====== Rate Limiter ======
RATE_LIMITER_LIMIT=50
RATE_LIMITER_PERIOD=1s

# ====== Circuit Breaker ======
CB_MAX_REQUEST=3
CB_TIMEOUT=5s
CB_FAILURE_RATE=0.5

# ====== Kafka ======
KAFKA_BROKERS=kafka:29092
USER_CREATED_TOPIC_NAME=user.created
USER_DELETED_TOPIC_NAME=user.deleted
USER_CONSUMER_GROUP_ID=user-consumer-group

# ====== Redis (Users Cache) ======
CACHE_HOST=redis-users
CACHE_PORT=6379
CACHE_CONNECTION_TIMEOUT=5s
CACHE_MAX_IDLE=10
CACHE_IDLE_TIMEOUT=300s

# ====== JWT (для валидации токенов) ======
REFRESH_TOKEN_SECRET=refresh_secret
ACCESS_TOKEN_SECRET=access_secret
```

## API

### gRPC Endpoints

Сервис предоставляет следующие gRPC методы (описаны в `api/chat/v1/chat.proto`):

**Управление чатами:**
- **CreateChat** - создание нового чата (требует JWT)
- **DeleteChat** - удаление чата (требует JWT, роль owner)
- **GetChat** - получение информации о чате с сообщениями (требует JWT)
- **ConnectChat** - подключение к чату через streaming (требует JWT)

**Управление участниками:**
- **AddUser** - добавление пользователя в чат (требует JWT, роль admin/owner)
- **RemoveUser** - удаление пользователя из чата (требует JWT, роль admin/owner)
- **UpdateUserRole** - изменение роли участника (требует JWT, роль owner)

**Управление сообщениями:**
- **SendMessage** - отправка сообщения (требует JWT)
- **EditMessage** - редактирование сообщения (требует JWT)
- **DeleteMessage** - удаление сообщения (требует JWT)
- **PinMessage** - закрепление/открепление сообщения (требует JWT, роль admin/owner)

### HTTP Endpoints (gRPC-Gateway)

HTTP запросы автоматически маппятся на gRPC методы через gRPC-Gateway. Маппинг определён в `.proto` файле с помощью аннотаций `google.api.http`.

**Примеры HTTP эндпоинтов:**

**Чаты:**
- `POST /api/v1/chats` → `CreateChat`
- `DELETE /api/v1/chats/{chat_id}` → `DeleteChat`
- `GET /api/v1/chats/{chat_id}` → `GetChat`

**Участники:**
- `POST /api/v1/chats/{chat_id}/add` → `AddUser`
- `DELETE /api/v1/chats/{chat_id}/remove` → `RemoveUser`
- `PATCH /api/v1/chats/{chat_id}/update` → `UpdateUserRole`

**Сообщения:**
- `POST /api/v1/chats/{chat_id}` → `SendMessage`
- `DELETE /api/v1/messages/{message_id}` → `DeleteMessage`
- `PATCH /api/v1/messages/{message_id}` → `EditMessage`
- `PATCH /api/v1/messages/{message_id}/pin` → `PinMessage`

## Роли участников

Система поддерживает три роли участников чата:

| Роль | Права доступа |
|------|---------------|
| **USER** | Отправка, редактирование и удаление своих сообщений |
| **ADMIN** | USER + добавление/удаление участников, закрепление сообщений |
| **OWNER** | ADMIN + удаление чата, изменение ролей участников |

## Real-time Streaming

### ConnectChat (gRPC Server-Side Streaming)

Подключение к чату для получения сообщений в реальном времени через gRPC streaming.

**Процесс:**
1. Клиент вызывает `ConnectChat` с `chat_id`
2. Сервер создаёт stream соединение
3. При отправке нового сообщения в чат (через `SendMessage`), все подключённые клиенты получают его через stream
4. Stream остаётся открытым до отключения клиента или завершения работы сервера

**Параметры:**
- `STREAMING_BUFFER_SIZE=5` - размер буфера сообщений для каждого stream

**Graceful Shutdown:**
- При завершении работы сервера все активные streams корректно закрываются
- Клиенты получают уведомление о закрытии соединения

## Авторизация

Все эндпоинты требуют **JWT токен** для аутентификации.

**Реализация:**
- **Middleware** - JWT validation middleware на уровне gRPC interceptor
- **TokenVerifier** - проверка подписи токена из платформенной библиотеки
- **Claims extraction** - извлечение user_id из токена для проверки прав доступа

**Токены выдаются:** AuthService

## Redis Cache

### Назначение

Redis используется для **кэширования множества ID активных пользователей** из UserServer.

**Цель:** избежать частых запросов к UserServer для проверки существования пользователя.

### Синхронизация через Kafka

**User Created Event:**
- При создании пользователя в UserServer
- Consumer добавляет `user_id` в Redis set

**User Deleted Event:**
- При удалении пользователя в UserServer
- Consumer выполняет каскадное удаление:
  1. Удаление `user_id` из Redis
  2. Закрытие всех активных streaming соединений пользователя
  3. Удаление пользователя из всех чатов
  4. Если пользователь - владелец чата:
     - Удаление чата из БД
     - Закрытие всех streaming соединений к этому чату

## Event-Driven Architecture

### Kafka Consumers

Сервис подписан на события пользователей из UserServer:

| Топик | Событие | Обработка |
|-------|---------|-----------|
| `user.created` | Создание пользователя | Добавление `user_id` в Redis cache |
| `user.deleted` | Удаление пользователя | Каскадное удаление: Redis → Streams → Chats → DB |

**Consumer Group:** `user-consumer-group`

**Каскадное удаление при `user.deleted`:**
1. Удаление из Redis cache
2. Закрытие активных streaming соединений пользователя
3. Удаление из таблицы `chat_participants`
4. Если пользователь - owner чата:
   - Удаление записи из таблицы `chats`
   - Удаление всех сообщений чата из таблицы `messages`
   - Закрытие streaming соединений всех участников чата

## База данных (PostgreSQL)

### Схема

**Таблица `chats`:**
- `id` (BIGSERIAL PRIMARY KEY)
- `owner_id` (BIGINT NOT NULL) - владелец чата
- `name` (VARCHAR)
- `description` (TEXT)
- `created_at` (TIMESTAMP)

**Таблица `chat_participants`:**
- `chat_id` (BIGINT REFERENCES chats)
- `user_id` (BIGINT NOT NULL)
- `role` (ENUM) - USER/ADMIN/OWNER
- `joined_at` (TIMESTAMP)
- PRIMARY KEY (chat_id, user_id)

**Таблица `messages`:**
- `id` (BIGSERIAL PRIMARY KEY)
- `chat_id` (BIGINT REFERENCES chats)
- `sender_id` (BIGINT NOT NULL)
- `text` (TEXT NOT NULL)
- `is_pinned` (BOOLEAN DEFAULT false)
- `send_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## Мониторинг и Observability

### OpenTelemetry (OTEL)

Интегрирован мониторинг из платформенной библиотеки:

**Traces:**
- Распределённая трассировка gRPC методов
- Trace context propagation через metadata

**Metrics:**
- RPS (requests per second)
- Latency percentiles (p50, p95, p99)
- Active streaming connections count

**Logs:**
- Структурированное логирование (uber-go/zap)
- Dual output: stdout + OTEL Collector → Elasticsearch → Kibana

## Отказоустойчивость

### Rate Limiter

- **Реализация:** middleware из платформенной библиотеки
- **Лимит:** 50 запросов/сек
- **Защита:** от DDoS и перегрузки

### Circuit Breaker

Защита при взаимодействии с внешними сервисами (если потребуется в будущем):
- `CB_MAX_REQUEST=3` - макс. запросов в half-open state
- `CB_TIMEOUT=5s` - таймаут перед переходом в half-open
- `CB_FAILURE_RATE=0.5` - порог ошибок (50%)

### Graceful Shutdown

При завершении работы сервиса:
1. Остановка приёма новых запросов
2. Завершение обработки текущих запросов
3. **Закрытие всех активных streaming соединений** с уведомлением клиентов
4. Закрытие соединения с PostgreSQL
5. Закрытие соединения с Redis
6. Остановка Kafka consumers
7. Остановка gRPC сервера

## Сетевое взаимодействие

Все внешние запросы к сервису проходят через **Envoy proxy** (конфигурация в платформенной библиотеке):
- gRPC-JSON transcoding
- Load balancing
- TLS termination
- Rate limiting на уровне proxy

## Обработка ошибок

Кастомные доменные ошибки с маппингом в gRPC статус коды:

| Доменная ошибка | gRPC Code | Описание |
|-----------------|-----------|----------|
| `CHAT_NOT_FOUND` | `NotFound` | Чат не найден |
| `MESSAGE_NOT_FOUND` | `NotFound` | Сообщение не найдено |
| `PERMISSION_DENIED` | `PermissionDenied` | Недостаточно прав (роль) |
| `USER_NOT_IN_CHAT` | `PermissionDenied` | Пользователь не является участником |
| `INVALID_ROLE` | `InvalidArgument` | Некорректная роль |
| `UNAUTHORIZED` | `Unauthenticated` | Невалидный JWT токен |

## Зависимости

### Внешние сервисы

- **AuthService** - выдача JWT токенов
- **UserServer** - источник событий о пользователях
- **PostgreSQL** - хранение чатов и сообщений
- **Redis** - кэш активных пользователей
- **Kafka** - event streaming
- **OTEL Collector** - мониторинг
- **Envoy Proxy** - маршрутизация

### Платформенная библиотека

Переиспользуемые компоненты:
- JWT Verifier (валидация токенов)
- OTEL instrumentation (traces, metrics, logs)
- Rate Limiter middleware
- Circuit Breaker
- Kafka wrapper (producer/consumer)
- Envoy configuration
