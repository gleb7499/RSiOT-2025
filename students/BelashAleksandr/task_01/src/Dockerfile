# Базовая сборка (будет использоваться для установки зависимостей)
FROM node:18-alpine AS builder

WORKDIR /app

# Копируем package.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm ci --only=production

# Финальный образ
FROM node:18-alpine

# Метаданные студента
LABEL org.bstu.student.fullname="Белаш Александр Олегович"
LABEL org.bstu.student.id="220031"
LABEL org.bstu.group="АС-64"
LABEL org.bstu.variant="25"
LABEL org.bstu.course="RSiOT"

# Создаем непривилегированного пользователя
RUN addgroup -g 65532 -S appgroup && \
    adduser -u 65532 -S appuser -G appgroup

WORKDIR /app

# Копируем зависимости из builder
COPY --from=builder --chown=appuser:appgroup /app/node_modules ./node_modules

# Копируем исходный код
COPY --chown=appuser:appgroup server.js .
COPY --chown=appuser:appgroup package.json .

# Переключаемся на непривилегированного пользователя
USER 65532

# Открываем порт
EXPOSE 8031

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:8031/ping', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

# Запуск приложения
CMD ["node", "server.js"]
