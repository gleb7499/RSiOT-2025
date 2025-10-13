FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY src/server.go .
RUN apk add --no-cache git && \
    go mod init server && \
    go get github.com/redis/go-redis/v9 && \
    go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server server.go

FROM alpine:3.20
RUN apk add --no-cache redis
RUN adduser -D -u 65532 nonroot
WORKDIR /app
COPY --from=builder /app/server .
RUN chown nonroot:nonroot /app/server
USER nonroot
ENV PORT=8071
LABEL org.bstu.student.fullname="Kolodich Maksim Pavlovich" \
      org.bstu.student.id="220013" \
      org.bstu.group="as-63" \
      org.bstu.variant="09" \
      org.bstu.course="RSIOT"
EXPOSE $PORT
HEALTHCHECK --interval=15s --timeout=3s --start-period=3s --retries=2 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:$PORT/healthz || exit 1
CMD ["./server"]
