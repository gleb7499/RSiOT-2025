# Лабораторная работа №1. Контейнеризация и Docker

<p align="center">Министерство образования Республики Беларусь</p>
<p align="center">Учреждение образования</p>
<p align="center">"Брестский Государственный технический университет"</p>
<p align="center">Кафедра ИИТ</p>
<br><br><br><br><br><br>
<p align="center"><strong>Лабораторная работа №1</strong></p>
<p align="center"><strong>По дисциплине:</strong> "Распределенные системы и облачные технологии"</p>
<p align="center"><strong>Тема:</strong> Контейнеризация и Docker</p>
<br><br><br><br><br><br>
<p align="right"><strong>Выполнил:</strong></p>
<p align="right">Студент 4 курса</p>
<p align="right">Группы АС-63</p>
<p align="right">Кульбеда Кирилл Александрович</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Иванин Д.Н.</p>
<br><br><br><br><br>
<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Научиться собирать минимальные образы (multi-stage) и запускать контейнеры под непривилегированным пользователем. Закрепить основы docker-compose: зависимости (БД/кэш), volume, сети. Настроить healthcheck.

---

### Вариант №12

## Метаданные студента

- **ФИО:** Кульбеда Кирилл Александрович
- **Группа:** АС-63
- **№ студенческого (StudentID):** 220016
- **Email (учебный):** <AS006313@g.bstu.by>
- **GitHub username:** fr0ogi
- **Вариант №:** 12
- **ОС и версия:** Windows 11 21H4, Docker Desktop v4.53.0

**Slug:** `as63-220016-v12`

---

## Окружение и инструменты

- **Язык:** Go (net/http)
- **База данных:** PostgreSQL 15
- **Порт приложения:** 8074
- **Health endpoint:** /ready
- **Volume:** data-as63-220016-v12
- **UID:** 10001
- **Образ тег:** stu-220016-v12

---

## Структура репозитория c описанием содержимого

```text
task_01/
├── src/
│   ├── main.go              # Исходный код HTTP-сервиса
│   ├── go.mod               # Зависимости Go
│   ├── go.sum               # Checksums зависимостей
│   ├── Dockerfile           # Multi-stage образ
│   └── docker-compose.yml   # Оркестрация сервисов
└── doc/
    └── README.md            # Документация и отчёт
```

---

## Подробное описание выполнения

### 1. Создание минимального HTTP-сервиса на Go

Реализован простой HTTP-сервер с двумя endpoints:

- `/` - основной endpoint, сохраняет запросы в PostgreSQL
- `/ready` - health check endpoint для проверки готовности

Сервис логирует метаданные студента при старте (STU_ID, STU_GROUP, STU_VARIANT).

### 2. Dockerfile с multi-stage сборкой

Dockerfile состоит из двух стадий:

- **Builder stage:** компиляция Go приложения
- **Final stage:** минимальный alpine образ (~20-30 MB)

Особенности:

- Использован непривилегированный пользователь (UID 10001)
- Добавлены LABEL с метаданными студента
- Настроен HEALTHCHECK на endpoint `/ready`
- EXPOSE 8074

### 3. Docker Compose

Создан `docker-compose.yml` с двумя сервисами:

- **app:** Go приложение
- **postgres:** база данных PostgreSQL

Конфигурация:

- Volume `data-as63-220016-v12` для персистентности данных БД
- Сеть `net-as63-220016-v12` для изоляции
- Зависимость app от postgres (depends_on)
- Labels с метаданными на каждом сервисе

### 4. Запуск и проверка

**Сборка и запуск:**

```bash
cd src
docker-compose up --build
```

**Проверка работы:**

```bash
# Основной endpoint
curl http://localhost:8074/

# Health check
curl http://localhost:8074/ready
```

**Остановка:**

```bash
docker-compose down
```

---

## Контрольный список (checklist)

- [✅] README с полными метаданными студента
- [✅] Dockerfile (multi-stage, non-root, labels)
- [✅] docker-compose.yml с PostgreSQL
- [✅] Volume для данных БД
- [✅] Сеть для изоляции сервисов
- [✅] HEALTHCHECK в Dockerfile
- [✅] Переменные окружения для конфигурации
- [✅] Логирование метаданных при старте
- [❌] Graceful shutdown (не реализовано)
- [❌] Оптимизация кэширования зависимостей

---

## Примеры логов

**Старт приложения:**

```text
app-as63-220016-v12  | 2025/12/01 10:00:00 Starting application...
app-as63-220016-v12  | 2025/12/01 10:00:00 Student ID: 220016
app-as63-220016-v12  | 2025/12/01 10:00:00 Group: АС-63
app-as63-220016-v12  | 2025/12/01 10:00:00 Variant: 12
app-as63-220016-v12  | 2025/12/01 10:00:00 Connected to PostgreSQL
app-as63-220016-v12  | 2025/12/01 10:00:00 Server listening on port 8074
```

**Обработка запроса:**

```text
app-as63-220016-v12  | 2025/12/01 10:01:15 Request received: GET /
```

---

## Вывод

В ходе выполнения лабораторной работы был создан минимальный HTTP-сервис на Go с использованием PostgreSQL в качестве базы данных. Реализован multi-stage Dockerfile с непривилегированным пользователем, настроен docker-compose для оркестрации сервисов. Освоены базовые навыки контейнеризации: сборка образов, работа с volume, сетями и зависимостями между контейнерами. Приложение успешно запускается, обрабатывает запросы и использует healthcheck для проверки готовности.
