# ЛР01 — Контейнеризация и Docker (вариант 19)

Кратко: Node/Express приложение с зависимостью Redis, упакованное в многоступенчатый Docker-образ и запускаемое через `docker-compose`.

## Метаданные студента

- ФИО: Соколова Маргарита Александровна
- Группа: АС-63
- StudentID: 220024
- Email: `as006321@g.bstu.by`
- GitHub username: Ritkas33395553
- Вариант: 19
- Дата выполнения: 15.10.2025
- ОС: Windows 10.0.19044.3086
- Docker Desktop/Engine: 28.4.0

## Короткое описание проекта

- Стек: Node.js (Express)
- Внутренний порт приложения: 8053
- Health (liveness): `GET /healthz` → 200 OK
- Readiness: `GET /ready` → проверяет доступность Redis, возвращает 200 OK если готов
- Зависимость: Redis (в `docker-compose.yml`)
- Тома: `data_v19` (используется как volume для приложения и Redis)
- Приложение реализует корректный graceful shutdown (обработка SIGTERM/SIGINT, закрытие HTTP сервера и Redis-клиента)

Файлы, которые важны:

- `src/src/Dockerfile` — multi-stage Dockerfile (USER 65532, LABELs с метаданными)
- `src/src/server.js` — исходный код сервера (Express, endpoints и graceful shutdown)
- `src/docker-compose.yml` — Compose-файл для локального запуска (сервис `app` и `redis`)
- `src/src/package.json` — зависимости (express, redis)

## Именование и метаданные (соответствие требованиям ЛР)

- IMAGE TAG: ritkas33395553/lr01-node-v19:v19 (как указано в `docker-compose.yml`)
- SLUG: `as-63-220024-v19` (используется в labels и именовании ресурсов)
- Имена ресурсов в `docker-compose.yml`:
  - Контейнер приложения: `node-app-220024-v19`
  - Контейнер Redis: `redis-v19`
  - Тома: `data_v19`
  - Labels на сервисах: `org.bstu.owner: "Ritkas33395553"`, `org.bstu.student.slug: "as-63-220024-v19"`

IMAGE LABELS в `Dockerfile` (обязательные):

- org.bstu.student.fullname="Соколова Маргарита Александровна"
- org.bstu.student.id="220024"
- org.bstu.group="АС-63"
- org.bstu.variant="19"
- org.bstu.course="RSIOT"

ENV-переменные, используемые приложением:

- `PORT` — порт приложения (по умолчанию 8053)
- `REDIS_HOST`, `REDIS_PORT` — параметры подключения к Redis
- Рекомендуется добавить (по заданию) ENV: `STU_ID`, `STU_GROUP`, `STU_VARIANT` при необходимости логирования/конфигурации

## Быстрый старт (PowerShell, Windows)

1. Перейдите в корень задания:

```powershell
Set-Location -Path 'd:\Предметы и литература\РСОТ\RSiOT-2025\students\SokolovaMargarita\task_01'
```

1. Сборка и запуск через docker-compose (встроенная сборка образа по Dockerfile в `src`):

```powershell
docker compose -f .\src\docker-compose.yml up --build -d
```

1. Просмотр логов сервиса (в отдельном окне/терминале):

```powershell
docker compose -f .\src\docker-compose.yml logs -f app
```

1. Проверка health и readiness (локально, с хоста):

```powershell
# liveness
curl http://127.0.0.1:8053/healthz

# readiness (проверяет Redis)
curl http://127.0.0.1:8053/ready
```

Ожидаемый ответ для `/healthz`: `OK` (HTTP 200)

Ожидаемый ответ для `/ready`: `READY` (HTTP 200) — если Redis доступен.

## Проверка graceful shutdown

1. Отправьте SIGTERM контейнеру приложения (Compose):

```powershell
docker kill --signal=SIGTERM node-app-220024-v19
```

1. Откройте логи и убедитесь, что приложение выведет строки, подобные:

- "SIGTERM/SIGINT получен — начинаем корректный shutdown..."
- "Redis корректно завершён"
- "HTTP сервер остановлен"

Если контейнер корректно завершился на 0 и в логах показан успешный shutdown — требование выполнено.

## Как собрать образ локально вручную (альтернативный путь)

```powershell
# Перейти в директорию с Dockerfile
Set-Location -Path .\src\src

# Собрать образ с нужными метками/тегом (пример)
docker build -f Dockerfile -t ritkas33395553/lr01-node-v19:stu-220024-v19 .
```

Примечание: в Dockerfile уже есть LABELs с метаданными студента.

## Тонкости реализации (соответствие требованиям задания)

- Multi-stage Dockerfile: `src/src/Dockerfile` использует multi-stage сборку и кеширование зависимостей (слой `deps`).
- Non-root: процесс запускается от `USER 65532:65532`.
- HEALTHCHECK: описан в Dockerfile, проверяет `http://127.0.0.1:8053/healthz`.
- Конфигурация через ENV: `PORT`, `REDIS_HOST`, `REDIS_PORT`. При желании добавьте `STU_ID`, `STU_GROUP`, `STU_VARIANT` в `docker-compose.yml`.
- Volume: `data_v19` используется для хранения данных и поделен между сервисами в compose.

## Список команд для отладки и верификации

- Просмотр статуса контейнеров:

```powershell
docker compose -f .\src\docker-compose.yml ps
```

- Просмотр логов всех сервисов:

```powershell
docker compose -f .\src\docker-compose.yml logs --no-color
```

- Остановить и удалить контейнеры/тома:

```powershell
docker compose -f .\src\docker-compose.yml down -v
```

## Замечания по оформлению и сдаче

- Убедитесь, что в `Dockerfile` и `docker-compose.yml` актуальные значения `org.bstu.*` labels и `org.bstu.student.slug` совпадают с README.
- Перед отправкой PR замените плейсхолдер Email и, при необходимости, ОС/версии Docker в разделе Метаданные студента.

## Файлы проекта

- `src/src/Dockerfile` — Dockerfile
- `src/docker-compose.yml` — Compose
- `src/src/server.js` — приложение
- `src/src/package.json` — package manifest

---

Если хотите, я могу:

- дополнительно добавить в `docker-compose.yml` переменные окружения STU_ID/STU_GROUP/STU_VARIANT;
- сгенерировать `.dockerignore` или улучшить Dockerfile для уменьшения размера образа;
- собрать и запустить окружение у себя (только если разрешите запуск команд локально).

Готов продолжать — скажите, что делать дальше.
