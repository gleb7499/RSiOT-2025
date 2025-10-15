# ЛР01 — Контейнеризация и Docker

## Вариант 26
- Стек: Python/Flask
- Порт: 8032
- Health: /health
- Зависимость: Postgres
- Volume: data_v26
- UID: 10001
- Тег: v26

---

## Метаданные студента
- **ФИО:** Кужир Владислав Витальевич
- **Группа:** AS-64
- **№ студенческого:** 220047
- **Email:** vkuzir7@gmail.com
- **GitHub username:** XD-cods
- **Вариант №:** 26
- **Дата выполнения:** 23.09.2025
- **ОС:** Windows 10, Docker Desktop 4.46.0

---




Задания

1. Собрать минимальный образ для простого HTTP‑сервиса (Python/Flask, Node/Express, Go net/http на выбор):
        Multi‑stage; финальный образ ≤ 150MB
        USER ненулевой; EXPOSE/HEALTHCHECK корректны
        Конфигурация через переменные окружения
        ![        Multi‑stage; финальный образ ≤ 150MB](image_weight.png)
2. Оформить docker-compose.yml: приложение + зависимость (Redis/Postgres) + volume для данных
3. Реализовать graceful shutdown (SIGTERM), проверить корректное завершение
4. Настроить кэширование зависимостей для ускорения повторной сборки


## Шаги сборки и запуска

1. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/XD-cods/RSiOT-2025
   cd RSiOT-2025
   ```
2. Соберите и запустите контейнеры:
   ```sh
   docker-compose up -d --build
   ```
3. Проверьте логи приложения:
   ```sh
   docker-compose logs app
   ```
4. Проверьте работу сервиса:
   - Главная: [http://localhost:8032/](http://localhost:8032/)
   - Health: [http://localhost:8032/health](http://localhost:8032/health)

---

## Структура проекта
- `Dockerfile` — многоступенчатая сборка, LABEL с метаданными
- `docker-compose.yml` — сервис Flask, Postgres, volume, labels, slug
- `app.py` — Flask HTTP-сервис, graceful shutdown (SIGTERM)
- `requirements.txt` — зависимости
- `init.sql` — инициализация БД
- `.dockerignore` — исключение лишних файлов

---

## Примеры логов
```
Starting Flask app | STU_ID=220047 | GROUP=AC-64 | VARIANT=26

 * Serving Flask app 'app'

 * Debug mode: off

WARNING: This is a development server. Do not use it in a production deployment. Use a production WSGI server instead.

 * Running on all addresses (0.0.0.0)

 * Running on http://127.0.0.1:8032⁠

 * Running on http://172.18.0.3:8032⁠
Press CTRL+C to quit

192.168.65.1 - - [23/Sep/2025 06:43:42] "GET / HTTP/1.1" 200 -
192.168.65.1 - - [23/Sep/2025 06:43:42] "GET /favicon.ico HTTP/1.1" 404 -
192.168.65.1 - - [23/Sep/2025 06:43:48] "GET /health HTTP/1.1" 200 -
192.168.65.1 - - [23/Sep/2025 06:43:49] "GET /favicon.ico HTTP/1.1" 404 -
Received SIGTERM, shutting down gracefully...

```

---

## Проверка shutdown
Для graceful shutdown выполните:
```sh
docker-compose stop app
```
В логах появится сообщение о корректном завершении.
