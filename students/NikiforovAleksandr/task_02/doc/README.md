# РСиОТ Лабораторная работа №2 - Kubernetes Deployment

## Описание

HTTP-сервис на Node.js/Express с эндпоинтами для Kubernetes (`/`, `/health`), использующий PostgreSQL для хранения данных.
Проект упакован в Docker-контейнер с multi-stage build и развернут в Kubernetes кластере с использованием Deployment, Service, ConfigMap, Secret.
Сервер поддерживает graceful shutdown, liveness/readiness probes и ресурсные ограничения.

## Структура проекта

```bash
task_02/
├── doc/                         # Документация проекта
│   └── README.md                # Основная документация
├── k8s/                         # Kubernetes манифесты
│   ├── namespace.yaml           # Namespace app16
│   ├── configmap.yaml           # Конфигурация приложения
│   ├── secret.yaml              # Секреты (DATABASE_URL)
│   ├── postgres.yaml            # PostgreSQL Deployment & Service
│   ├── deployment.yaml          # App Deployment с RollingUpdate
│   ├── service.yaml             # ClusterIP Service
│   └── ingress.yaml             # Ingress с nginx
├── src/
│   └── server.js                # Код приложения на Node.js
├── scripts/
│   ├── smoke-test.sh            # Скрипт для тестирования (Linux/Mac)
│   └── smoke-test.bat           # Скрипт для тестирования (Windows)
├── Dockerfile                   # Multi-stage build (≤ 150MB)
├── docker-compose.yml           # Локальный запуск
└── package.json, package-lock.json # Зависимости Node.js
```

## Требования

- Node.js 22+
- Docker + Docker Compose (для локального тестирования)
- Kubernetes кластер (Kind, Minikube или Docker Desktop Kubernetes)
- kubectl
- PostgreSQL 16+

## Быстрый запуск

### Локально с Docker Compose

```bash
docker compose up --build
```

Сервер доступен: `http://localhost:8064`

### В Kubernetes с Kind/Minikube

```bash
# Создать кластер (Kind)
kind create cluster --name app16-cluster

# Собрать образ
docker build -t lab01-node-express-pg:stu-220020-v16 .

# Загрузить образ в Kind
kind load docker-image lab01-node-express-pg:stu-220020-v16 --name app16-cluster

# Деплой приложения
kubectl apply -f k8s/

# Проверить статус
kubectl get all -n app16

# Доступ через port-forward
kubectl port-forward -n app16 service/web16-service 8064:8064
```

## Endpoints

- `GET /` - информация о сервисе и студенте (возвращает JSON с временной меткой)
- `GET /health` - проверка работоспособности PostgreSQL (используется для liveness/readiness probes)

## Конфигурация Kubernetes

### Deployment

- **Replicas:** 3
- **Strategy:** RollingUpdate (maxSurge: 1, maxUnavailable: 0)
- **Resources:** requests: 100m CPU/128Mi RAM, limits: 200m CPU/256Mi RAM
- **Liveness probe:** `/health` (initialDelay: 30s, interval: 10s, timeout: 5s)
- **Readiness probe:** `/health` (initialDelay: 5s, interval: 5s, timeout: 3s)
- **Graceful shutdown:** preStop hook с задержкой 10 секунд

### Service

- **Type:** ClusterIP (port: 8064)
- **Namespace:** app16

### Ingress

- **Host:** web16.local
- **Ingress Class:** nginx
- **Path:** /

### ConfigMap

- `PORT`: "8064"
- `NODE_ENV`: "production"
- `LOG_LEVEL`: "info"

### Secret

- `DATABASE_URL`: "postgres://app:app@postgres-service:5432/app_220020_v16"

### PostgreSQL Deployment

- **Image:** postgres:16-alpine
- **Database:** app_220020_v16
- **User:** app
- **Password:** app
- **Port:** 5432

## Особенности реализации

- **Multi-stage Docker build** (финальный образ ≤ 150MB)
- **Non-root пользователь** в контейнере (UID 10001, user: app)
- **Graceful shutdown** с обработкой SIGTERM/SIGINT для Kubernetes
- **Liveness/Readiness probes** с проверкой подключения к PostgreSQL
- **Resource limits** согласно варианту (CPU: 200m, Memory: 256Mi)
- **Environment variables** из ConfigMap и Secret
- **RollingUpdate strategy** для бесшовных обновлений
- **Health checks** в Dockerfile для локального тестирования

## Проверка деплоя

```bash
# Статус подов
kubectl get pods -n app16 -o wide

# Логи приложения
kubectl logs -n app16 -l app=web16 --tail=10

# Описание deployment
kubectl describe deployment web16-deployment -n app16

# Smoke-тест через скрипт
./scripts/smoke-test.sh  # Linux/Mac
scripts\smoke-test.bat   # Windows

# Ручной тест
kubectl port-forward -n app16 service/web16-service 8064:8064 &
curl http://localhost:8064/
curl http://localhost:8064/health
```

## Student Metadata

- **Full Name:** Александр Никифоров
- **Group:** AS-63
- **Student ID:** 220020
- **Email (Academic):** <AS63006320@g.bstu.by>
- **GitHub Username:** woqhy
- **Variant №:** 16
- **K8s Namespace:** app16
- **Deployment:** web16-deployment
- **Service:** web16-service
- **Service Port:** 8064
- **Database:** app_220020_v16
- **Docker Image:** lab01-node-express-pg:stu-220020-v16

## Технические спецификации

- **Operating System:** Windows 11 / Linux
- **Docker Version:** Docker Desktop 4.30+
- **Node.js Version:** 22-alpine (в контейнере)
- **PostgreSQL Version:** 16-alpine
- **Kubernetes Version:** 1.28+
- **Express Version:** 4.21.2
- **PostgreSQL Client:** pg 8.16.3

## Метки в Kubernetes

Все ресурсы содержат метки:

- `org.bstu.course: RSIOT`
- `org.bstu.variant: "16"`
- `org.bstu.student.id: "220020"`
- `org.bstu.group: AS-63`
- `org.bstu.owner: woqhy`
- `org.bstu.student.slug: as-63-220020-v16`
- `org.bstu.student.fullname: "Nikiforov-Alexandr"`

## Health Checks

- **Docker HEALTHCHECK:** `CMD wget -qO- http://127.0.0.1:8064/health`
- **K8s Liveness Probe:** HTTP GET на `/health` (initialDelay: 30s)
- **K8s Readiness Probe:** HTTP GET на `/health` (initialDelay: 5s)
