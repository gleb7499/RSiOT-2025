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
<p align="right">Группы АС-64</p>
<p align="right">Котковец К. В.</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А. Н.</p>
<br><br><br><br><br>
<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Научиться готовить Kubernetes-манифесты для простого HTTP-сервиса (Deployment + Service), настроить liveness/readiness probes и политику обновления (rolling update), подготовить конфигурацию через ConfigMap и научиться запускать кластер локально.

---

### Вариант №35

## Метаданные студента

- **ФИО:** Котковец Кирилл Викторович
- **Группа:** АС-64
- **№ студенческого (StudentID):** 220044
- **Email (учебный):** <as006412@g.bstu.by>
- **GitHub username:** Kirill-Kotkovets
- **Вариант №:** 35
- **ОС и версия:** Windows 11 21H3

---

## Окружение и инструменты

- **Docker Desktop:** v4.53.0
- **Kubernetes:** встроенный в Docker Desktop
- **kubectl:** v1.28+
- **Язык:** Go 1.21
- **Образ:** Alpine 3.18

Параметры варианта:

- namespace: `app35`
- имя сервиса: `web35`
- replicas: `2`
- port: `8013`
- cpu: `150m`
- memory: `128Mi`

---

## Структура репозитория c описанием содержимого

```text
task_02/
├── src/
│   ├── main.go              # HTTP-сервис на Go
│   ├── go.mod               # Go модуль
│   ├── Dockerfile           # Multi-stage сборка образа
│   └── k8s/                 # Kubernetes манифесты
│       ├── namespace.yaml   # Namespace app35
│       ├── configmap.yaml   # ConfigMap с переменными
│       ├── deployment.yaml  # Deployment с пробами
│       └── service.yaml     # Service (NodePort)
└── doc/
    └── README.md            # Документация
```

---

## Подробное описание выполнения

### 1. Подготовка HTTP-сервиса

Создан простой HTTP-сервис на Go с endpoints:

- `/` - основной endpoint
- `/health` - liveness probe
- `/ready` - readiness probe

Сервис логирует старт, метаданные студента (STU_ID, STU_GROUP, STU_VARIANT) и graceful shutdown при получении SIGTERM.

### 2. Создание Dockerfile

Создан multi-stage Dockerfile:

- Этап сборки: golang:1.21-alpine
- Финальный образ: alpine:3.18 (≤ 150 MB)
- Запуск от non-root пользователя (appuser)
- LABEL с метаданными студента
- EXPOSE 8013

### 3. Kubernetes манифесты

Созданы манифесты:

**namespace.yaml** - создание namespace `app35`

**configmap.yaml** - конфигурация с переменными:

- STU_ID: "220044"
- STU_GROUP: "АС-64"
- STU_VARIANT: "35"

**deployment.yaml** - деплой с параметрами:

- replicas: 2
- RollingUpdate стратегия (maxSurge: 1, maxUnavailable: 0)
- resources: cpu 150m, memory 128Mi
- livenessProbe: /health
- readinessProbe: /ready

**service.yaml** - Service типа NodePort на порту 30035

### 4. Проверка работы

**Сборка образа:**

```bash
cd src
docker build -t web35:latest .
```

**Применение манифестов:**

```bash
kubectl apply -f src/k8s/namespace.yaml
kubectl apply -f src/k8s/configmap.yaml
kubectl apply -f src/k8s/deployment.yaml
kubectl apply -f src/k8s/service.yaml
```

**Проверка статуса:**

```bash
kubectl get pods -n app35
kubectl get svc -n app35
kubectl describe deployment app-as64-220044-v35 -n app35
```

**Тестирование endpoint:**

```bash
curl http://localhost:30035/
curl http://localhost:30035/health
```

---

## Контрольный список (checklist)

- [✅] README с полными метаданными студента
- [✅] Dockerfile (multi-stage, non-root, labels)
- [✅] Kubernetes манифесты (namespace, configmap, deployment, service)
- [✅] Health/Liveness/Readiness probes
- [✅] Старт/остановка: логирование и graceful shutdown
- [✅] RollingUpdate стратегия
- [✅] Ресурсные лимиты (cpu/memory)

---

## Вывод

В рамках лабораторной работы был создан минимальный HTTP-сервис на Go и подготовлены Kubernetes-манифесты для его деплоя. Настроены liveness и readiness проbes, политика RollingUpdate и ресурсные ограничения. Сервис успешно разворачивается в локальном Kubernetes кластере с 2 репликами и доступен через NodePort.
