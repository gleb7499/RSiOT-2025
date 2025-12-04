# Лабораторная работа №02

<p align="center">Министерство образования Республики Беларусь</p>
<p align="center">Учреждение образования</p>
<p align="center">"Брестский Государственный технический университет"</p>
<p align="center">Кафедра ИИТ</p>
<br><br><br><br><br><br>
<p align="center"><strong>Лабораторная работа №02</strong></p>
<p align="center"><strong>По дисциплине:</strong> "Распределенные системы и облачные технологии"</p>
<p align="center"><strong>Тема:</strong> Kubernetes: базовый деплой</p>
<br><br><br><br><br><br>
<p align="right"><strong>Выполнил:</strong></p>
<p align="right">Студент 4 курса</p>
<p align="right">Группы АС-63</p>
<p align="right">Логинов Г. О.</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А. Н.</p>
<br><br><br><br><br>
<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Научиться готовить Kubernetes-манифесты для простого HTTP-сервиса (Deployment + Service), настроить liveness/readiness probes и политику обновления (rolling update), подготовить конфигурацию через ConfigMap/Secret, смонтировать volume для данных при необходимости, научиться запускать кластер локально (Kind/Minikube) и проверять корректность деплоя.

---

### Вариант №14

```text
ns=app14, name=web14, replicas=3, port=8062, ingressClass=nginx, cpu=200m, mem=192Mi
```

## Метаданные студента

| Поле | Значение |
|------|----------|
| **ФИО** | Логинов Глеб Олегович |
| **Группа** | АС-63 |
| **№ студенческого (StudentID)** | 220018 |
| **Email (учебный)** | <as006315@g.bstu.by> |
| **GitHub username** | gleb7499 |
| **Вариант №** | 14 |
| **Дата выполнения** | 28.11.2025 |
| **ОС и версия** | Windows 11 24H2 |
| **Docker Desktop** | v4.45.0 |
| **kubectl** | v1.31.0 |
| **Kind** | v0.24.0 |
| **Minikube** | v1.34.0 |

### Slug и Labels

- **slug:** `as63-220018-v14`
- **Префиксы ресурсов:** `app-<slug>`, `data-<slug>`, `net-<slug>`

### Labels/Annotations в манифестах

```yaml
labels:
  org.bstu.owner: gleb7499
  org.bstu.student.slug: as63-220018-v14
  org.bstu.course: RSIOT
  org.bstu.student.id: "220018"
  org.bstu.group: "АС-63"
  org.bstu.variant: "14"
  org.bstu.student.fullname: "Логинов Глеб Олегович"

annotations:
  org.bstu.student.fullname: "Логинов Глеб Олегович"
  org.bstu.description: "Flask HTTP service for Kubernetes lab (variant 14)"
```

---

## Окружение и инструменты

| Инструмент | Версия | Назначение |
|------------|--------|------------|
| Docker Desktop | v4.45.0 | Контейнеризация |
| kubectl | v1.31.0 | CLI для Kubernetes |
| Kind | v0.24.0 | Локальный Kubernetes кластер |
| Minikube | v1.34.0 | Альтернативный локальный кластер |
| Kustomize | встроен в kubectl | Управление манифестами |
| Python | 3.11 | Язык сервиса |
| Flask | 3.x | HTTP фреймворк |
| PostgreSQL | 16-alpine | База данных |

---

## Структура репозитория с описанием содержимого

```text
task_02/
├── doc/
│   └── README.md               # Документация (этот файл)
├── materials/
│   ├── extra-info.md           # Дополнительная информация от преподавателя
│   └── previous-promt.md       # Методичка и каркас отчета
└── src/
    ├── app/                    # Исходный код приложения
    │   ├── .dockerignore       # Исключения для Docker
    │   ├── Dockerfile          # Multi-stage Dockerfile
    │   ├── requirements.txt    # Зависимости Python
    │   └── src/
    │       └── app.py          # Flask HTTP-сервис
    └── k8s/                    # Kubernetes манифесты
        ├── kustomization.yaml  # Kustomize конфигурация
        ├── namespace.yaml      # Namespace app14
        ├── configmap.yaml      # Конфигурация приложения (ENV)
        ├── secret.yaml         # Секреты (DB credentials)
        ├── pvc.yaml            # PersistentVolumeClaim для PostgreSQL
        ├── db-deployment.yaml  # Deployment для PostgreSQL
        ├── db-service.yaml     # Service для PostgreSQL (ClusterIP)
        ├── deployment.yaml     # Deployment для Flask приложения
        ├── service.yaml        # Service для Flask (ClusterIP)
        └── ingress.yaml        # Ingress (nginx) для внешнего доступа
```

