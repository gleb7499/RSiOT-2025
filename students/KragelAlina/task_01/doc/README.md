# Node/Express Service with Postgres

## Description

A minimal Node/Express HTTP service running on port 8072, with a /ping health endpoint. Connects to Postgres for data persistence. Built with security in mind: non-root user (UID 10001), multi-stage Docker build, graceful shutdown on SIGTERM.

Version: v10

## Метаданные студента

- ФИО: Крагель Алина Максимовна
- Группа: АС-63
- № студенческого (StudentID): 220046
- Email (учебный): AS006417@g.bstu.by
- GitHub username: Alina529
- Вариант №: 10
- Дата выполнения: October 1, 2025
- ОС (версия): сборка ОС 19045.6332
- Версия Docker Desktop/Engine: 4.46.0 (204649)/28.4.0

## Setup and Run

1. Ensure Docker and Docker Compose are installed.
2. Build and start: `docker-compose up --build`
3. Access: http://localhost:8072/ (hello message) or http://localhost:8072/ping (health check)
4. Environment vars: Configurable in docker-compose.yml (e.g., DB creds).
5. Stop: Ctrl+C (triggers graceful shutdown) or `docker-compose down`

## Image Tag

Uses `node-express-service:stu-220046-v10`.

## Volumes

- `data_v10`: Persists Postgres data.

## Testing Graceful Shutdown

- Run the service.
- Send SIGTERM: `docker kill --signal=SIGTERM app-AC-63-220046-v10`
- Check logs for graceful close (HTTP server and DB pool).

## Build Optimization

Dockerfile caches npm dependencies by copying package.json first.

## Логи старта, обработки запросов и корректного shutdown

### Пример логов старта (из docker logs app-AC-63-220046-v10)

```
STU_ID: 220046, STU_GROUP: АС-63, STU_VARIANT: 10
Server running on port 8072
Connected to Postgres
```

### Пример логов обработки запросов

```
GET / 200 - - 1.234 ms
GET /ping 200 - - 5.678 ms
```

### Пример логов shutdown (после SIGTERM)

```
SIGTERM received. Shutting down gracefully...
HTTP server closed.
DB pool closed.
```
