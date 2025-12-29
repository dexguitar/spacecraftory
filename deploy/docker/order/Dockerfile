# ============================
# Stage 1: Build stage
# ============================

# Используем официальный образ Go на базе Alpine — лёгкий, быстрый и безопасный
FROM golang:1.25.2-alpine AS builder

# Устанавливаем git — нужен для загрузки зависимостей из приватных и публичных репозиториев
RUN apk add --no-cache git

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.work и его checksum — это позволяет Go правильно связать модули
COPY go.work go.work.sum ./

# Копируем go.mod и go.sum всех зависимых модулей — нужно для кэширования go mod download
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY order/go.mod order/go.sum ./order/

# Отключаем go.work (он ссылается на все модули, а нам нужен только order)
ENV GOWORK=off

# Загружаем все зависимости, указанные в модульных файлах
RUN go mod download -C order

# Копируем исходный код всех модулей — включая зависимости, так как они нужны при компиляции
COPY platform ./platform
COPY shared ./shared
COPY order ./order

# Собираем бинарный файл Order сервиса для Linux-архитектуры без CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C order -o /app/app-order ./cmd/http_server/main.go


# ============================
# Stage 2: Final image
# ============================

# Используем чистый Alpine образ как итоговый контейнер — он будет лёгким и безопасным
FROM alpine:3.21.3

# Устанавливаем curl для health checks
RUN apk add --no-cache curl

# Создаём системного пользователя без root-прав — best practice для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app/app-order .

# Копируем директорию миграций
COPY order/migrations ./order/migrations

# Запускаем приложение под non-root пользователем
USER appuser

# Экспонируем порт HTTP-сервиса Order
EXPOSE 8080

# Устанавливаем команду запуска — запускаем наш бинарь
ENTRYPOINT ["./app-order"]