---

## Подробное описание выполнения

### 1. Подготовка HTTP-сервиса и контейнерного образа

#### Dockerfile (multi-stage build)

Dockerfile находится в `task_02/src/app/Dockerfile` и использует multi-stage сборку для минимизации размера образа:

```dockerfile
# ---- Builder stage ----
FROM python:3.11-alpine AS builder

ARG PIP_NO_CACHE_DIR=1
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

RUN apk add --no-cache build-base=0.5-r3

WORKDIR /app

COPY requirements.txt ./
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

# ---- Final stage ----
FROM python:3.11-alpine AS final

LABEL org.bstu.student.fullname="Логинов Глеб Олегович" \
      org.bstu.student.id="220018" \
      org.bstu.group="АС-63" \
      org.bstu.variant="14" \
      org.bstu.course="RSIOT" \
      org.bstu.owner="gleb7499" \
      org.bstu.student.slug="as63-220018-v14"

ENV STU_ID=220018 \
    STU_GROUP=АС-63 \
    STU_VARIANT=14 \
    APP_PORT=8062 \
    APP_HOST=0.0.0.0 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# Создание non-root пользователя для безопасности
RUN adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder /install /usr/local
COPY src/app.py ./

EXPOSE 8062

RUN apk add --no-cache wget=1.25.0-r1 \
    && echo 'hosts: files dns' > /etc/nsswitch.conf

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-hsts -qO- http://127.0.0.1:${APP_PORT}/healthz || exit 1

# Запуск от non-root пользователя (UID 10001)
USER 10001

ENTRYPOINT ["python", "app.py"]
```

**Особенности Dockerfile:**

- **Multi-stage build** — первый stage (builder) устанавливает зависимости, второй stage (final) содержит только runtime
- **Non-root пользователь** — UID 10001 (appuser) для безопасности
- **Health endpoint** — `/healthz` для liveness/readiness проверок
- **Labels** — все необходимые метаданные org.bstu.*
- **Финальный образ ≤ 150 MB** — благодаря Alpine и multi-stage

#### Сборка образа

```bash
# Из директории task_02/src/app
cd task_02/src/app
docker build -t gleb7499/lab1-v14:stu-220018-v14 .

# Проверка размера образа (должен быть ≤ 150MB)
docker images gleb7499/lab1-v14:stu-220018-v14

# Запуск для проверки
docker run -d -p 8062:8062 --name test-app gleb7499/lab1-v14:stu-220018-v14

# Проверка health endpoint
curl http://localhost:8062/healthz

# Проверка логов (должны быть STU_ID, STU_GROUP, STU_VARIANT)
docker logs test-app

# Проверка graceful shutdown
docker stop test-app

# Очистка
docker rm test-app
```

#### Flask приложение (app.py)

Приложение реализует:

- **Health endpoint** `/healthz` — возвращает JSON `{"status": "ok", "timestamp": "..."}`
- **Логирование при запуске** — выводит STU_ID, STU_GROUP, STU_VARIANT
- **Graceful shutdown** — корректная обработка SIGTERM/SIGINT
- **Echo endpoint** `/echo` — для тестирования POST запросов

**Образ:** `gleb7499/lab1-v14:stu-220018-v14`

### 2. Kubernetes-манифесты

#### 2.1 Namespace (namespace.yaml)

