# Инструкция по использованию

## Как запустить контейнер

### Запускаем сервер с фронтом

`docker-compose up --build`

## Копируем генеренную API
`docker cp db_server:app/api/api.gen.go ./backend/api/api.gen.go`

## Vendor
`docker cp db_server:app/vendor ./backend/vendor`

### Расположение систем

**frontend** - репозиторий фронта
**backend** - репозиторий  для бэка