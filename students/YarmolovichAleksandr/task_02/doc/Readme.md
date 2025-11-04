# Министерство образования Республики Беларусь

<p align="center">Учреждение образования</p>
<p align="center">“Брестский Государственный технический университет”</p>
<p align="center">Кафедра ИИТ</p>
<br><br><br><br><br><br>
<p align="center"><strong>Лабораторная работа №2</strong></p>
<p align="center"><strong>По дисциплине:</strong> “Распределенные системы и облачные технологии”</p>
<p align="center"><strong>Тема:</strong>Kubernetes: базовый деплой</p>
<br><br><br><br><br><br>
<p align="right"><strong>Выполнил:</strong></p>
<p align="right">Студент 4 курса</p>
<p align="right">Группы AC-63</p>
<p align="right">Ярмолович Александр Сергеевич</p>
<p align="right"><strong>Проверил:</strong></p>
<p align="right">Несюк А.Н.</p>
<br><br><br><br><br>
<p align="center"><strong>Брест 2025</strong></p>

# «Метаданные студента»

- ФИО - Ярмолович Александр Сергеевич
- Группа - АС-63
- № студенческого/зачетной книжки (StudentID) - 220029
- Email (учебный) -as006326@g.bstu.by 
- GitHub username - yarmolov
- Вариант № - 24
- Дата выполнения - 04.11.2025
- ОС (версия), версия Docker Desktop/Engine - Windows 10, Docker version 28.4.0, kubectl 1.30.5

# RSOT Проект

Это пример минимального HTTP-сервиса и набор Kubernetes-манифестов для лабораторной работы. Включает:

* Dockerfile (multi-stage, финальный образ ≤150 MB, не root)
* Приложение: минимальный HTTP-сервер с эндпоинтами `/`, `/healthz`, `/ready`.
* Логи запуска, остановки и корректного завершения
* Манифесты Kubernetes: Deployment (RollingUpdate + ресурсы), Service, ConfigMap, Secret, (опционально) PVC
* Настройки liveness/readiness probes
* Инструкции локального тестирования (kubectl)

## 1. Запуск образа

```bash
- docker compose up -d --build
```

## 2. Развертывание HTTP-сервиса в Kubernetes

### 2.1 Проверка текущего контекста Kubernetes

```bash
- kubectl config current-context
- kubectl get nodes
```

Смотрим с каким класстером мы работаем и доступны ли ноды.

### 2.2 Сборка Docker-образа

```bash
- docker build -t http-service:local .
```

Создаём локальный образ сервиса.

### 2.3 Установка Ingress NGINX

```bash
- kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.11.1/deploy/static/provider/cloud/deploy.yaml
```

Установка NGINX.

### 2.4 Создание Namespace

```bash
kubectl apply -f k8s/namespace.yaml
```

Создание отдельного пространства имён.

### 2.5 Применение манифестов приложения

```bash
kubectl apply -f k8s/web-configmap.yaml
kubectl apply -f k8s/postgres-secret.yaml
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/web-deployment.yaml
kubectl apply -f k8s/web-0service.yaml
```

Создайте все необходимые ресурсы Kubernetes в нашем namespace.

### 2.6 Применение Ingress

```bash
kubectl apply -f k8s/ingress.yaml
```

Настройка доступа к приложению через Ingress.

### 2.7 Добавление записи в /etc/hosts

```bash
sudo sh -c "echo '127.0.0.1 web24.local' >> /etc/hosts"
```

Добавляет локальную запись для доступа к нашему приложению.

### 2.8 Просмотр работоспособности подов

```bash
kubectl get pods -n app24
```

## 3. Подтвердить приложение

### 3.1 Liveness

```bash
curl http://web24.local/healthz
```

Content           : OK

### 3.2 Readiness

```bash
curl http://web24.local/ready
```

Content           : READY

### 3.3 Posgres

```bash
http://web24.local/visit
```

Студент: 24, Группа: feis, Вариант: v24, Кол-во визитов: 4

## 4. Просмотр логов

### 4.1 Просмотр логов бд

```bash
kubectl logs db-74d5cd88bd-s5ttt -n app24
```

### 4.2 Просмотр логов аpp

```bash
kubectl logs <app-name> -n app24
```

## 5. Очистка ресурсов

```bash
kubectl delete -f k8s/ingress.yaml
kubectl delete -f k8s/service.yaml
kubectl delete -f k8s/deployment.yaml
kubectl delete -f k8s/pvc.yaml
kubectl delete -f k8s/secret.yaml
kubectl delete -f k8s/configmap.yaml
kubectl delete -f k8s/namespace.yaml
```

## 6. Краткое описание проделанных действий

* Deployment — создан с типом стратегии RollingUpdate, добавлены ресурсные лимиты и запросы (resources.limits/requests) для контейнера.
* Service — настроен для доступа к приложению.
* ConfigMap и Secret — используются для передачи конфигурации и чувствительных данных в контейнер.
* Добавлен PersistentVolumeClaim и подключён volume для хранения данных приложения.
* Настроены livenessProbe и readinessProbe (HTTP), проверена их корректная работа.
* Подготовлены инструкции для локального тестирования:
1. Создание локального кластера;
2. Применение всех манифестов;
3. Проверка статуса Pod;
4. Выполнение smoke-теста через HTTP-запрос к эндпоинту приложения.