Создан namespace `app14` согласно варианту:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: app14
  labels:
    org.bstu.owner: gleb7499
    org.bstu.student.slug: as63-220018-v14
    org.bstu.course: RSIOT
    org.bstu.student.id: "220018"
    org.bstu.group: "АС-63"
    org.bstu.variant: "14"
    org.bstu.student.fullname: "Логинов Глеб Олегович"
```

#### 2.2 Deployment для Flask приложения (deployment.yaml)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-as63-220018-v14
  namespace: app14
spec:
  replicas: 3  # Согласно варианту
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0  # Нулевой downtime при обновлениях
      maxSurge: 1
  template:
    spec:
      containers:
        - name: web14
          image: gleb7499/lab1-v14:stu-220018-v14
          ports:
            - containerPort: 8062  # Порт согласно варианту
          resources:
            requests:
              memory: "96Mi"
              cpu: "100m"
            limits:
              memory: "192Mi"  # Согласно варианту
              cpu: "200m"      # Согласно варианту
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8062
            initialDelaySeconds: 10
            periodSeconds: 15
          readinessProbe:
            httpGet:
              path: /healthz
              port: 8062
            initialDelaySeconds: 5
            periodSeconds: 10
          securityContext:
            runAsNonRoot: true
            runAsUser: 10001
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false
```

**Ключевые параметры согласно варианту 14:**

- `replicas: 3`
- `containerPort: 8062`
- `limits.memory: 192Mi`
- `limits.cpu: 200m`

#### 2.3 Service (service.yaml)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: net-as63-220018-v14
  namespace: app14
spec:
  type: ClusterIP
  selector:
    app: web14
  ports:
    - name: http
      port: 80
      targetPort: 8062
```

#### 2.4 Ingress (ingress.yaml)

Согласно варианту используется `ingressClass=nginx`:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: net-ingress-as63-220018-v14
  namespace: app14
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: web14.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: net-as63-220018-v14
                port:
                  number: 80
```

#### 2.5 ConfigMap (configmap.yaml)

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config-as63-220018-v14
  namespace: app14
data:
  STU_ID: "220018"
  STU_GROUP: "АС-63"
  STU_VARIANT: "14"
  APP_PORT: "8062"
  APP_HOST: "0.0.0.0"
```

#### 2.6 Secret (secret.yaml)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: data-secret-as63-220018-v14
  namespace: app14
type: Opaque
stringData:
  DATABASE_URL: "postgresql://app_user:app_pass@data-db-as63-220018-v14:5432/app_220018_v14"
  POSTGRES_USER: "app_user"
  POSTGRES_PASSWORD: "app_pass"
  POSTGRES_DB: "app_220018_v14"
```

#### 2.7 PersistentVolumeClaim (pvc.yaml)

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: data-pvc-as63-220018-v14
  namespace: app14
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

#### 2.8 PostgreSQL Deployment (db-deployment.yaml)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: data-db-as63-220018-v14
  namespace: app14
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: postgres
          image: postgres:16-alpine
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: data-secret-as63-220018-v14
                  key: POSTGRES_USER
            # ... other env vars from secret
          livenessProbe:
            exec:
              command: ["pg_isready", "-U", "app_user"]
          readinessProbe:
            exec:
              command: ["pg_isready", "-U", "app_user"]
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql/data
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: data-pvc-as63-220018-v14
```

### 3. Liveness и Readiness Probes

**Flask приложение:**

- `livenessProbe`: HTTP GET `/healthz`, port 8062
  - initialDelaySeconds: 10
  - periodSeconds: 15
  - timeoutSeconds: 3
  - failureThreshold: 3
- `readinessProbe`: HTTP GET `/healthz`, port 8062
  - initialDelaySeconds: 5
  - periodSeconds: 10
  - timeoutSeconds: 3
  - failureThreshold: 3

**PostgreSQL:**

- `livenessProbe`: exec `pg_isready -U app_user`
  - initialDelaySeconds: 30
  - periodSeconds: 20
- `readinessProbe`: exec `pg_isready -U app_user`
  - initialDelaySeconds: 10
  - periodSeconds: 10

### 4. Kustomize для управления манифестами (БОНУС)

Используется Kustomize для параметризации и централизованного управления labels/annotations:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: app14

labels:
  - pairs:
      org.bstu.owner: gleb7499
      org.bstu.student.slug: as63-220018-v14
      org.bstu.course: RSIOT
      org.bstu.student.id: "220018"
      org.bstu.group: "АС-63"
      org.bstu.variant: "14"

commonAnnotations:
  org.bstu.student.fullname: "Логинов Глеб Олегович"

resources:
  - namespace.yaml
  - configmap.yaml
  - secret.yaml
  - pvc.yaml
  - db-deployment.yaml
  - db-service.yaml
  - deployment.yaml
  - service.yaml
  - ingress.yaml
```

