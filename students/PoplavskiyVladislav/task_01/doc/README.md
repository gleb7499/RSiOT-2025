# «Метаданные студента»

Full name: Vladislav Vladimirovich Poplavsky
Group: ASOI 63
Student ID: 220021
Email: as006318@g.bstu.by
GitHub: ImRaDeR1
Variant: 17
Date: 2025-11-03
OS: Windows 10 / Docker Desktop 4.33.0

## RSOT

Проект Минимальное веб-приложение на Python с использованием Flask и подключением к базе данных Redis.
Контейнеризировано для быстрого развёртывания через Docker, использует порт 8051, точку проверки `/Ready` и собственный volume `data_v17`.

### 1. Сбор и запуск контейнера

```bash
docker compose up -d --build
```

### 2 Readiness

```bash
Invoke-RestMethod http://localhost:8044/ready
```

status ok

Студент: 17, Группа: feis, Вариант: v17, Кол-во визитов: 4

### 3. Остановка и удаление контейнера

```bash
docker compose down
```
