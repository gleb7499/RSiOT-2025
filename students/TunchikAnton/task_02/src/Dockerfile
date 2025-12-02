# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -X main.buildVersion=v2.0.0-k8s" -o /app/server ./cmd/server

# Final stage
FROM alpine:3.18

RUN adduser -D -u 1000 appuser

WORKDIR /app
COPY --from=builder --chown=appuser:appuser /app/server /app/server

USER appuser

EXPOSE 8041

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8041/health || exit 1

CMD ["/app/server"]