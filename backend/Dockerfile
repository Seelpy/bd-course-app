# Используем официальный образ Go
FROM golang:latest AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Устанавливаем oapi-codegen
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Копируем исходный код, включая директорию backend/api
COPY . .

# Создаем директорию для сгенерированных файлов, если она не существует
RUN mkdir -p ./api

# Генерируем код из OpenAPI спецификации
RUN /go/bin/oapi-codegen -generate types,server,spec -package api ./api/api.yaml > ./api/api.gen.go

RUN go mod tidy

RUN go mod download

RUN go mod vendor

# Собираем приложение
RUN go build -o ./bin/main ./cmd/main.go

# Используем другой образ для финального результата
FROM golang:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем сгенерированные файлы из builder
COPY --from=builder /app/api/api.yaml ./api/api.yaml
COPY --from=builder /app/api/api.gen.go ./api/api.gen.go
COPY --from=builder /app/bin/main ./bin/main
COPY --from=builder /app/vendor ./vendor

# Экспортируем порт 8082
EXPOSE 8082

# Запускаем приложение
CMD ["./bin/main"]