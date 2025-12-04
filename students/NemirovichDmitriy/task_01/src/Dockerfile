# Базовый multi-stage build без дополнительных оптимизаций
FROM node:18 AS builder

WORKDIR /app

# Копируем package.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем исходники
COPY *.js .


# Финальный образ
FROM node:18-slim

# Метаданные студента (LABEL)
LABEL org.bstu.student.fullname="Немирович Дмитрий Александрович"
LABEL org.bstu.student.id="220050"
LABEL org.bstu.group="АС-64"
LABEL org.bstu.variant="37"
LABEL org.bstu.course="RSiOT"

WORKDIR /app

# Копируем из builder
COPY --from=builder /app /app

# Переменные окружения по умолчанию
ENV PORT=8001
ENV STU_ID=220050
ENV STU_GROUP=АС-64
ENV STU_VARIANT=37

# EXPOSE порт
EXPOSE 8001

# HEALTHCHECK
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:8001/ready', (r) => { process.exit(r.statusCode === 200 ? 0 : 1); }).on('error', () => process.exit(1));"

# Непривилегированный пользователь
USER 65532:65532

# Запуск
CMD ["node", "app.js"]
