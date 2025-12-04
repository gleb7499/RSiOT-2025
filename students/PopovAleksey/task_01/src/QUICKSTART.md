# Инструкция по запуску

## Быстрый старт

1. Перейдите в директорию с исходным кодом:

```bash
cd src
```

2. Запустите приложение:

```bash
docker-compose up --build
```

3. Проверьте работу приложения:

```bash
# Главная страница
curl http://localhost:8002/

# Health check
curl http://localhost:8002/live

# Список запросов
curl http://localhost:8002/requests
```

4. Остановка с graceful shutdown:

```bash
docker-compose down
```

## Проверка метаданных

Проверьте LABEL в образе:

```bash
docker inspect flask-app:stu-220051-v38 | grep -A 5 Labels
```

Проверьте контейнеры:

```bash
docker ps
```

Проверьте volume:

```bash
docker volume ls | grep data_v38
```

Проверьте сеть:

```bash
docker network ls | grep net-as64-220051-v38
```
