# Лабораторная работа №1

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
<p align="right">Группы АС-64</p>
<p align="right">Котковец К. В.</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А. Н.</p>
<br><br><br><br><br>
<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Научиться собирать минимальные образы (multi-stage) и запускать контейнеры под непривилегированным пользователем. Закрепить основы docker-compose: зависимости (БД/кэш), volume, сети. Настроить healthcheck и graceful shutdown.

---

### Вариант №35

## Метаданные студента

- **ФИО:** Котковец Кирилл Викторович
- **Группа:** АС-64
- **№ студенческого (StudentID):** 220044
- **Email (учебный):** <as006412@g.bstu.by>
- **GitHub username:** Kirill-Kotkovets
- **Вариант №:** 35
- **ОС и версия:** Windows 11 21H3, Docker Desktop v4.53.0

**Параметры варианта:**

- Стек: Python/Flask
- Порт: 8013
- Health endpoint: /ping
- Зависимость: Redis
- Volume: data_v35
- UID: 65532
- Тег: v35

---

## Окружение и инструменты

В данной лабораторной работе использовались:

- **Docker Desktop** v4.53.0 — для контейнеризации приложения
- **Python 3.11** — язык программирования
- **Flask 3.0.0** — веб-фреймворк
- **Redis 7-alpine** — база данных в памяти для кэширования
- **docker-compose** — для оркестрации контейнеров

---

## Структура репозитория с описанием содержимого

```text
task_01/
├── src/
│   ├── app.py                # Flask приложение с graceful shutdown
│   ├── requirements.txt      # Python зависимости
│   ├── Dockerfile            # Multi-stage Dockerfile
│   └── docker-compose.yml    # Конфигурация Docker Compose
└── doc/
    └── README.md             # Документация и отчет
```

---

## Подробное описание выполнения

### 1. Создание минимального HTTP-сервиса (Flask)

Создан простой Flask сервис (`src/app.py`) с двумя endpoints:

- `/` — основной endpoint, показывает информацию о студенте и счетчик посещений через Redis
- `/ping` — health check endpoint для проверки работоспособности

Приложение:

- Читает конфигурацию из переменных окружения
- Подключается к Redis с использованием префикса `stu:220044:v35`
- Логирует метаданные студента при старте
- Обрабатывает сигналы SIGTERM и SIGINT для graceful shutdown

### 2. Создание Dockerfile с multi-stage build

Dockerfile (`src/Dockerfile`) реализован в два этапа:

- **Builder stage**: установка зависимостей Python
- **Final stage**: создание минимального образа с установленными пакетами

Особенности:

- Базовый образ: `python:3.11-slim`
- Непривилегированный пользователь: `USER 65532`
- LABEL с метаданными студента
- HEALTHCHECK на endpoint `/ping`
- EXPOSE порта 8013

### 3. Создание docker-compose.yml

Файл `src/docker-compose.yml` включает:

- Сервис `app` — Flask приложение
- Сервис `redis` — Redis с персистентным хранилищем
- Volume `data_v35` для хранения данных Redis
- Сеть `net-as64-220044-v35` для связи между контейнерами
- Labels с метаданными студента (slug: `as64-220044-v35`)

### 4. Graceful shutdown

В `app.py` реализована обработка сигналов:

- При получении SIGTERM/SIGINT приложение:
  - Логирует информацию о получении сигнала
  - Закрывает соединение с Redis
  - Корректно завершает работу

### 5. Сборка и запуск

Команды для сборки и запуска:

```bash
cd src
docker-compose build
docker-compose up
```

Для проверки работы:

```bash
curl http://localhost:8013/
curl http://localhost:8013/ping
```

Для остановки с graceful shutdown:

```bash
docker-compose down
```

---

## Контрольный список (checklist)

- [✅] README с полными метаданными студента
- [✅] Dockerfile (multi-stage, non-root USER 65532, labels)
- [✅] docker-compose.yml с Redis и volume
- [✅] Health check endpoint /ping
- [✅] Graceful shutdown (обработка SIGTERM)
- [✅] Конфигурация через переменные окружения
- [✅] Логирование старта и shutdown
- [✅] Именование контейнеров/томов/сетей по slug
- [✅] Redis с префиксом ключей stu:220044:v35

---

## Вывод

В ходе выполнения лабораторной работы был создан контейнеризованный Flask сервис с использованием Docker и docker-compose. Освоены следующие навыки:

- Создание multi-stage Dockerfile для уменьшения размера образа
- Запуск контейнеров под непривилегированным пользователем (UID 65532)
- Настройка docker-compose для оркестрации нескольких сервисов
- Работа с volume для персистентного хранения данных Redis
- Реализация graceful shutdown для корректного завершения приложения
- Настройка healthcheck для мониторинга состояния контейнера
- Работа с переменными окружения для конфигурации приложения

Приложение успешно работает на порту 8013, взаимодействует с Redis, корректно обрабатывает запросы и завершается при получении сигнала остановки.
