# Лабораторная работа №1 — Контейнеризация и Docker (Вариант 14)

## 1. Цель работы
Закрепить навыки контейнеризации приложений с использованием Docker: упаковка Flask-сервиса, настройка healthcheck, graceful shutdown, организация multi-stage сборки, работа с docker compose, использование переменных окружения и метаданных.

## 2. Шаги сборки и запуска

### Предварительно
Требуется установленный **Docker Desktop v4.45.0** (или совместимый движок) на **Windows 11 24H2**.

### Быстрый старт с Make (рекомендуемый способ)
```powershell
# Просмотр всех доступных команд
make help

# Сборка и запуск
make build
make up

# Проверка здоровья приложения
make health

# Просмотр логов
make logs

# Остановка приложения
make stop

# Полная остановка и удаление сервисов
make down
```

### Альтернативно через Docker команды
```powershell
# Сборка образа вручную с требуемым тегом варианта/студента
docker build -t gleb7499/lab1-v14:stu-220018-v14 .

# Запуск через docker compose
docker compose up -d --build

# Просмотр логов приложения
docker compose logs -f app

# Остановка
docker compose down
```

## 3. Пример логов работы
```
2025-09-16 10:15:04,120 | INFO | ==== Application Startup ==== 
2025-09-16 10:15:04,121 | INFO | Student ID: 220018
2025-09-16 10:15:04,121 | INFO | Student Group: АС-63
2025-09-16 10:15:04,121 | INFO | Student Variant: 14
2025-09-16 10:15:04,122 | INFO | ENV STU_ID=220018
2025-09-16 10:15:04,122 | INFO | ENV STU_GROUP=АС-63
2025-09-16 10:15:04,122 | INFO | ENV STU_VARIANT=14
2025-09-16 10:15:04,123 | INFO | ================================================
2025-09-16 10:15:04,123 | INFO | Starting Flask server on 0.0.0.0:8062
# Запрос на /healthz
2025-09-16 10:15:10,500 | INFO | 200 GET /healthz
# Завершение (docker stop / SIGTERM)
2025-09-16 10:15:20,010 | WARNING | Received signal 15 - initiating graceful shutdown...
2025-09-16 10:15:20,311 | INFO | Stop accepting new connections. Shutdown flag set.
2025-09-16 10:15:20,512 | INFO | Graceful shutdown complete.
```

## 4. Пример запроса к `/healthz`
```powershell
# С помощью Make (рекомендуемый способ)
make health

# Или вручную через PowerShell
Invoke-RestMethod -Uri http://localhost:8062/healthz -Method GET

# Или с помощью curl
curl http://localhost:8062/healthz
```
Пример ответа:
```json
{
  "status": "ok",
  "timestamp": "2025-09-16T07:15:10.500Z"
}
```

## 5. Структура репозитория
```
.
├── Dockerfile
├── docker-compose.yml
├── requirements.txt
├── .dockerignore
├── Makefile
├── README.md
└── src/
    └── app.py
```

## 6. Описание компонентов
- **Flask приложение** (`src/app.py`): маршруты `/` (информация), `/healthz` (health check), `/echo` (POST echo). Логирует учебные ENV и корректно завершается по SIGTERM/SIGINT (graceful shutdown).
- **Dockerfile**: multi-stage (builder + final), Alpine, непривилегированный пользователь UID 10001, HEALTHCHECK, LABELS (учебные метаданные), pip без кэша.
- **docker-compose.yml**: сервисы `db` (Postgres 16 Alpine) и `app`; именование по шаблону (app-as63-220018-v14, db-as63-220018-v14, сеть net-as63-220018-v14, том data-as63-220018-v14), метки (labels) на обоих сервисах.


## 7. Метаданные студента (обязательно)

В отчёте (README в корне проекта) укажите:

- **ФИО (полностью):** Логинов Глеб Олегович
- **Группа:** АС-63
- **№ студенческого (StudentID):** 220018
- **Email (учебный):** as006315@g.bstu.by
- **GitHub username:** gleb7499
- **Вариант №:** 14
- **Дата выполнения:** 08.10.2025
- **ОС (версия), версия Docker Desktop/Engine:** Windows 11 24H2, Docker Desktop v4.45.0

В артефактах:

- **Dockerfile → LABEL:**
  - org.bstu.student.fullname = Логинов Глеб Олегович
  - org.bstu.student.id = 220018
  - org.bstu.group = АС-63
  - org.bstu.variant = 14
  - org.bstu.course = RSIOT
- **docker-compose.yml → labels на сервисах:**
  - org.bstu.owner = gleb7499
  - org.bstu.student.slug = as63-220018-v14
- **slug** = as63-220018-v14

## 8. Дополнительно
- Образ создаётся через multi-stage: зависимости устанавливаются отдельно и копируются в финальный слой для оптимизации размера (<150MB).
- Используется `HEALTHCHECK`, чтобы оркестратор мог отслеживать состояние контейнера.
- `ENTRYPOINT` в exec-форме обеспечивает доставку сигналов приложению (graceful shutdown).
- Логи пишутся в stdout для интеграции с Docker logging driver.

## 9. Команды для быстрого теста

### С помощью Make
```powershell
# Проверка состояния контейнеров
make ps

# Проверка размера образа
make image-size

# Тест graceful shutdown
make graceful-test

# Открыть shell в контейнере
make shell

# Очистка Docker ресурсов
make clean
```

### Прямые Docker команды
```powershell
# Проверка состояния контейнеров
docker ps

# Просмотр health статуса
docker inspect --format='{{json .State.Health}}' app-as63-220018-v14

# Отправка POST на /echo
curl -X POST http://localhost:8062/echo -H "Content-Type: application/json" -d '{"msg":"hello"}'
```

## 10. Имя проекта Docker Compose

Иногда при запуске `make up` в директории с пробелами или не-ASCII символами (например, кириллицей) Docker Compose может вывести ошибку:

```
project name must not be empty
```

Чтобы избежать этой проблемы, в `Makefile` добавлена переменная `COMPOSE_PROJECT_NAME`, по умолчанию равная `as63-220018-v14`. Все команды `docker compose` теперь вызываются с явным указанием `--project-name`. Это гарантирует стабильные имена контейнеров, сети и volume независимо от пути к каталогу.

Переопределить имя проекта можно так:

```powershell
COMPOSE_PROJECT_NAME=myexp make up

# или для PowerShell (переменная только на время команды)
$env:COMPOSE_PROJECT_NAME='testproj'; make up
```

Проверить, что используется нужное имя:

```powershell
make ps
```

Если нужно полностью удалить созданные ресурсы (контейнеры/сеть/том) для альтернативного имени проекта — укажите то же имя и выполните:

```powershell
COMPOSE_PROJECT_NAME=testproj make down
```

Причина ошибки: Docker пытается вычислить имя проекта из имени директории, но в редких случаях (особенно Windows + локали) это может вернуть пустое или некорректное значение. Явное указание имени делает процесс детерминированным.
