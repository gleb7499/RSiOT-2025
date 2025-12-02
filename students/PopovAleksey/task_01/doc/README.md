# Лабораторная работа №1

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
<p align="right">Попов Алексей Сергеевич</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А.Н.</p>
<br><br><br><br><br>
<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Освоить базовые навыки работы с Docker: создание образов, запуск контейнеров, работа с volume и сетью. Научиться собирать минимальные образы (multi-stage) и запускать контейнеры под непривилегированным пользователем. Закрепить основы docker-compose: зависимости (БД/кэш), volume, сети. Настроить healthcheck и graceful shutdown.

---

### Вариант №38

## Метаданные студента

- **ФИО:** Попов Алексей Сергеевич
- **Группа:** АС-64
- **№ студенческого (StudentID):** 220051
- **Email (учебный):** <as006416@g.bstu.by>
- **GitHub username:** LexusxdsD
- **Вариант №:** 38
- **ОС и версия:** Windows 11 21H2, Docker Desktop v4.52.0

**Параметры варианта:**

- Стек: Python/Flask
- Порт: 8002
- Health endpoint: /live
- Зависимость: Postgres
- Volume: data_v38
- UID: 10001
- Тег: v38

---

## Окружение и инструменты

В лабораторной работе использовались следующие технологии и инструменты:

- **Docker Desktop v4.52.0** — контейнеризация приложения
- **Python 3.11** — язык программирования
- **Flask 3.0.0** — веб-фреймворк для создания HTTP-сервиса
- **PostgreSQL 15** — реляционная база данных
- **psycopg2** — адаптер для работы с PostgreSQL из Python
- **docker-compose** — оркестрация multi-container приложения

---

## Структура репозитория c описанием содержимого

```text
task_01/
├── src/                           # исходники сервиса
│   ├── app.py                     # Flask приложение с endpoints и graceful shutdown
│   ├── requirements.txt           # Python зависимости
│   ├── Dockerfile                 # multi-stage образ с USER 10001
│   └── docker-compose.yml         # оркестрация: app + postgres + volume
└── doc/
    └── README.md                  # документация и отчет
```

---

## Подробное описание выполнения

### 1. Создание Flask приложения

Реализован простой HTTP-сервис с использованием Flask, который:

- Логирует метаданные студента при старте (STU_ID, STU_GROUP, STU_VARIANT)
- Предоставляет endpoint `/live` для healthcheck
- Подключается к PostgreSQL и создает таблицу для хранения запросов
- Записывает информацию о запросах в БД
- Реализует graceful shutdown через обработку сигнала SIGTERM

**Endpoints:**

- `GET /` — главная страница с информацией о студенте
- `GET /live` — healthcheck endpoint (проверяет подключение к БД)
- `GET /requests` — получение последних 10 запросов из БД

### 2. Создание Dockerfile (multi-stage)

Dockerfile состоит из двух стадий:

**Stage 1 (builder):**

- Использует образ `python:3.11-slim`
- Устанавливает зависимости из `requirements.txt`

**Stage 2 (runtime):**

- Копирует установленные зависимости из builder
- Добавляет LABEL с метаданными студента
- Создает непривилегированного пользователя с UID 10001
- Переключается на USER 10001
- Настраивает EXPOSE 8002 и HEALTHCHECK для `/live`

### 3. Настройка docker-compose.yml

Сконфигурирован docker-compose с двумя сервисами:

**app:**

- Собирается из локального Dockerfile
- Тег образа: `flask-app:stu-220051-v38`
- Имя контейнера: `app-as64-220051-v38`
- Переменные окружения для подключения к БД
- Зависит от сервиса `db`
- Labels с метаданными студента

**db:**

- Образ: `postgres:15-alpine`
- Имя контейнера: `db-as64-220051-v38`
- Volume: `data_v38` для персистентности данных
- Настройки БД соответствуют требованиям именования

**Сеть:** `net-as64-220051-v38`

### 4. Реализация graceful shutdown

В `app.py` реализован обработчик сигнала SIGTERM:

```python
def signal_handler(signum, frame):
    logger.info(f"Получен сигнал {signum}. Начинаем graceful shutdown...")
    shutdown_flag = True
    logger.info("Приложение корректно завершено")
    sys.exit(0)

signal.signal(signal.SIGTERM, signal_handler)
```

При получении SIGTERM приложение:

1. Логирует получение сигнала
2. Устанавливает флаг shutdown
3. Логирует корректное завершение
4. Завершает работу

### 5. Запуск и проверка

**Сборка и запуск:**

```bash
cd src
docker-compose up --build
```

**Проверка healthcheck:**

```bash
curl http://localhost:8002/live
```

Ожидаемый ответ:

```json
{"status": "healthy", "database": "connected"}
```

**Проверка graceful shutdown:**

```bash
docker stop app-as64-220051-v38
```

В логах должно появиться:

```text
Получен сигнал 15. Начинаем graceful shutdown...
Приложение корректно завершено
```

---

## Контрольный список (checklist)

- [✅] README с полными метаданными студента
- [✅] Dockerfile (multi-stage, non-root USER 10001, LABEL с метаданными)
- [✅] docker-compose.yml (app + postgres, volume data_v38, сеть, labels)
- [✅] EXPOSE порт 8002
- [✅] HEALTHCHECK на endpoint /live
- [✅] Конфигурация через переменные окружения (STU_ID, STU_GROUP, STU_VARIANT)
- [✅] Graceful shutdown (обработка SIGTERM)
- [✅] Volume для данных PostgreSQL
- [✅] Зависимости между сервисами (depends_on)
- [✅] Именование по требованиям (slug: as64-220051-v38)
- [✅] Логирование старта и завершения

---

## Вывод

В ходе выполнения лабораторной работы были освоены базовые навыки работы с Docker:

- Создан multi-stage Dockerfile для минимизации размера финального образа
- Настроен запуск контейнера под непривилегированным пользователем (UID 10001)
- Реализовано взаимодействие между несколькими контейнерами через docker-compose
- Настроен volume для персистентности данных PostgreSQL
- Реализован healthcheck для мониторинга состояния приложения
- Добавлена поддержка graceful shutdown через обработку сигнала SIGTERM
- Освоена работа с переменными окружения для конфигурации приложения
- Изучены требования к именованию и метаданным в Docker

Приложение полностью функционально, соответствует требованиям варианта №38 и готово к развертыванию.
