# Лабораторная работа №1

<p align="center">Министерство образования Республики Беларусь</p>
<p align="center">Учреждение образования</p>
<p align="center">"Брестский Государственный технический университет"</p>
<p align="center">Кафедра ИИТ</p>
<br><br><br><br><br><br>
<p align="center"><strong>Лабораторная работа №1</strong></p>
<p align="center"><strong>По дисциплине:</strong> "Распределенные системы и облачные технологии"</p>
<p align="center"><strong>Тема:</strong> "Контейнеризация и Docker"</p>
<br><br><br><br><br><br>
<p align="right"><strong>Выполнил:</strong></p>
<p align="right">Студент 4 курса</p>
<p align="right">Группы АС-63</p>
<p align="right">Козловская А.Г.</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А.Н.</p>
<br><br><br><br><br>
<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Научиться собирать минимальные образы (multi-stage) и запускать контейнеры под непривилегированным пользователем, закрепить основы docker-compose (зависимости БД/кэш, volume, сети), настроить healthcheck и graceful shutdown.

---

### Вариант №8

## Метаданные студента

- **ФИО:** Козловская Анна Геннадьевна
- **Группа:** АС-63
- **№ студенческого (StudentID):** 220012
- **Email (учебный):** <as006309@g.bstu.by>
- **GitHub username:** annkrq
- **Вариант №:** 8
- **ОС и версия:** Windows 10 1809, Docker Desktop v4.53.0

---

## Окружение и инструменты

**Технологический стек:**

- **Язык программирования:** Python 3.11
- **Веб-фреймворк:** Flask
- **База данных:** PostgreSQL 16 Alpine
- **Контейнеризация:** Docker (multi-stage build)
- **Оркестрация:** Docker Compose
- **Базовый образ:** python:3.11-alpine
- **Автоматизация:** Makefile (бонусное задание)

**Конфигурация:**

- **Порт приложения:** 8094
- **Порт БД:** 5432 (внутри Docker сети)
- **Система именования (slug):** as63-220012-v8
- **Тег образа:** annkrq/lab1-v8:stu-220012-v8
- **Переменные окружения:** STU_ID=220012, STU_GROUP=АС-63, STU_VARIANT=8

**Метаданные в артефактах:**

- **Dockerfile → LABEL:** org.bstu.student.fullname, org.bstu.student.id, org.bstu.group, org.bstu.variant, org.bstu.course
- **docker-compose.yml → labels:** org.bstu.owner, org.bstu.student.slug
- **Именование ресурсов:** app-as63-220012-v8, db-as63-220012-v8, net-as63-220012-v8, data-as63-220012-v8
- **База данных:** app_220012_v8

---

## Структура репозитория c описанием содержимого

```
task_01/
├── doc/
│   └── README.md              # отчет по лабораторной работе
└── src/
    ├── Dockerfile             # multi-stage сборка образа
    ├── docker-compose.yml     # конфигурация сервисов (app + db)
    ├── Makefile               # автоматизация команд сборки и запуска
    ├── requirements.txt       # зависимости Python (Flask)
    ├── .dockerignore          # исключения для сборки образа
    └── src/
        └── app.py             # Flask-приложение с graceful shutdown
```

**Описание компонентов:**

- **Dockerfile**: Multi-stage сборка (builder + final), непривилегированный пользователь UID 10001, HEALTHCHECK, LABELS с метаданными студента, минимизация размера образа
- **docker-compose.yml**: Два сервиса (app + db), именованные volume и network по шаблону slug, labels с метаданными, healthcheck для БД
- **app.py**: Flask-сервис с маршрутами /, /healthz, /echo; graceful shutdown по SIGTERM/SIGINT; логирование метаданных при старте
- **Makefile**: Команды для сборки, запуска, проверки health, просмотра логов, graceful shutdown теста

---

## Подробное описание выполнения

### Задание 1. Собрать минимальный образ для Flask HTTP-сервиса

