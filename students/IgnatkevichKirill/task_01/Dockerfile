# ---------- Stage 1: build ----------
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY src/ ./src/
WORKDIR /app/src
RUN go mod tidy && go build -o /server main.go

# ---------- Stage 2: final image ----------
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /server .
EXPOSE 8071
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD [ "/busybox/wget", "-q", "--spider", "http://localhost:8071/healthz" ]
ENV STU_ID=220042 STU_GROUP=AS-64 STU_VARIANT=9
LABEL org.bstu.student.fullname="Игнаткевич Кирилл Сергеевич" \
      org.bstu.student.id="220042" \
      org.bstu.group="AS-64" \
      org.bstu.variant="9" \
      org.bstu.course="RSIOT"
ENTRYPOINT ["/app/server"]