**Преимущества Kustomize:**

- Централизованное управление labels на всех ресурсах
- Автоматическое добавление annotations
- Параметризация namespace
- Стандартный синтаксис без нестандартных расширений

---

## Инструкции для локального тестирования

### Развертывание с Kind

#### Создание кластера

```bash
# Создание кластера с именем lab2
kind create cluster --name lab2

# Проверка подключения
kubectl cluster-info --context kind-lab2
```

#### Загрузка образа в Kind

```bash
# Сборка образа (из task_02/src/app)
cd task_02/src/app
docker build -t gleb7499/lab1-v14:stu-220018-v14 .

# Загрузка образа в кластер Kind
kind load docker-image gleb7499/lab1-v14:stu-220018-v14 --name lab2
```

#### Применение манифестов

```bash
# С использованием Kustomize (рекомендуется)
kubectl apply -k task_02/src/k8s/

# Или применение манифестов по отдельности
kubectl apply -f task_02/src/k8s/namespace.yaml
kubectl apply -f task_02/src/k8s/configmap.yaml
kubectl apply -f task_02/src/k8s/secret.yaml
kubectl apply -f task_02/src/k8s/pvc.yaml
kubectl apply -f task_02/src/k8s/db-deployment.yaml
kubectl apply -f task_02/src/k8s/db-service.yaml
kubectl apply -f task_02/src/k8s/deployment.yaml
kubectl apply -f task_02/src/k8s/service.yaml
kubectl apply -f task_02/src/k8s/ingress.yaml
```

### Развертывание с Minikube

#### Создание кластера Minikube

```bash
# Запуск Minikube
minikube start --driver=docker

# Включение Ingress контроллера
minikube addons enable ingress

# Проверка статуса
minikube status
```

#### Загрузка образа

```bash
# Использование Docker daemon Minikube
eval $(minikube docker-env)
cd task_02/src/app
docker build -t gleb7499/lab1-v14:stu-220018-v14 .
```

#### Применение манифестов в Minikube

```bash
kubectl apply -k task_02/src/k8s/
```

#### Доступ через Ingress

```bash
# Получение IP Minikube
MINIKUBE_IP=$(minikube ip)
echo "Minikube IP: $MINIKUBE_IP"

# Добавление записи в hosts файл:
# Linux/macOS: sudo sh -c "echo '$MINIKUBE_IP web14.local' >> /etc/hosts"
# Windows (PowerShell от администратора):
#   Add-Content -Path "C:\Windows\System32\drivers\etc\hosts" -Value "$MINIKUBE_IP web14.local"
#
# Пример записи: 192.168.49.2 web14.local

# Проверка доступа
curl http://web14.local/healthz
```

### Проверка статусов

```bash
# Все ресурсы в namespace
kubectl get all -n app14

# Подробности о подах
kubectl get pods -n app14 -o wide

# Статус deployments
kubectl get deployments -n app14

# Статус services
kubectl get svc -n app14

# Статус ingress
kubectl get ingress -n app14

# Описание пода (для диагностики)
kubectl describe pod -n app14 -l app=web14
```

### Smoke-тест проверка HTTP-эндпоинта

```bash
# Port-forward для доступа к приложению
kubectl port-forward -n app14 svc/net-as63-220018-v14 8062:80 &

# Проверка healthz
curl http://localhost:8062/healthz

# Проверка главной страницы
curl http://localhost:8062/

# Проверка echo endpoint
curl -X POST http://localhost:8062/echo -H "Content-Type: application/json" -d '{"test":"hello"}'
```

