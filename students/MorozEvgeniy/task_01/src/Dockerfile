FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY main.go ./

RUN CGO_ENABLED=0 go build -o /server main.go

FROM alpine:3.18

LABEL org.bstu.student.fullname="Мороз Евгений Васильевич"
LABEL org.bstu.student.id="220019"
LABEL org.bstu.group="АС-63"
LABEL org.bstu.variant="15"
LABEL org.bstu.course="RSIOT"

WORKDIR /

COPY --from=builder /server /server

RUN apk --no-cache add curl=8.4.0-r0

EXPOSE 8063

RUN addgroup -S appgroup && adduser -S appuser -G appgroup -u 65532

USER appuser

HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8063/ping || exit 1

ENTRYPOINT ["/server"]