**Требования:**

- Multi-stage сборка
- Финальный образ ≤ 150MB
- USER ненулевой
- EXPOSE/HEALTHCHECK корректны
- Конфигурация через переменные окружения

**Реализация:**

Создан Dockerfile с двумя стадиями:

**Builder stage:**

- Базовый образ: python:3.11-alpine
- Установка build-зависимостей (build-base)
- Установка Python-пакетов из requirements.txt в отдельный префикс (/install)
- Использование pip --no-cache-dir для уменьшения размера

**Final stage:**

- Базовый образ: python:3.11-alpine
- LABELS с метаданными студента (fullname, id, group, variant, course)
- ENV переменные (STU_ID, STU_GROUP, STU_VARIANT, APP_PORT, APP_HOST)
- Копирование только установленных пакетов из builder stage
- Создание непривилегированного пользователя (UID 10001)
- Копирование исходного кода приложения
- EXPOSE 8094
- HEALTHCHECK с проверкой /healthz эндпоинта
- USER 10001 (запуск от непривилегированного пользователя)
- ENTRYPOINT в exec-форме для корректной обработки сигналов

**Результат:** Образ размером менее 150MB, все требования выполнены.

`скриншот сборки образа`

### Задание 2. Оформить docker-compose.yml: приложение + зависимость (Postgres) + volume

**Требования:**

- Сервис приложения + сервис БД
- Volume для персистентности данных
- Именование по slug
- Labels с метаданными

**Реализация:**

Создан docker-compose.yml с двумя сервисами, именованным volume и сетью:

**Volumes:**

- `data-as63-220012-v8` — для персистентности данных PostgreSQL

**Networks:**

- `net-as63-220012-v8` — bridge-сеть для изоляции сервисов

**Сервис db (PostgreSQL):**

- Образ: postgres:16-alpine
- Имя контейнера: db-as63-220012-v8
- Volume: data-as63-220012-v8 → /var/lib/postgresql/data
- Healthcheck: pg_isready -U app_user (интервал 20s)
- ENV: POSTGRES_DB=app_220012_v8, POSTGRES_USER, POSTGRES_PASSWORD
- Labels: org.bstu.owner=annkrq, org.bstu.student.slug=as63-220012-v8

**Сервис app (Flask):**

- Сборка из Dockerfile с тегом annkrq/lab1-v8:stu-220012-v8
- Имя контейнера: app-as63-220012-v8
- Зависимость от db с условием service_healthy
- Порты: 8094:8094
- ENV: STU_ID=220012, STU_GROUP=АС-63, STU_VARIANT=8, APP_PORT=8094
- DATABASE_URL: postgresql://app_user:app_pass@db:5432/app_220012_v8
- Labels: org.bstu.owner=annkrq, org.bstu.student.slug=as63-220012-v8
- Security: read_only=true, tmpfs для /tmp

**Результат:** Оркестрация двух сервисов с корректными зависимостями и именованием.

`скриншот docker-compose up`

### Задание 3. Реализовать graceful shutdown (SIGTERM), проверить корректное завершение

**Требования:**

- Обработка сигнала SIGTERM
- Корректное завершение без потери данных
- Логирование процесса завершения

**Реализация:**

В app.py реализован обработчик сигналов SIGTERM и SIGINT:

```python
def initiate_graceful_shutdown(signum: int, _frame):
    logger.warning("Received signal %s - initiating graceful shutdown...", signum)
    if not shutdown_requested.is_set():
        logger.info("Stop accepting new connections. Shutdown flag set.")
    shutdown_requested.set()
```

**Механизм работы:**

1. При получении SIGTERM (docker stop) или SIGINT (Ctrl+C) вызывается обработчик
2. Логируется предупреждение с номером сигнала
3. Устанавливается флаг shutdown_requested
4. Сервер завершает обработку текущих запросов
5. Прекращает прием новых соединений
6. Логируется сообщение о завершении
7. Приложение корректно завершается

