FROM golang:1.25 AS builder

WORKDIR /app

COPY go.* ./
RUN test -f go.sum || : > go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server

FROM alpine:3.20
LABEL org.bstu.student.fullname="Nikita Horkauchuk" \
      org.bstu.student.id="220038" \
      org.bstu.group="as-64" \
      org.bstu.variant="6" \
      org.bstu.course="AIPS"

RUN adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /app/server /app/server

RUN apk add --no-cache ca-certificates curl

EXPOSE 8092

USER 10001

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -fsS http://localhost:8092/health || exit 1

ENTRYPOINT ["/app/server"]
