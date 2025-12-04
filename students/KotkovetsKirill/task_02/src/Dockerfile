# Базовый multi-stage Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o web35 main.go

# Финальный образ
FROM alpine:3.18

# Создаем пользователя appuser
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Копируем бинарник
COPY --from=builder /build/web35 .

# Метаданные студента
LABEL org.bstu.student.fullname="Котковец Кирилл Викторович"
LABEL org.bstu.student.id="220044"
LABEL org.bstu.group="АС-64"
LABEL org.bstu.variant="35"
LABEL org.bstu.course="RSiOT"

# Переключаемся на non-root пользователя
USER appuser

EXPOSE 8013

CMD ["./web35"]
