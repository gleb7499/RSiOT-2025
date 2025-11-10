# syntax=docker/dockerfile:1.7

# ========== Build stage ==========
FROM golang:1.23-alpine AS build

# Metadata labels (student)
LABEL org.bstu.student.fullname="Савко Павел Станиславович" \
      org.bstu.student.id="220023" \
      org.bstu.group="АС-63" \
      org.bstu.variant="18" \
      org.bstu.course="RSIOT"

WORKDIR /src

# Minimal build deps
# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

# Cache go modules first
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy sources
COPY . .

# Build with cache for go build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/app ./cmd/app

# ========== Runtime stage ==========
FROM gcr.io/distroless/static-debian12:nonroot

# Use an explicit non-root user with UID 10001 as per requirement
USER 10001:10001

WORKDIR /app

# Copy CA roots and binary
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/app /app/app

# Expose port (per variant: 8052)
EXPOSE 8052

# Healthcheck: self-probe via binary flag (no curl in distroless)
HEALTHCHECK --interval=10s --timeout=2s --retries=5 --start-period=5s \
  CMD ["/app/app", "-test.healthcheck"]

# Default envs (can be overridden)
ENV PORT=8052 \
    HOST=0.0.0.0 \
    STU_ID=220023 \
    STU_GROUP="АС-63" \
    STU_VARIANT=18 \
    PGHOST=db \
    PGPORT=5432 \
    PGUSER=stuuser \
    PGPASSWORD=stupass \
    PGDATABASE=app_220023_v18 \
    PGSSLMODE=disable

# Entry point
ENTRYPOINT ["/app/app"]