#!/usr/bin/env bash
set -euo pipefail

APP_URL="${APP_URL:-http://127.0.0.1:8041}"

echo "==> build & up"
docker compose build --no-cache
docker compose up -d

echo "==> wait for health"
for i in $(seq 1 30); do
  if curl -sf "${APP_URL}/health" >/dev/null; then
    echo "healthy"
    break
  fi
  sleep 1
done

echo "==> sample requests"
curl -s "${APP_URL}/" | jq .
curl -s "${APP_URL}/hit" | jq .

echo "==> stop app (SIGTERM) and watch logs"
docker compose stop app
docker compose logs app | tail -n 50

echo "OK"
