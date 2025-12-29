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
COPY iam/go.mod iam/go.sum ./iam/

# Отключаем go.work (он ссылается на все модули, а нам нужен только iam)
ENV GOWORK=off

# Загружаем все зависимости, указанные в модульных файлах
RUN go mod download -C iam

# Копируем исходный код всех модулей — включая зависимости, так как они нужны при компиляции
COPY platform ./platform
COPY shared ./shared
COPY iam ./iam

# Скачиваем grpc-health-probe — это утилита для проверки состояния gRPC-сервиса
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

# Делаем скачанный файл исполняемым
RUN chmod +x grpc-health-probe

# Собираем бинарный файл IAM сервиса для Linux-архитектуры без CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C iam -o /app/app-iam ./cmd/grpc_server/main.go


# ============================
# Stage 2: Final image
# ============================

# Используем чистый Alpine образ как итоговый контейнер — он будет лёгким и безопасным
FROM alpine:3.21.3

# Создаём системного пользователя без root-прав — best practice для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app/app-iam .

# Копируем бинарь grpc-health-probe в /bin
COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

# Копируем директорию миграций
COPY iam/migrations ./iam/migrations

# Запускаем приложение под non-root пользователем
USER appuser

# Экспонируем порт gRPC сервиса IAM
EXPOSE 50053

# Устанавливаем команду запуска — запускаем наш gRPC сервер
ENTRYPOINT ["./app-iam"]

