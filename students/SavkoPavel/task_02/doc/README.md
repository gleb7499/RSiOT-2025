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
<p align="right">Савко П.С.</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А.Н.</p>

<p align="center"><strong>Брест 2025</strong></p>

---

## Цель работы

Научиться готовить Kubernetes-манифесты для простого HTTP-сервиса (Deployment + Service). Подготовить конфигурацию через ConfigMap/Secret и смонтировать volume для данных при необходимости.

---

### Вариант №18

| Параметр            | Значение                  |
|---------------------|---------------------------|
| Namespace           | `app18`                   |
| Имя приложения      | `web18`                   |
| Количество реплик   | 3                         |
| Порт приложения     | 8052                      |
| Ingress Class       | `nginx`                   |
| CPU limit           | 200m                      |
| Memory limit        | 192Mi                     |
| Healthcheck endpoint| `/live`                   |
| Номер варианта      | 18                        |

---

### Ход выполнения работы

### 1. Структура проекта

├── Dockerfile # Multi-stage сборка Go-приложения
├── go.mod, go.sum # Go зависимости
├── src/main.go # Исходный код HTTP-сервиса
└── k8s/ # Kubernetes манифесты
├── namespace.yaml # Namespace app18
├── configmap.yaml # Конфигурация приложения
├── secret.yaml # Секреты (БД)
├── deployment.yaml # Deployment с 3 репликами
├── service.yaml # ClusterIP Service
└── ingress.yaml # Ingress с nginx контроллером

---

### 2. Последовательность запуска

#### Проверка версий

docker --version
kubectl version --client

#### Проверка кластера Kubernetes

kubectl cluster-info
kubectl get nodes

#### Установка Ingress контроллера

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=120s

---

### 3. Сборка Docker образа

docker build -t web18-app:v1 .
docker images web18-app:v1 --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

---

### 4. Развертывание в Kubernetes

kubectl apply -f k8s\namespace.yaml
kubectl apply -f k8s\configmap.yaml
kubectl apply -f k8s\secret.yaml
kubectl apply -f k8s\deployment.yaml
kubectl apply -f k8s\service.yaml
kubectl apply -f k8s\ingress.yaml

#### Проверка создания ресурсов

kubectl get all -n app18
kubectl get configmap,secret,ingress -n app18

---

### 5. Тестирование приложения

#### Port-forward для локального доступа

kubectl port-forward -n app18 svc/web18-service 8052:8052 &

#### Тестирование endpoints

curl http://localhost:8052/live
Ответ: ok

curl http://localhost:8052/
Ответ: JSON с информацией о сервисе

---

### 6. Проверка требований

| Критерий | Выполнено | Комментарий |
|-----------|-----------|--------------|
| Multi-stage Docker build | ✅ | Образ 29.4 MB|
| Размер образа ≤ 150 MB | ✅ | 29.4 MB < 150 MB|
| Non-root пользователь | ✅ | UID: 10001|
| Resource limits | ✅ | CPU=200m, Memory=192Mi|
| Health endpoints | ✅ | /live возвращает ok |
| Liveness/Readiness probes | ✅ | Настроены в deployment|
| Graceful shutdown | ✅ | Обработка SIGTERM/SIGINT|
| 3 реплики | ✅ | Deployment с 3 подами|
| RollingUpdate стратегия | ✅ | Настроена в deployment|
| ConfigMap/Secret | ✅ | Разделение конфигурации|
| Service (ClusterIP) | ✅ | Сервис создан и работает|
| Ingress с nginx | ✅ | Ingress настроен|
| Метаданные студента | ✅ | Все labels присутствуют|

---

## Вывод

В ходе лабораторной работы:
- Разработан multi-stage Dockerfile для сборки минимального Go-приложения размером 29.4 MB
- Создан полный набор Kubernetes манифестов
- Настроены health checks: liveness и readiness probes через endpoint /live
- Реализован graceful shutdown с логированием получения сигналов и корректным завершением работы
- Установлены resource limits согласно варианту 18: CPU=200m, Memory=192Mi