**Проверка:**

Команда для тестирования:

```powershell
make graceful-test
```

Логи при корректном завершении:

```text
2025-11-28 10:15:20,010 | WARNING | Received signal 15 - initiating graceful shutdown...
2025-11-28 10:15:20,311 | INFO | Stop accepting new connections. Shutdown flag set.
2025-11-28 10:15:20,512 | INFO | Graceful shutdown complete.
```

**Результат:** Приложение корректно завершается по SIGTERM без потери данных.

`скриншот graceful shutdown`

### Задание 4. Настроить кэширование зависимостей для ускорения повторной сборки

**Требования:**

- Использование механизмов кэширования Docker
- Ускорение повторной сборки при изменении кода

**Реализация:**

Применены следующие подходы для оптимизации сборки:

**1. Multi-stage build:**

- Зависимости устанавливаются в отдельной builder stage
- Только результат копируется в финальный образ
- При изменении кода не требуется переустановка зависимостей

**2. Правильный порядок COPY команд:**

```dockerfile
# Сначала копируем только requirements.txt
COPY requirements.txt ./
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

# Затем копируем код (при изменении кода слой с зависимостями берется из кэша)
COPY src/app.py ./
```

**3. Использование .dockerignore:**

- Исключение ненужных файлов из контекста сборки
- Уменьшение размера контекста
- Ускорение передачи контекста Docker daemon

**4. pip --no-cache-dir:**

- Исключение кэша pip из образа
- Уменьшение финального размера образа

**Результат:** При изменении только кода приложения, зависимости берутся из кэша, что значительно ускоряет сборку.

`скриншот повторной сборки с использованием кэша`

---

### Дополнительно: Автоматизация с помощью Makefile (бонус)

Создан Makefile для автоматизации команд:

**Доступные команды:**

- `make build` — сборка Docker образа
- `make up` / `make start` — запуск сервисов (app + db)
- `make health` — проверка health endpoint
- `make logs` — просмотр логов приложения
- `make stop` — остановка приложения
- `make down` — остановка и удаление всех сервисов
- `make ps` — статус контейнеров
- `make graceful-test` — тест graceful shutdown с проверкой логов
- `make shell` — открыть shell в контейнере приложения
- `make image-size` — проверка размера образа
- `make clean` — очистка Docker ресурсов

**Переменные Makefile:**

- SLUG=as63-220012-v8
- IMAGE=annkrq/lab1-v8:stu-220012-v8
- PORT=8094
- COMPOSE_PROJECT_NAME для корректной работы на Windows

`скриншот make help`

---

### Проверка работоспособности

**Логирование метаданных при старте:**

```text
2025-11-28 10:15:04,120 | INFO | ==== Application Startup ==== 
2025-11-28 10:15:04,121 | INFO | Student ID: 220012
2025-11-28 10:15:04,121 | INFO | Student Group: АС-63
2025-11-28 10:15:04,121 | INFO | Student Variant: 8
2025-11-28 10:15:04,122 | INFO | ENV STU_ID=220012
2025-11-28 10:15:04,122 | INFO | ENV STU_GROUP=АС-63
2025-11-28 10:15:04,122 | INFO | ENV STU_VARIANT=8
```

`скриншот логов старта`

**Health endpoint проверка:**

Команда:

```powershell
make health
```

Ответ:

```json
{
  "status": "ok",
  "timestamp": "2025-11-28T10:15:10.500Z"
}
```

`скриншот health check`

---

## Контрольный список (checklist)

### Основные требования (100 баллов)

**Корректность контейнеризации и образа (30 баллов):**

- [✅] Multi-stage сборка Dockerfile
- [✅] Размер финального образа ≤ 150MB
- [✅] Непривилегированный пользователь (USER 10001)
- [✅] EXPOSE 8094 корректно настроен
- [✅] HEALTHCHECK реализован и работает
- [✅] LABELS с метаданными студента

