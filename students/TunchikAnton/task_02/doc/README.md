# РСиОТ Лабораторная работа №2 - Kubernetes Deployment

## Описание

HTTP-сервис на Go с эндпоинтами для Kubernetes (`/live`, `/ready`, `/metrics`), использующий Redis для хранения состояния.
Проект упакован в Docker-контейнер с multi-stage build и развернут в Kubernetes кластере с использованием Deployment, Service, ConfigMap.
Сервер поддерживает graceful shutdown, liveness/readiness probes и ресурсные ограничения.

## Структура проекта

```
anton-lab-1/
├── k8s/                         # Kubernetes манифесты
│   ├── namespace.yaml           # Namespace app21
│   ├── configmap.yaml           # Конфигурация приложения
│   ├── redis.yaml               # Redis Deployment & Service
│   ├── deployment.yaml          # App Deployment с RollingUpdate
│   ├── service.yaml             # NodePort Service
│   └── ingress.yaml             # Ingress с nginx
├── cmd/server/
│   └── main.go                  # Обновленный код с K8s probes
├── scripts/
│   └── smoke.sh                 # Скрипт для тестирования
├── Dockerfile                   # Multi-stage build
├── docker-compose.yml           # Локальный запуск
├── go.mod, go.sum               # Зависимости Go
├── .env                         # Переменные окружения
└── README.md                    # Документация
```

## Требования

- Docker + Docker Compose (для локального тестирования)
- Kubernetes кластер (Kind, Minikube или Docker Desktop Kubernetes)
- kubectl

## Быстрый запуск

### Локально с Docker Compose

```bash
docker compose up --build
```

Сервер доступен: `http://localhost:8041`

### В Kubernetes с Kind

```bash
# Создать кластер
kind create cluster --name app21-cluster

# Собрать и загрузить образ
docker build -t rsiot-app:stu-220026-v21 .
kind load docker-image rsiot-app:stu-220026-v21 --name app21-cluster

# Деплой приложения
kubectl apply -f k8s/

# Проверить статус
kubectl get all -n app21

# Доступ через port-forward
kubectl port-forward -n app21 service/web21-service 8080:8041
```

## Endpoints

- `GET /` - информация о сервисе и студенте
- `GET /health` - проверка работоспособности Redis (совместимость с ЛР01)
- `GET /live` - **liveness probe** для Kubernetes
- `GET /ready` - **readiness probe** для Kubernetes
- `GET /hit` - увеличение счетчика в Redis
- `GET /metrics` - метрики приложения для мониторинга

## Конфигурация Kubernetes

### Deployment

- **Replicas:** 2
- **Strategy:** RollingUpdate (maxSurge: 1, maxUnavailable: 0)
- **Resources:** requests: 100m CPU/64Mi RAM, limits: 150m CPU/128Mi RAM
- **Liveness probe:** `/live` (interval: 10s, timeout: 5s)
- **Readiness probe:** `/ready` (interval: 5s, timeout: 3s)

### Service

- **Type:** NodePort (port: 8041, nodePort: 30081)
- **Namespace:** app21

### ConfigMap

- Конфигурация порта, студенческих данных, Redis адреса

## Особенности реализации

- **Multi-stage Docker build** (финальный образ ~45MB)
- **Non-root пользователь** в контейнере (UID 1000)
- **Graceful shutdown** с логированием для Kubernetes
- **Liveness/Readiness probes** с проверкой Redis
- **Resource limits** согласно варианту (CPU: 150m, Memory: 128Mi)
- **Structured logging** с K8s-метаданными
- **RollingUpdate strategy** для бесшовных обновлений

## Проверка деплоя

```bash
# Статус подов
kubectl get pods -n app21 -o wide

# Логи приложения
kubectl logs -n app21 -l app=web21 --tail=10

# Описание deployment
kubectl describe deployment web21-deployment -n app21

# Smoke-тест
kubectl port-forward -n app21 service/web21-service 8080:8041 &
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/metrics
```

## Student Metadata

- **Full Name:** Tunchik Anton Dmitrievich
- **Group:** AS-63
- **Student ID:** 220026
- **Email (Academic):** AS63006326@g.bstu.by
- **GitHub Username:** Stis25
- **Variant №:** 21
- **K8s Namespace:** app21
- **Deployment:** web21
- **Service Port:** 8041
- **Operating System:** Windows 11 23H2
- **Docker Version:** Docker Desktop 4.47.0
- **Kubernetes:** Kind v0.20.0 / Docker Desktop

## Метки в Kubernetes

Все ресурсы содержат метки:
- `org.bstu.course: RSIOT`
- `org.bstu.variant: "21"`
- `org.bstu.student.id: "220026"`
- `org.bstu.group: AS-63`
- `org.bstu.owner: Stis25`
- `org.bstu.student.slug: AS-63-220026-v21`
