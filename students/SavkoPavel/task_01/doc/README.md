# Министерство образования Республики Беларусь

<p align="center">Учреждение образования</p>
<p align="center">“Брестский Государственный технический университет”</p>
<p align="center">Кафедра ИИТ</p>

<p align="center"><strong>Лабораторная работа №1</strong></p>
<p align="center"><strong>По дисциплине:</strong> “РСиОТ”</p>
<p align="center"><strong>Тема:</strong> “Контейнеризация и Docker”</p>

<p align="right"><strong>Выполнил:</strong></p>
<p align="right">Студент 4 курса</p>
<p align="right">Группы АС-63</p>
<p align="right">Савко П.С.</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А.Н.</p>

<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Освоить **основы контейнеризации с использованием Docker**, научиться собирать **минимальные образы (multi-stage build)**, работать с **volume**, **сетями**, и **docker-compose**, а также реализовать **graceful shutdown** и **healthcheck** для сервисов.

---

### Вариант №18

| Параметр            | Значение                  |
|---------------------|---------------------------|
| Стек                | Go (net/http)             |
| Порт приложения     | 8052                      |
| Healthcheck endpoint| `/live`                   |
| Зависимость         | PostgreSQL                |
| Volume              | `data-as-63-220023-v18`   |
| UID                 | 10001                     |
| Тег                 | `v18`                     |

---

## Ход выполнения работы

### 1. Структура проекта

- `Dockerfile` — multi-stage сборка Go-приложения
- `docker-compose.yml` — оркестрация приложения и базы данных
- `cmd/app/` — исходный код HTTP-сервиса (Go net/http)
- `README.md` — отчёт и инструкции
- `logs/` — примеры логов старта и остановки контейнера

---

### 2. Dockerfile (основные моменты)

- Используется **multi-stage build**:
  `golang:1.23-alpine` → `gcr.io/distroless/static-debian12:nonroot`
- Установлен **USER 10001** (ненулевой UID)
- Минимальный размер финального образа ≤ 150 MB
- Настроено кэширование зависимостей:

  ```Dockerfile

  --mount=type=cache,target=/go/pkg/mod
  --mount=type=cache,target=/root/.cache/go-build

  ```

- Реализован HEALTHCHECK:

  ```Dockerfile

  HEALTHCHECK --interval=10s --timeout=2s --retries=5     CMD ["/app/app", "-test.healthcheck"]

  ```

- Переменные окружения:

  ```

  STU_ID=220023
  STU_GROUP=АС-63
  STU_VARIANT=18
  PGDATABASE=app_220023_v18

  ```

- LABEL’ы с метаданными студента:

  ```

  org.bstu.student.fullname="Савко Павел Станиславович"
  org.bstu.student.id="220023"
  org.bstu.group="АС-63"
  org.bstu.variant="18"
  org.bstu.course="RSIOT"

  ```

---

### 3. docker-compose.yml

Включает два сервиса:
- **app** — основное приложение
- **db** — PostgreSQL 16 (alpine)

Особенности:
- Связь через сеть `net-as-63-220023-v18`
- Volume для хранения данных: `data-as-63-220023-v18`
- HEALTHCHECK `/live`
- Переменные окружения передаются через ENV
- Метки:

  ```

  org.bstu.owner = "1nsirius"
  org.bstu.student.slug = "as-63-220023-v18"

  ```

---

### 4. Реализация graceful shutdown

В коде Go-приложения реализована обработка сигналов:

```go

sig := make(chan os.Signal, 1)
signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
<-sig
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
server.Shutdown(ctx)

```

При остановке контейнера (`docker compose down`) сервис корректно завершает работу, закрывая соединение с БД и завершая все активные запросы.

---

### 5. Проверка работы

#### Запуск

```bash

docker compose up --build

```

#### Логи старта

```

[INFO] Starting app on 0.0.0.0:8052
[INFO] Connected to Postgres db=app_220023_v18
[INFO] Health endpoint: /live

```

#### Проверка healthcheck

```

$ curl http://localhost:8052/live
OK

```

#### Корректное завершение

```

[INFO] Received SIGTERM
[INFO] Shutting down HTTP server gracefully...
[INFO] Shutdown complete.

```

---

### 6. Проверка требований

| Критерий | Выполнено | Комментарий |
|-----------|-----------|--------------|
| Multi-stage build | ✅ | минимальный образ |
| USER ненулевой | ✅ | UID=10001 |
| EXPOSE / HEALTHCHECK | ✅ | порт 8052, /live |
| ENV конфигурация | ✅ | все переменные заданы |
| docker-compose (app + db + volume + network) | ✅ | корректно оформлено |
| Graceful shutdown | ✅ | через `server.Shutdown` |
| Кэширование зависимостей | ✅ | с `--mount=type=cache` |
| LABEL / slug / теги образов | ✅ | все поля присутствуют |
| README и отчёт | ✅ | оформлено в соответствии с ТЗ |

---

## Метаданные студента

| Поле | Значение |
|------|-----------|
| **ФИО** | Савко Павел Станиславович |
| **Группа** | АС-63 |
| **StudentID** | 220023 |
| **Email (учебный)** | 220023@bsut.by |
| **GitHub username** | 1nsirius |
| **Вариант** | 18 |
| **ОС / Docker** | W10, Docker Desktop 4.45.0 |
| **Slug** | as-63-220023-v18 |

---

## Вывод

В ходе лабораторной работы:
- Собран **минимальный контейнеризированный Go-сервис** (≤150 MB)
- Настроен **PostgreSQL** как зависимость
- Реализован **graceful shutdown** и **healthcheck**
- Использованы **multi-stage build** и **кэширование зависимостей**
- Все требования по **метаданным**, **именованию**, **volume** и **сетям** выполнены

Работа продемонстрировала навыки практической контейнеризации и сборки production-ready образов с использованием Docker и Docker Compose.
