# «Метаданные студента»

- ФИО - Ярмолович Александр Сергеевич
- Группа - АС-63
- № студенческого/зачетной книжки (StudentID) - 220029
- Email (учебный) -as006326@g.bstu.by
- GitHub username - yarmolov
- Вариант № - 24
- Дата выполнения - 03.11.2025
- ОС (версия), версия Docker Desktop/Engine - Windows 10, Docker version 28.4.0

## RSOT Проект

Минимальное веб-приложение на Go с использованием стандартного пакета net/http и подключением к базе данных PostgreSQL. Контейнеризировано для быстрого развёртывания через Docker, использует порт 8044, точку проверки /healthz и собственный volume data_v24.

### 1. Клонировать репазиторий

```bash
git clone https://github.com/yarmolov/RSOT/Lab1.git
cd RSOT
```

### 2. Сбор и запуск контейнера

```bash
- docker compose up -d --build
```

### 3. Подтвердить приложение

#### 3.1 Healthcheck

```bash
Invoke-RestMethod http://localhost:8044/healthz
```

Content           : OK

#### 3.2 Readiness

```bash
Invoke-RestMethod http://localhost:8044/ready
```

Content           : READY

#### 3.3 Posgres

```bash
Invoke-RestMethod -Uri http://localhost:8044/visit
```

Студент: 24, Группа: feis, Вариант: v24, Кол-во визитов: 4

### 4. Остановка и удаление контейнера

```bash
docker compose down
```
