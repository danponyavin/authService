**Задание:**

Написать часть сервиса аутентификации.

**Запуск приложения:**

1. Запуск docker-compose:

docker-compose up

2. Применение миграций:

migrate -path ./migrations -database postgres://postgres:mypass@localhost:5440/test?sslmode=disable up

**Документация приложения:**

http://localhost:8000/swagger/index.html

**Запуск тестов:**

go test -v ./pkg/handler/...
