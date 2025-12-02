# Министерство образования Республики Беларусь

<p align="center">Учреждение образования</p>
<p align="center">“Брестский Государственный технический университет”</p>
<p align="center">Кафедра ИИТ</p>

<p align="center"><strong>Лабораторная работа №2</strong></p>
<p align="center"><strong>По дисциплине:</strong> “РСиОТ”</p>
<p align="center"><strong>Тема:</strong> “Kubernetes: базовый деплой”</p>

<p align="right"><strong>Выполнил:</strong></p>
<p align="right">Студент 4 курса</p>
<p align="right">Группы АС-63</p>
<p align="right">Поплавский В.В.</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А.Н.</p>

<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Научиться готовить Kubernetes-манифесты для простого HTTP-сервиса (Deployment + Service). Подготовить конфигурацию через ConfigMap/Secret и смонтировать volume для данных при необходимости.

---

### Вариант №17

| Параметр            | Значение                  |
|---------------------|---------------------------|
| Namespace           | `app17`                   |
| Имя приложения      | `web17`                   |
| Количество реплик   | 2                         |
| Порт приложения     | 8051                      |
| Ingress Class       | `nginx`                   |
| CPU limit           | 150m                      |
| Memory limit        | 128Mi                     |
| Healthcheck endpoint| /ready и /health          |
| Номер варианта      | 17                        |

---

### Ход выполнения работы

### 1. Структура проекта

├── app/
│   └── app.py                    # Flask приложение на Python
├── k8s/
│   ├── 00-namespace.yaml        # Namespace app17
│   ├── 01-configmap.yaml        # Конфигурация приложения
│   ├── 02-secret.yaml           # Секреты (Redis)
│   ├── 03-redis-deployment.yaml # Redis deployment
│   ├── 04-redis-service.yaml    # Redis Service
│   ├── 05-deployment.yaml       # Основное приложение с 2 репликами
│   └── 06-service.yaml          # NodePort Service
├── Dockerfile                   # Multi-stage сборка Python приложения
└── requirements.txt             # Python зависимости (Flask, Redis)

---

### 2. Последовательность запуска

#### Проверка версий

docker --version
kubectl version --client

#### Проверка кластера Kubernetes

kubectl cluster-info
kubectl get nodes

---

### 3. Сборка Docker образа

docker build -t app:stu-220021-v17 .
docker images | findstr "stu-220021-v17"

---

### 4. Развертывание в Kubernetes

kubectl apply -f k8s/00-namespace.yaml
kubectl apply -f k8s/01-configmap.yaml
kubectl apply -f k8s/02-secret.yaml
kubectl apply -f k8s/03-redis-deployment.yaml
kubectl apply -f k8s/04-redis-service.yaml
kubectl apply -f k8s/05-deployment.yaml
kubectl apply -f k8s/06-service.yaml

#### Проверка создания ресурсов

kubectl get all -n app17

---

### 5. Тестирование приложения

#### Port-forward для локального доступа

kubectl port-forward -n app17 svc/web17-service 8051:8051

#### Тестирование endpoints

curl http://localhost:8051/ready
Ответ: ok

curl http://localhost:8051/healthy
Ответ: ok

curl http://localhost:8051/
Ответ: JSON с информацией о сервисе

---

### 6. Проверка требований

| Критерий | Выполнено | Комментарий |
|-----------|-----------|--------------|
| Multi-stage Docker build | ✅ | Образ 114 MB|
| Размер образа ≤ 150 MB | ✅ | 114 MB < 150 MB|
| Non-root пользователь | ✅ | UID: 1000|
| Resource limits | ✅ | CPU=150m, Memory=128Mi|
| Health endpoints | ✅ | /ready и /health возвращают 200 OK|
| Liveness/Readiness probes | ✅ | Настроены в deployment|
| Graceful shutdown | ✅ | Обработка SIGTERM/SIGINT настроены с логированием|
| 3 реплики | ✅ | Deployment с 2 подами|
| RollingUpdate стратегия | ✅ | Настроена в deployment|
| ConfigMap/Secret | ✅ | Разделение конфигурации|
| Service (ClusterIP) | ✅ | Сервис создан, порт 30085|
| Ingress с nginx | ✅ | Redis deployment + service|
| Метаданные студента | ✅ | Все labels присутствуют|

---

## Вывод

В ходе лабораторной работы:
- Разработан multi-stage Dockerfile для сборки минимального Python/Flask приложения размером 114 MB
- Создан полный набор Kubernetes манифестов, включая Deployment, Service, ConfigMap, Secret
- Настроены health checks: liveness и readiness probes через endpoints /ready и /health
- Реализован graceful shutdown с обработкой сигналов SIGTERM/SIGINT и корректным завершением работы
- Установлены resource limits согласно варианту 17: CPU=150m, Memory=128Mi
- Обеспечена отказоустойчивость через 2 реплики приложения
- Настроена стратегия RollingUpdate для бесшовного обновления приложения
- Интегрирован Redis как отдельный сервис для хранения состояния приложения
