FROM golang:1.22-alpine AS builder
WORKDIR /src

COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w -X main.buildVersion=${BUILD_VERSION:-v21}" \
    -o /out/app ./cmd/server

FROM alpine:3.20

ARG BSTU_FULLNAME="Тунчик Антон Дмитриевич"
ARG BSTU_STUDENT_ID="220026"
ARG BSTU_GROUP="AS-63"
ARG BSTU_VARIANT="21"
ARG BSTU_COURSE="4"

LABEL org.bstu.student.fullname="${BSTU_FULLNAME}" \
      org.bstu.student.id="${BSTU_STUDENT_ID}" \
      org.bstu.group="${BSTU_GROUP}" \
      org.bstu.variant="${BSTU_VARIANT}" \
      org.bstu.course="${BSTU_COURSE}"

RUN adduser -D -u 65532 appuser

WORKDIR /app
COPY --from=builder /out/app /app/app

ENV PORT=8041
ENV REDIS_ADDR=redis:6379
ENV STU_VARIANT=21

EXPOSE 8041

HEALTHCHECK --interval=15s --timeout=2s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:${PORT}/health || exit 1

USER 65532
ENTRYPOINT ["/app/app"]
