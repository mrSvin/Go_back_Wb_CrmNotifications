# Используем официальный образ Go для сборки
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM debian:stable-slim

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение из стадии builder
COPY --from=builder /app/myapp .

# Открываем порт, на котором будет работать приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["./myapp"]