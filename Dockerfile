# Используем официальный образ Go
FROM golang:latest AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum (если есть) в контейнер
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main .

# Создаем новый образ для сервера
FROM nginx:alpine

# Копируем статические файлы React в Nginx
COPY --from=builder /app/frontend/build /usr/share/nginx/html

# Копируем исполняемый файл Go в Nginx
COPY --from=builder /app/main /usr/bin/main

# Экспортируем порт 80
EXPOSE 80

# Запускаем Nginx и Go-сервер
CMD ["sh", "-c", "/usr/bin/main & nginx -g 'daemon off;'"]