# Отчёт по лабораторной работе №1

## Метаданные студента

- **ФИО:** Иванюк Иван Александрович
- **Группа:** АС-64
- **№ студенческого/зачетной книжки (StudentID):** 220041
- **Email (учебный):** AS006412@g.bstu.by
- **GitHub username:** JonF1re
- **Вариант №:** 8
- **Дата выполнения:** 07.10.2025
- **ОС (версия):** Windows 11 Home 24H2
- **Версия Docker Desktop:** 4.49.0

## Шаги сборки и запуска

1. Прописываем в консоли docker compose up -d --build
2. Прописываем к консоли docker compose up -d

## Логи

Starting Flask application on 0.0.0.0:8094...
 * Serving Flask app 'app'
 * Debug mode: off
Running on http://127.0.0.1:8094⁠
172.18.0.1 - - [27/Nov/2025 15:20:13] "GET / HTTP/1.1" 200 -
172.18.0.1 - - [27/Nov/2025 15:23:23] "GET /slow HTTP/1.1" 200 -
172.18.0.1 - - [27/Nov/2025 15:23:12] "GET /live HTTP/1.1" 200 -
Received signal 15, initiating graceful shutdown...
Shutdown complete