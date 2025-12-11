# =====================================================
# AI Browser Agent - Production Dockerfile
# =====================================================
# Multi-stage build для оптимизации размера образа
# =====================================================

# ========== STAGE 1: Builder ==========
FROM golang:1.24-alpine AS builder

# Установка необходимых зависимостей для сборки
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

WORKDIR /build

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Компилируем приложение
# CGO_ENABLED=0 для статической линковки
# -ldflags="-s -w" для уменьшения размера бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty)" \
    -o /build/agent \
    ./cmd/agent

# ========== STAGE 2: Runtime ==========
FROM ubuntu:22.04

# Метаданные образа
LABEL maintainer="your-email@example.com"
LABEL description="AI Browser Agent - Автономный агент для управления браузером"
LABEL version="1.0"

# Устанавливаем переменные окружения
ENV DEBIAN_FRONTEND=noninteractive \
    LANG=C.UTF-8 \
    TZ=UTC

# Устанавливаем системные зависимости и Chrome
RUN apt-get update && apt-get install -y \
    # Основные утилиты
    ca-certificates \
    tzdata \
    curl \
    wget \
    # Зависимости для Chrome/Chromium
    chromium-browser \
    chromium-chromedriver \
    # Дополнительные библиотеки для GUI
    xvfb \
    x11vnc \
    fluxbox \
    # Шрифты для корректного отображения
    fonts-liberation \
    fonts-noto-color-emoji \
    # Очистка
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Создаём непривилегированного пользователя
RUN groupadd -r agent && useradd -r -g agent -u 1001 agent

# Создаём рабочую директорию
WORKDIR /app

# Копируем скомпилированный бинарник из builder stage
COPY --from=builder /build/agent /app/agent

# Создаём необходимые директории с правильными правами
RUN mkdir -p \
    /app/logs \
    /app/screenshots \
    /app/.browser-data \
    && chown -R agent:agent /app

# Переключаемся на непривилегированного пользователя
USER agent

# Устанавливаем Chrome как переменную окружения для Rod
ENV CHROME_BIN=/usr/bin/chromium-browser

# Healthcheck (опционально, для мониторинга)
# HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
#   CMD pgrep -x agent || exit 1

# Точка входа
ENTRYPOINT ["/app/agent"]

# Команда по умолчанию
CMD ["run"]

# =====================================================
# Volumes (рекомендуется монтировать при запуске):
# - /app/logs           - логи приложения
# - /app/screenshots    - скриншоты
# - /app/.browser-data  - persistent browser sessions
# =====================================================
