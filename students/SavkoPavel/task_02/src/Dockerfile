# Build stage
FROM golang:1.23-alpine AS builder

# Metadata labels
LABEL org.bstu.student.fullname="Савко Павел Станиславович" \
      org.bstu.student.id="220023" \
      org.bstu.group="АС-63" \
      org.bstu.variant="18" \
      org.bstu.course="RSIOT"

WORKDIR /app

# Copy go modules
COPY go.mod .

# Download dependencies
RUN go mod download

# Copy source code
COPY src/ ./src/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /main ./src

# Runtime stage
FROM alpine:3.19

# Metadata
LABEL org.bstu.student.fullname="Савко Павел Станиславович" \
      org.bstu.student.id="220023" \
      org.bstu.group="АС-63" \
      org.bstu.variant="18" \
      org.bstu.course="RSIOT"

# Install CA certificates and tzdata with pinned versions
RUN apk --no-cache add \
    ca-certificates=20230506-r0 \
    tzdata=2023d-r1 \
    && rm -rf /var/cache/apk/*

# Create non-root user and set up working directory in one RUN
RUN adduser -D -u 10001 appuser \
    && mkdir -p /home/appuser \
    && chown -R appuser:appuser /home/appuser

WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /main ./app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8052

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O- http://localhost:8052/live || exit 1

# Run the application
CMD ["./app"]