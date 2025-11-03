# «Метаданные студента»:

- ФИО - Ващук Анатолий Валерьевич
- Группа - АС-63
- № студенческого/зачетной книжки (StudentID) - 220036
- Email (учебный) - as006407@g.bstu.by
- GitHub username - skevet-avramuk
- Вариант № - 1
- Дата выполнения - 03.11.2025
- ОС (версия), версия Docker Desktop/Engine - Windows 11, Docker version 27.4.0

---

## RSOT Project

This project is a minimal Node/Express server with Redis integration, Dockerized for easy setup.

---

### 1. Build and start containers

```bash
- docker compose up -d --build
```

### 2. Verify the application

#### 2.1 Healthcheck

```bash
Invoke-RestMethod http://localhost:8081/health
# or
curl http://localhost:8081/health
```

Expected output: OK

#### 2.2 Readiness

```bash
Invoke-RestMethod http://localhost:8081/ready
```

Expected output: READY

#### 2.3 Redis

```bash
curl http://localhost:8081/visit
```

Expected output: Количество визитов: 1 (increments with each request)

#### 4. Stop and remove containers

```bash
docker compose down
```

If you need to remove Redis data volumes as well, add -v:

```bash
docker compose down -v
```
