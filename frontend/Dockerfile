# Используем облегченный образ Node.js для сборки React-приложения
FROM node:alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы package.json и package-lock.json
COPY ./package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем исходный код
COPY ./ ./

# Собираем приложение
RUN npm run build

# Используем официальный образ Nginx
FROM nginx:alpine

# Копируем статические файлы React в Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Экспортируем порт 80
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]