**Работа docker-compose (25 баллов):**

- [✅] docker-compose.yml с двумя сервисами (app + db)
- [✅] Зависимость PostgreSQL подключена
- [✅] Volume для персистентности данных (data-as63-220012-v8)
- [✅] Именованная сеть (net-as63-220012-v8)
- [✅] Healthcheck для БД (pg_isready)
- [✅] depends_on с условием service_healthy

**Graceful shutdown и логирование (20 баллов):**

- [✅] Обработка SIGTERM реализована
- [✅] Обработка SIGINT реализована
- [✅] Корректное завершение без потери данных
- [✅] Логирование процесса завершения
- [✅] Логирование метаданных при старте

**Именование и метаданные (15 баллов):**

- [✅] LABEL в Dockerfile (fullname, id, group, variant, course)
- [✅] labels в docker-compose (owner, student.slug)
- [✅] slug = as63-220012-v8
- [✅] ENV переменные: STU_ID, STU_GROUP, STU_VARIANT
- [✅] Тег образа: annkrq/lab1-v8:stu-220012-v8
- [✅] База данных: app_220012_v8
- [✅] Именование контейнеров по slug

**Оформление репозитория и README (10 баллов):**

- [✅] README с полными метаданными студента
- [✅] Описание шагов сборки и запуска
- [✅] Структура репозитория описана
- [✅] Логи старта и shutdown приложены

### Бонусные задания (+10 баллов)

**Оптимизация и автоматизация:**

- [✅] Кэширование зависимостей (.dockerignore, порядок COPY)
- [✅] Makefile для автоматизации команд
- [✅] Structured logging с timestamp
- [✅] Security настройки (read_only, tmpfs)
- [✅] HEALTHCHECK с настраиваемыми параметрами

**Итого:** 110 баллов (100 основных + 10 бонусных)

---

## Ссылки

Репозиторий с исходным кодом: `https://github.com/annkrq/RSiOT-2025`

---

## Вывод

В ходе выполнения лабораторной работы №1 были успешно освоены базовые навыки контейнеризации приложений с использованием Docker.

**Выполненные задачи:**

1. **Создан минимальный Docker образ** с использованием multi-stage сборки для Flask-приложения. Применены техники оптимизации: Alpine Linux как базовый образ, раздельная установка зависимостей, использование непривилегированного пользователя. Финальный размер образа составил менее 150MB.

2. **Настроена оркестрация сервисов** через docker-compose с двумя контейнерами: Flask-приложение и PostgreSQL база данных. Реализованы зависимости между сервисами, именованные volume для персистентности данных, изолированная сеть для безопасности.

3. **Реализован механизм graceful shutdown** с корректной обработкой сигналов SIGTERM и SIGINT, что обеспечивает безопасное завершение работы приложения без потери данных и соединений.

4. **Оптимизировано кэширование зависимостей** для ускорения повторной сборки образов при изменении кода приложения.

**Освоенные технологии и навыки:**

- Создание multi-stage Dockerfile для оптимизации размера образов
- Работа с непривилегированными пользователями (security best practices)
- Настройка HEALTHCHECK для мониторинга состояния контейнеров
- Организация микросервисной архитектуры через docker-compose
- Работа с Docker volume и networks
- Реализация graceful shutdown для production-ready приложений
- Соблюдение стандартов именования и метаданных (LABEL, labels, ENV)
- Автоматизация процессов через Makefile

**Использованные инструменты:**

- Docker Desktop v4.53.0
- Python 3.11 + Flask
- PostgreSQL 16 Alpine
- Docker Compose
- Makefile

Все задания выполнены в полном объеме согласно требованиям методических указаний. Дополнительно реализованы бонусные задачи: автоматизация через Makefile, оптимизация размера образа, structured logging, security настройки (read_only, tmpfs).
