# РСиОТ Лабораторная работа №1

## Описание

Это простой HTTP-сервис написанный на языке GoLang, с определенными endpoint-ами (/ready, /health) и поддержкой graceful shutdown. Данный проект был упакован в Docker-контейнер с использованием мельтиуровневого развёртывания для более удобного и работоспособного деплоя на другие устройства. Сервер имеет единственную зависимость - драйвера для работы с PostgreSQL.

## Структура проекта

src/ - Содержит исходный код сервера (server.go), а также зависимости.
Dockerfile - Мультуровневный Dockerfile для развертывания проекта.
.dockerignore - Исключает необязатльные файлы из создания контейнеров.
docker-compose.yml - Конфигурация для локального запуска проекта.
README.md - Документация проекта.

## Требования

- Docker
- Docker Compose

## Как запустить

Забилдить и запустить используя Docker Compose:
```docker-compose up --build```

Сервер будет доступен по адресу http://localhost:8092.

Чтобы вручную построить Docker-image:

```
docker build -t go-server .
docker run -p 8092:8092 go-server
```

## Endpoints

GET /ready - Readiness check
GET /health - Health check

## Graceful Shutdown

Сервер поддерживает graceful shutdown, ожидая сигнала типа SIGNIT или SIGTERM, позволяя системе закончить все запросы в течении 5 секунд. Для теста, пошлите какой-либо SIGTERM-сигнал (например Ctrl+C в консоли) и ожидайте "Shutting down server..." и "Server exiting" в коносли.

## Student Metadata

```
Full Name: Карпеш Никита Петрович
Group: АС-63
Student ID: 220009
Email (Academic): AS63006306@g.bstu.by
GitHub Username: Frosyka
Variant №: 6
Completion Date: 15/09/2025
Operating System: Windows 10 Pro 22H2, Ubuntu 22.04
Docker Version: Docker Desktop 4.45.0 / Engine 28.3.3
```
