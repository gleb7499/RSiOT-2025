# syntax=docker/dockerfile:1
#########################
# Stage 1: build
#########################
FROM golang:1.25-alpine AS build
WORKDIR /app

# Кэш зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем
COPY . .
RUN go build -o server ./src

#########################
# Stage 2: runtime
#########################
FROM alpine:3.20

# Добавляем curl для healthcheck
RUN apk add --no-cache curl=8.2.1-r0

WORKDIR /app
USER 10001:10001

COPY --from=build /app/server .

EXPOSE 8044

LABEL org.bstu.student.fullname="Yarmolovich Alexander Sergeevich" \
      org.bstu.student.id="24" \
      org.bstu.group="feis" \
      org.bstu.variant="v24" \
      org.bstu.course="RSIOT"

ENV STU_ID=24 \
    STU_GROUP=feis \
    STU_VARIANT=v24 \
    PORT=8044

# Рабочий healthcheck
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8044/healthz || exit 1

ENTRYPOINT ["/app/server"]
