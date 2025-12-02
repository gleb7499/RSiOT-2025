# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go

# Final stage
FROM alpine:3.19

# Метаданные студента
LABEL org.bstu.student.fullname="Кульбеда Кирилл Александрович"
LABEL org.bstu.student.id="220016"
LABEL org.bstu.group="АС-63"
LABEL org.bstu.variant="12"
LABEL org.bstu.course="RSiOT"

# Создание непривилегированного пользователя
RUN adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder /build/app .

# Переключение на непривилегированного пользователя
USER 10001

EXPOSE 8074

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8074/ready || exit 1

CMD ["./app"]
