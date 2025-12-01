# Лабораторная работа 02: Kubernetes базовый деплой

## Метаданные студента

- **ФИО:** Стельмашук Иван Александрович
- **Группа:** АС-63
- **№ студенческого (StudentID):** 220025
- **Email (учебный):** [AS006322@g.bstu.by](AS006322@g.bstu.by)
- **GitHub username:** KulibinI
- **Вариант №:** 20
- **Дата выполнения:** 25.11.2025
- **ОС (версия):** Windows 10 Home 22H2
- **Версия Docker Desktop:** 28.3.3
- **Версия kubectl:** 1.28.0
- **Версия Minikube:** 1.31.2

## Описание

Эта лабораторная работа реализует деплой HTTP-сервиса (Flask app из ЛР01) в Kubernetes с использованием Minikube. Включает Deployment, Service, Ingress, probes, ресурсы, rolling update, ConfigMap/Secret и PVC для данных. Также развернута зависимость Postgres как StatefulSet.

## Шаги сборки и запуска образа (из ЛР01)

1. Соберите образ: `docker build -t kulibini/myapp:stu-220025-v20 .`
2. Push в реестр: `docker push kulibini/myapp:stu-220025-v20`
3. Тестируйте локально: `docker run -p 8054:8054 -e DB_HOST=localhost kulibini/myapp:stu-220025-v20`

## Шаги деплоя и проверки в Kubernetes (Minikube)

1. Запустите Minikube: `minikube start`
2. Включите Ingress: `minikube addons enable ingress`
3. Примените манифесты: `kubectl apply -f k8s/`
4. Проверьте статус: `kubectl get all -n app20`
5. Добавьте в /etc/hosts: `127.0.0.1 web20.local`
6. Smoke-test: `curl http://web20.local/ping` (должен вернуть {"status": "healthy", ...})
7. Проверьте логи: `kubectl logs -f deployment/web20 -n app20` (ищите STU_* ENV и startup)
8. Тест graceful shutdown: `kubectl delete pod <pod-name> -n app20` (проверьте логи на shutdown)
9. Rolling update: Измените tag образа в deployment-app.yaml на новый (соберите/push новую версию), `kubectl apply -f k8s/deployment-app.yaml`. Наблюдайте: `kubectl get pods -n app20 -w` (замена без простоя).
10. Очистка: `kubectl delete -f k8s/`
