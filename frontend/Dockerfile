# Используем официальный образ Nginx
FROM nginx:alpine

# Копируем статические файлы
COPY ./dist /usr/share/nginx/html

# Экспортируем порт 80
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]