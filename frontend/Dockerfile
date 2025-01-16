# Используем официальный образ Nginx
FROM nginx:alpine

# Копируем статические файлы
COPY ./dist /usr/share/nginx/html

# Копируем пользовательскую конфигурацию Nginx
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

# Экспортируем порт 80
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]