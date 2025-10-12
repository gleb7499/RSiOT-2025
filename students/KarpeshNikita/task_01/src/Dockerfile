FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY src/server.go .
RUN apk add --no-cache git && \
    go mod init server && \
    go get github.com/lib/pq && \
    go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server server.go

FROM alpine:3.20
RUN apk add --no-cache postgresql-client
RUN adduser -D -u 10001 nonroot
WORKDIR /app
COPY --from=builder /app/server .
RUN chown nonroot:nonroot /app/server
USER nonroot
ENV PORT=8092
LABEL org.bstu.student.fullname="Karpesh Nikita Petrovich" \
      org.bstu.student.id="220009" \
      org.bstu.group="as-63" \
      org.bstu.variant="06" \
      org.bstu.course="RSIOT"
EXPOSE $PORT
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:$PORT/health || exit 1
CMD ["./server"]