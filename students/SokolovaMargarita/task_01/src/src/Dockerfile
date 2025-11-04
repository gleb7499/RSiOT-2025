# syntax=docker/dockerfile:1


FROM node:20-alpine AS deps
WORKDIR /app
COPY package*.json ./
# Устанавливаем зависимости — кэшируется, если package*.json не менялись
RUN npm ci --omit=dev


FROM node:20-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
# Устанавливаем права на /app под UID 65532 (будет использовано в финальном образе)
RUN chown -R 65532:65532 /app


FROM node:20-alpine AS runtime
WORKDIR /app
# Копируем файлы из build и сразу задаём владельца 65532:65532
COPY --from=build --chown=65532:65532 /app /app

# Metadata (обязательные поля по заданию)
LABEL org.bstu.student.fullname="Соколова Маргарита Александровна" \
      org.bstu.student.id="220024" \
      org.bstu.group="АС-63" \
      org.bstu.variant="19" \
      org.bstu.course="RSIOT"

# Работаем от непривилегированного пользователя (UID: 65532)
USER 65532:65532

# Порт приложения внутри контейнера (мы слушаем на 8053)
EXPOSE 8053

# Healthcheck (используем встроенный wget в busybox)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD wget -qO- http://127.0.0.1:8053/healthz || exit 1

# Запуск
CMD ["node", "server.js"]