### Просмотр логов

```bash
# Логи Flask приложения
kubectl logs -n app14 -l app=web14 --tail=50

# Логи PostgreSQL
kubectl logs -n app14 -l app=postgres-db --tail=50

# Стриминг логов
kubectl logs -n app14 -l app=web14 -f
```

### Rolling Update проверка

```bash
# Обновление образа (симуляция rolling update)
kubectl set image deployment/app-as63-220018-v14 web14=gleb7499/lab1-v14:stu-220018-v14 -n app14

# Наблюдение за обновлением
kubectl rollout status deployment/app-as63-220018-v14 -n app14

# Проверка истории обновлений
kubectl rollout history deployment/app-as63-220018-v14 -n app14
```

### Удаление ресурсов

```bash
# Удаление всех ресурсов через Kustomize
kubectl delete -k task_02/src/k8s/

# Или удаление namespace (удалит все ресурсы внутри)
kubectl delete namespace app14

# Удаление кластера Kind
kind delete cluster --name lab2

# Удаление кластера Minikube
minikube delete
```

---

## Пример логов работы

```text
2025-11-28 10:15:04,120 | INFO | ==== Application Startup ==== 
2025-11-28 10:15:04,121 | INFO | Student ID: 220018
2025-11-28 10:15:04,121 | INFO | Student Group: АС-63
2025-11-28 10:15:04,121 | INFO | Student Variant: 14
2025-11-28 10:15:04,122 | INFO | DATABASE_URL: postgresql://app_user:***@data-db-as63-220018-v14:5432/app_220018_v14
2025-11-28 10:15:04,122 | INFO | ENV STU_ID=220018
2025-11-28 10:15:04,122 | INFO | ENV STU_GROUP=АС-63
2025-11-28 10:15:04,122 | INFO | ENV STU_VARIANT=14
2025-11-28 10:15:04,123 | INFO | ================================
2025-11-28 10:15:04,123 | INFO | Starting Flask server on 0.0.0.0:8062
```

## Пример запроса к `/healthz`

```bash
curl http://localhost:8062/healthz
```

Пример ответа:

```json
{
  "status": "ok",
  "timestamp": "2025-11-28T07:15:10.500Z"
}
```

---

## Описание компонентов

### Deployment Flask App (app-as63-220018-v14)

| Параметр | Значение |
|----------|----------|
| Replicas | 3 |
| Image | `gleb7499/lab1-v14:stu-220018-v14` |
| Port | 8062 |
| CPU Requests | 100m |
| CPU Limits | 200m |
| Memory Requests | 96Mi |
| Memory Limits | 192Mi |
| Liveness Probe | HTTP GET /healthz (period: 15s) |
| Readiness Probe | HTTP GET /healthz (period: 10s) |
| Security | runAsNonRoot, runAsUser: 10001, readOnlyRootFilesystem |
| Strategy | RollingUpdate (maxUnavailable: 0, maxSurge: 1) |

### Deployment PostgreSQL (data-db-as63-220018-v14)

| Параметр | Значение |
|----------|----------|
| Replicas | 1 |
| Image | `postgres:16-alpine` |
| Port | 5432 |
| CPU Requests | 100m |
| CPU Limits | 500m |
| Memory Requests | 128Mi |
| Memory Limits | 256Mi |
| Storage | PVC 1Gi |
| Liveness Probe | pg_isready (period: 20s) |
| Readiness Probe | pg_isready (period: 10s) |

### Service Flask App (net-as63-220018-v14)

| Параметр | Значение |
|----------|----------|
| Type | ClusterIP |
| Port | 80 |
| TargetPort | 8062 |

### Service PostgreSQL (data-db-as63-220018-v14)

| Параметр | Значение |
|----------|----------|
| Type | ClusterIP |
| Port | 5432 |
| TargetPort | 5432 |

### Ingress (net-ingress-as63-220018-v14)

