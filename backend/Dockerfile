# Используем официальный образ Go
FROM golang:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o ./bin/main ./cmd/main.go

# Экспортируем порт 8080
EXPOSE 8082

# Запускаем приложение
CMD ["./bin/main"]