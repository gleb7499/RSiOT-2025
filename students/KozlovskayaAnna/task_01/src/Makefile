# Makefile для лабораторной работы (Вариант 8)

SLUG=as63-220012-v8
COMPOSE_PROJECT_NAME?=$(SLUG)
IMAGE=annkrq/lab1-v8:stu-220012-v8
PORT=8094
APP_CONTAINER=app-$(SLUG)
DB_CONTAINER=db-$(SLUG)
COMPOSE=docker compose --project-name $(COMPOSE_PROJECT_NAME)

POWERSHELL=powershell -NoLogo -NoProfile -Command "[Console]::OutputEncoding=[System.Text.Encoding]::UTF8; $$OutputEncoding = New-Object System.Text.UTF8Encoding;"

.DEFAULT_GOAL := help

.PHONY: help build up start logs app-log db-log ps stop down rebuild clean prune health kill-term graceful-test shell image-size

help:
	@echo "Available targets:"
	@echo "  build          - Build docker image"
	@echo "  up | start     - Start services (app + db)"
	@echo "  health         - Check /healthz endpoint"
	@echo "  logs           - Application streaming logs"
	@echo "  stop           - Stop application"
	@echo "  kill-term      - Send SIGTERM to app (graceful)"
	@echo "  graceful-test  - SIGTERM + wait for shutdown"
	@echo "  down           - Stop and remove services"
	@echo "  image-size     - Image size"
	@echo "  clean | prune  - Prune docker resources"
	@echo "Example: make build && make up && make health"

build:
	@echo "==> Building image $(IMAGE)"
	@docker build -t $(IMAGE) -f Dockerfile .

up start:
	@echo "==> Starting stack (app + db) (project=$(COMPOSE_PROJECT_NAME))"
	@$(COMPOSE) up -d --build
	@$(COMPOSE) ps

logs:
	@$(COMPOSE) logs -f app

app-log:
	@$(COMPOSE) logs app

db-log:
	@$(COMPOSE) logs db

ps:
	@$(COMPOSE) ps

stop:
	@echo "==> Stopping app container"
	@docker stop $(APP_CONTAINER) || true

kill-term:
	@echo "==> Sending SIGTERM to app (graceful)"
	@docker kill --signal=SIGTERM $(APP_CONTAINER)

graceful-test: kill-term
	@echo "Waiting for graceful shutdown..."
	@$(POWERSHELL) "$$ProgressPreference='SilentlyContinue'; for($$i=0; $$i -lt 30; $$i++){ $$logs = docker logs $(APP_CONTAINER) 2>&1; if($$logs -match 'Graceful shutdown complete.'){ Write-Host 'Graceful shutdown complete' -ForegroundColor Green; exit 0 }; Start-Sleep -Seconds 1 }; Write-Host 'Graceful shutdown timeout' -ForegroundColor Yellow; exit 1"

health:
	@echo "==> GET /healthz"
	@$(POWERSHELL) "try { (Invoke-RestMethod -Uri 'http://localhost:$(PORT)/healthz' -Method GET) | ConvertTo-Json -Depth 5 } catch { Write-Host 'Health endpoint unavailable' -ForegroundColor Red; exit 1 }"

rebuild:
	@echo "==> Rebuild without cache"
	@docker build --no-cache -t $(IMAGE) -f Dockerfile .

image-size:
	@docker images annkrq/lab1-v8 --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | findstr stu-220012-v8 || echo Image tag not found

shell:
	@echo "==> Opening shell in app container"
	@docker exec -it $(APP_CONTAINER) sh || true

down:
	@echo "==> Stopping and removing stack (project=$(COMPOSE_PROJECT_NAME))"
	@$(COMPOSE) down

clean prune:
	@echo "==> Pruning unused docker data"
	@docker system prune -f