| Параметр | Значение |
|----------|----------|
| IngressClass | nginx |
| Host | web14.local |
| Path | / |
| Backend | net-as63-220018-v14:80 |

### ConfigMap (app-config-as63-220018-v14)

| Ключ | Значение |
|------|----------|
| STU_ID | 220018 |
| STU_GROUP | АС-63 |
| STU_VARIANT | 14 |
| APP_PORT | 8062 |
| APP_HOST | 0.0.0.0 |

### Secret (data-secret-as63-220018-v14)

| Ключ | Описание |
|------|----------|
| DATABASE_URL | URL подключения к PostgreSQL |
| POSTGRES_USER | Имя пользователя БД |
| POSTGRES_PASSWORD | Пароль пользователя БД |
| POSTGRES_DB | Имя базы данных |

### PersistentVolumeClaim (data-pvc-as63-220018-v14)

| Параметр | Значение |
|----------|----------|
| AccessModes | ReadWriteOnce |
| Storage | 1Gi |

---

## Контрольный список (checklist)

| Требование | Статус |
|------------|--------|
| README с полными метаданными студента | ✅ |
| Dockerfile (multi-stage, non-root, labels) в task_02 | ✅ |
| Финальный образ ≤ 150 MB | ✅ |
| Non-root пользователь (UID 10001) | ✅ |
| Health endpoints (/healthz) | ✅ |
| Graceful shutdown (SIGTERM/SIGINT) | ✅ |
| Логирование STU_ID, STU_GROUP, STU_VARIANT | ✅ |
| Kubernetes Deployment с replicas=3 | ✅ |
| Kubernetes Deployment с RollingUpdate strategy | ✅ |
| Kubernetes Deployment с resources limits (cpu=200m, mem=192Mi) | ✅ |
| Kubernetes Service (ClusterIP) | ✅ |
| Kubernetes Ingress (ingressClass=nginx) | ✅ |
| Kubernetes ConfigMap | ✅ |
| Kubernetes Secret | ✅ |
| Kubernetes PersistentVolumeClaim | ✅ |
| Liveness Probe (HTTP) | ✅ |
| Readiness Probe (HTTP) | ✅ |
| Labels org.bstu.* на всех ресурсах | ✅ |
| Annotations org.bstu.student.fullname | ✅ |
| Namespace app14 (согласно варианту) | ✅ |
| Именование ресурсов с префиксами app-/data-/net- | ✅ |
| Инструкции для Kind | ✅ |
| Инструкции для Minikube | ✅ |
| Smoke-test проверка | ✅ |
| Kustomize для управления манифестами (БОНУС) | ✅ |
| PVC + демонстрация использования для PostgreSQL (БОНУС) | ✅ |

---

## Вывод

В ходе выполнения лабораторной работы №02 были закреплены навыки развертывания приложений в Kubernetes:

1. **Подготовлены Kubernetes-манифесты** для HTTP-сервиса Flask, включая Deployment, Service, Ingress, ConfigMap, Secret и PersistentVolumeClaim.

2. **Настроены liveness и readiness probes** для автоматического отслеживания состояния приложения и перезапуска при сбоях.

3. **Настроена стратегия RollingUpdate** с параметрами `maxUnavailable: 0` и `maxSurge: 1` для обеспечения нулевого downtime при обновлениях.

4. **Установлены ресурсные лимиты** согласно варианту 14 (cpu=200m, mem=192Mi).

5. **Создан Ingress** с ingressClass=nginx для внешнего доступа к приложению.

6. **Использован Kustomize** для централизованного управления labels, annotations и параметризации манифестов.

7. **Подготовлены инструкции** для локального тестирования с использованием Kind и Minikube.

Все ресурсы именованы согласно требованиям с использованием префиксов `app-`, `data-`, `net-` и slug студента `as63-220018-v14`. Все манифесты содержат необходимые labels и annotations согласно методическим указаниям.

**Освоенные инструменты:** Docker, kubectl, Kubernetes (Deployment, Service, Ingress, ConfigMap, Secret, PVC), Kustomize, Kind, Minikube.
