# РСиОТ Лабораторная работа №1

## Описание

HTTP-сервис на Go с эндпоинтами `/`, `/health`, `/hit`, использующий Redis для хранения счетчика посещений.
Проект упакован в Docker-контейнеры с использованием Docker Compose для развертывания приложения и Redis.
Сервер поддерживает graceful shutdown и health checks.

## Структура проекта

- `cmd/server/` - исходный код сервера (main.go)
- `scripts/` - скрипты для тестирования (smoke.sh)
- `Dockerfile` - для сборки образа приложения
- `docker-compose.yml` - конфигурация для запуска проекта
- `.env` - переменные окружения
- `go.mod`, `go.sum` - зависимости Go
- `.gitignore` - исключения для Git

## Требования

- Docker
- Docker Compose

## Как запустить

Забилдить и запустить используя Docker Compose:

```bash
docker compose up --build
```

Сервер будет доступен по адресу `http://localhost:8041`

Для тестирования можно выполнить скрипт:

```bash
./scripts/smoke.sh
```

## Endpoints

- `GET /` - основная информация о приложении и студенте
- `GET /health` - проверка работоспособности Redis
- `GET /hit` - увеличение счетчика посещений в Redis

## Graceful Shutdown

Сервер поддерживает graceful shutdown при получении SIGTERM или SIGINT сигналов.
При завершении работы сервер корректно закрывает HTTP-соединения и Redis-клиент в течение 10 секунд.

## Особенности реализации

- Использует Redis для хранения счетчика с префиксом `stu:220026:v21:hits`
- Health check проверяет доступность Redis
- Мультиконтейнерная архитектура (app + redis)
- Переменные окружения для конфигурации
- Логирование ключевых событий работы приложения

## Student Metadata

- **Full Name:** Tunchik Anton Dmitrievich
- **Group:** AS-63
- **Student ID:** 220026
- **Email (Academic):** AS63006326@g.bstu.by
- **GitHub Username:** Stis25
- **Variant №:** 21
- **Operating System:** Windows 11 23H2
- **Docker Version:** Docker Desktop 4.47.0
