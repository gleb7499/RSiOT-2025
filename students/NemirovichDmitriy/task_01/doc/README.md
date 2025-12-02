# Лабораторная работа №01

<p align="center">Министерство образования Республики Беларусь</p>
<p align="center">Учреждение образования</p>
<p align="center">"Брестский Государственный технический университет"</p>
<p align="center">Кафедра ИИТ</p>
<br><br><br><br><br><br>
<p align="center"><strong>Лабораторная работа №01</strong></p>
<p align="center"><strong>По дисциплине:</strong> "Распределенные системы и облачные технологии"</p>
<p align="center"><strong>Тема:</strong> Контейнеризация и Docker</p>
<br><br><br><br><br><br>
<p align="right"><strong>Выполнил:</strong></p>
<p align="right">Студент 4 курса</p>
<p align="right">Группы АС-64</p>
<p align="right">Немирович Д.А.</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А.Н.</p>
<br><br><br><br><br>
<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Научиться собирать минимальные Docker-образы с использованием multi-stage сборки, запускать контейнеры под непривилегированным пользователем. Закрепить основы docker-compose: работа с зависимостями (Redis), volume и сети. Настроить healthcheck и graceful shutdown.

---

### Вариант №37

**Стек:** Node/Express  
**Порт:** 8001  
**Health:** /ready  
**Зависимость:** Redis  
**Volume:** data_v37  
**UID:** 65532  
**Тег:** v37

## Метаданные студента

- **ФИО:** Немирович Дмитрий Александрович
- **Группа:** АС-64
- **№ студенческого (StudentID):** 220050
- **Email (учебный):** <as006415@g.bstu.by>
- **GitHub username:** goryachiy-ugolek
- **Вариант №:** 37
- **Дата выполнения:** 29.11.2025
- **ОС и версия:** Windows 10 1809, Docker Desktop v4.53.0

---

## Окружение и инструменты

В лабораторной работе использовались:

- **Node.js** v18 — платформа для запуска JavaScript на сервере
- **Express.js** v4.18.2 — веб-фреймворк для Node.js
- **Redis** v7 (alpine) — in-memory база данных для кэширования
- **Docker** Desktop v4.53.0 — для контейнеризации приложения
- **docker-compose** v3.8 — для оркестрации multi-container приложений

Согласно варианту №37:

- HTTP-сервис на порту 8001
- Health check endpoint: `/ready`
- Зависимость: Redis с volume `data-as64-220050-v37`
- Непривилегированный пользователь: UID 65532

---

## Структура репозитория c описанием содержимого

```text
task_01/
├── README.md                # Полная документация и отчет
└── src/                     # Все файлы проекта
    ├── app.js              # Основной файл приложения Express
    ├── package.json        # Зависимости Node.js
    ├── Dockerfile          # Multi-stage образ приложения
    ├── docker-compose.yml  # Оркестрация app + Redis
    └── .gitignore          # Игнорирование файлов
```

---

## Подробное описание выполнения

### 1. Создание простого HTTP-сервиса (Node/Express)

Создан простой Express-приложение (`src/app.js`), которое:

- Слушает порт 8001
- Имеет endpoint `/ready` для health check
- Подключается к Redis и использует префикс ключей `stu:220050:v37:`
- Логирует метаданные студента (STU_ID, STU_GROUP, STU_VARIANT) при старте
- Обрабатывает запросы и инкрементирует счетчик в Redis

### 2. Создание Dockerfile (multi-stage)

Dockerfile (расположен в `src/Dockerfile`) состоит из двух стадий:

- **Builder stage**: использует `node:18`, устанавливает зависимости
- **Production stage**: использует `node:18-slim` для уменьшения размера образа
- Установлен USER 65532:65532 (непривилегированный пользователь)
- Добавлены LABEL с метаданными студента
- Настроен HEALTHCHECK с проверкой endpoint `/ready`
- Конфигурация через переменные окружения (PORT, STU_ID, STU_GROUP, STU_VARIANT)

### 3. Создание docker-compose.yml

Файл `src/docker-compose.yml` настраивает два сервиса:

- **app**: приложение Node.js, зависит от Redis
- **redis**: Redis v7-alpine с persistent volume

Именование согласно требованиям:

- Контейнеры: `app-as64-220050-v37`, `redis-as64-220050-v37`
- Volume: `data-as64-220050-v37`
- Network: `net-as64-220050-v37`
- Image tag: `rsiot-v37:stu-220050-v37`
- Labels: `org.bstu.owner`, `org.bstu.student.slug`

### 4. Реализация graceful shutdown

В `app.js` реализована обработка сигналов SIGTERM и SIGINT:

- При получении сигнала устанавливается флаг `isShuttingDown`
- Закрывается HTTP сервер
- Закрывается подключение к Redis
- Логируется информация о завершении работы

### 5. Сборка и запуск

**Команды для запуска:**

```cmd
# Переход в директорию src
cd src

# Сборка образа
docker-compose build

# Запуск контейнеров
docker-compose up -d

# Просмотр логов
docker-compose logs -f app

# Остановка (для проверки graceful shutdown)
docker-compose down
```

**Проверка работы:**

```cmd
# Проверка health endpoint
curl http://localhost:8001/ready

# Проверка основного endpoint
curl http://localhost:8001/
```

### Логи старта приложения

```text
[INFO] Starting application...
[INFO] StudentID: 220050
[INFO] Group: АС-64
[INFO] Variant: 37
[INFO] Port: 8001
[INFO] Redis Host: redis:6379
[INFO] Redis Prefix: stu:220050:v37:
[INFO] Redis connected
[INFO] Server started on port 8001
[INFO] Application ready to handle requests
```

### Логи graceful shutdown

```text
[INFO] SIGTERM received, starting graceful shutdown...
[INFO] HTTP server closed
[INFO] Redis connection closed
[INFO] Graceful shutdown completed
```

---

## Контрольный список (checklist)

- [ ✅ ] README с полными метаданными студента
- [ ✅ ] Dockerfile (multi-stage, non-root USER 65532, labels)
- [ ✅ ] docker-compose.yml с app + Redis
- [ ✅ ] Volume для данных Redis (data-as64-220050-v37)
- [ ✅ ] Health check endpoint /ready
- [ ✅ ] HEALTHCHECK в Dockerfile
- [ ✅ ] Graceful shutdown (SIGTERM/SIGINT)
- [ ✅ ] Логирование метаданных студента при старте
- [ ✅ ] Конфигурация через переменные окружения
- [ ✅ ] Именование согласно требованиям (slug, теги, labels)
- [ ✅ ] Интеграция с Redis (префикс ключей stu:220050:v37:)

---

## Вывод

В ходе лабораторной работы были освоены базовые навыки работы с Docker: создание multi-stage образов для минимизации размера, запуск контейнеров под непривилегированным пользователем (UID 65532).

Реализовано простое HTTP-приложение на Node.js/Express с интеграцией Redis, настроен docker-compose для оркестрации сервисов с использованием volume и изолированной сети.

Настроен health check endpoint `/ready` и HEALTHCHECK в Dockerfile для мониторинга состояния контейнера. Реализован graceful shutdown с корректной обработкой сигналов SIGTERM/SIGINT и закрытием всех соединений.

Все артефакты оформлены согласно требованиям методички: добавлены метаданные студента в LABEL, labels в docker-compose, корректное именование контейнеров, volume и сетей с использованием slug `as64-220050-v37`.
