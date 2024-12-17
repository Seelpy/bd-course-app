# Используем образ Go для сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только файлы с зависимостями для кэширования
COPY go.mod go.sum ./

# Устанавливаем зависимости и инструменты
RUN go mod download && \
    go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Копируем исходный код
COPY . .

# Генерируем код из OpenAPI спецификации
RUN /go/bin/oapi-codegen -generate types,server,spec -package api ./api/api.yaml > ./api/api.gen.go

# Упрощаем и объединяем go мод команды
RUN go mod tidy && \
    go mod vendor

# Компиляция приложения
RUN go build -o ./bin/main ./cmd/main.go

# Финальный образ с минимальным размером
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем необходимые файлы из билд-образа
COPY --from=builder /app/bin/main ./bin/main
COPY --from=builder /app/vendor ./vendor
COPY --from=builder /app/api/api.yaml ./api/api.yaml
COPY --from=builder /app/api/api.gen.go ./api/api.gen.go

# Экспортируем порт
EXPOSE 8082

# Запускаем приложение
CMD ["./bin/main"]
