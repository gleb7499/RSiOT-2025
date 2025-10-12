# ЛР01 — Контейнеризация и Docker

Цель
- Освоить базовые навыки работы с Docker: создание образов, запуск контейнеров, работа с volume и сетью.

Необходимые знания и ПО
- Установлены Docker Desktop/Engine.
- Базовые знания Linux CLI.

Задания
1) Соберите минимальный Docker-образ для простого HTTP-сервиса (на выбор: Python/Flask, Node/Express, Go net/http). Требования:
   - Многоступенчатая сборка (multi-stage), финальный образ ≤150MB.
   - Непривилегированный пользователь (USER), корректный EXPOSE/HEALTHCHECK.
   - Конфигурация через переменные окружения.
2) Оформите docker-compose.yml c сервисом, его зависимостью (например, Redis/Postgres) и volume для данных.
3) Реализуйте graceful shutdown (SIGTERM) и проверьте корректное завершение.
4) Настройте кэширование зависимостей, чтобы повторная сборка была быстрее.


Критерии приемки
- Репозиторий с Dockerfile и docker-compose.yml; README с шагами сборки/запуска.
- Логи демонстрируют старт, обработку запросов и корректный shutdown.
 - В README и артефактах присутствуют корректные метаданные студента (см. раздел ниже) и именование по варианту.
 
Формат сдачи
- PR в основную ветку. https://github.com/andreiNiasiuk как ревьювер.

Материалы
- Методическое руководство: [Лабораторная работа 01 — Методические материалы](./Лабораторная_работа_01_Методические_материалы.md)

## Метаданные студента (обязательно)
Укажите полную информацию о студенте и добавьте её в артефакты.

В отчете (README в корне репозитория) раздел «Метаданные студента»:
- ФИО (полностью)
- Группа
- № студенческого/зачетной книжки (StudentID)
- Email (учебный)
- GitHub username
- Вариант №
- Дата выполнения
- ОС (версия), версия Docker Desktop/Engine

В артефактах (требуется):
- Dockerfile → LABEL:
   - org.bstu.student.fullname = <ФИО>
   - org.bstu.student.id = <StudentID>
   - org.bstu.group = <Группа>
   - org.bstu.variant = <Номер варианта>
   - org.bstu.course = RSIOT
- docker-compose.yml → labels на сервисах:
   - org.bstu.owner = <GitHub username>
   - org.bstu.student.slug = <slug>

Где slug = <группа>-<StudentID>-v<вариант> (например, feis-41-12345-v07).

## Требования к именованию (ключи и ресурсы)
- Имена образов/тегов: добавьте суффикс с вариантом, напр.: :stu-<StudentID>-v<вариант>.
- Имена контейнеров/томов/сетей в compose: включают slug, напр.: container: app-<slug>, volume: data-<slug>, network: net-<slug>.
- ENV переменные конфигурации: добавьте переменные STU_ID, STU_GROUP, STU_VARIANT и используйте их в логе старта.
- Если используется Postgres: имя БД = app_<StudentID>_v<вариант> (или схема s<StudentID>_v<вариант>), пользователь/пароль без спецсимволов пробелов.
- Если используется Redis: префикс ключей = stu:<StudentID>:v<вариант>:<entity>.

## Варианты (24)
Укажите номер варианта в отчете и примените параметры ниже.

1) Стек: Node/Express; порт: 8081; health: /health; зависимость: Redis; volume: data_v1; UID: 65532; тег: v1
2) Стек: Python/Flask; порт: 8082; health: /healthz; зависимость: Postgres; volume: data_v2; UID: 10001; тег: v2
3) Стек: Go net/http; порт: 8083; health: /ready; зависимость: Redis; volume: data_v3; UID: 65532; тег: v3
4) Стек: Node/Express; порт: 8084; health: /live; зависимость: Postgres; volume: data_v4; UID: 10001; тег: v4
5) Стек: Python/Flask; порт: 8091; health: /ping; зависимость: Redis; volume: data_v5; UID: 65532; тег: v5
6) Стек: Go net/http; порт: 8092; health: /health; зависимость: Postgres; volume: data_v6; UID: 10001; тег: v6
7) Стек: Node/Express; порт: 8093; health: /ready; зависимость: Redis; volume: data_v7; UID: 65532; тег: v7
8) Стек: Python/Flask; порт: 8094; health: /live; зависимость: Postgres; volume: data_v8; UID: 10001; тег: v8
9) Стек: Go net/http; порт: 8071; health: /healthz; зависимость: Redis; volume: data_v9; UID: 65532; тег: v9
10) Стек: Node/Express; порт: 8072; health: /ping; зависимость: Postgres; volume: data_v10; UID: 10001; тег: v10
11) Стек: Python/Flask; порт: 8073; health: /health; зависимость: Redis; volume: data_v11; UID: 65532; тег: v11
12) Стек: Go net/http; порт: 8074; health: /ready; зависимость: Postgres; volume: data_v12; UID: 10001; тег: v12
13) Стек: Node/Express; порт: 8061; health: /live; зависимость: Redis; volume: data_v13; UID: 65532; тег: v13
14) Стек: Python/Flask; порт: 8062; health: /healthz; зависимость: Postgres; volume: data_v14; UID: 10001; тег: v14
15) Стек: Go net/http; порт: 8063; health: /ping; зависимость: Redis; volume: data_v15; UID: 65532; тег: v15
16) Стек: Node/Express; порт: 8064; health: /health; зависимость: Postgres; volume: data_v16; UID: 10001; тег: v16
17) Стек: Python/Flask; порт: 8051; health: /ready; зависимость: Redis; volume: data_v17; UID: 65532; тег: v17
18) Стек: Go net/http; порт: 8052; health: /live; зависимость: Postgres; volume: data_v18; UID: 10001; тег: v18
19) Стек: Node/Express; порт: 8053; health: /healthz; зависимость: Redis; volume: data_v19; UID: 65532; тег: v19
20) Стек: Python/Flask; порт: 8054; health: /ping; зависимость: Postgres; volume: data_v20; UID: 10001; тег: v20
21) Стек: Go net/http; порт: 8041; health: /health; зависимость: Redis; volume: data_v21; UID: 65532; тег: v21
22) Стек: Node/Express; порт: 8042; health: /ready; зависимость: Postgres; volume: data_v22; UID: 10001; тег: v22
23) Стек: Python/Flask; порт: 8043; health: /live; зависимость: Redis; volume: data_v23; UID: 65532; тег: v23
24) Стек: Go net/http; порт: 8044; health: /healthz; зависимость: Postgres; volume: data_v24; UID: 10001; тег: v24
