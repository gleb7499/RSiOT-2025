# ---- Builder stage ----
FROM python:3.11-alpine AS builder

ARG PIP_NO_CACHE_DIR=1
ENV PYTHONDONTWRITEBYTECODE=1 \
  PYTHONUNBUFFERED=1

RUN apk add --no-cache build-base

WORKDIR /app

COPY requirements.txt ./
# --no-cache-dir для pip исключает кэш, сокращая размер промежуточного слоя
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

# ---- Final stage ----
FROM python:3.11-alpine AS final

LABEL org.bstu.student.fullname="Логинов Глеб Олегович" \
      org.bstu.student.id="220018" \
      org.bstu.group="АС-63" \
      org.bstu.variant="14" \
      org.bstu.course="RSIOT"

ENV STU_ID=220018 \
  STU_GROUP=АС-63 \
  STU_VARIANT=14 \
  APP_PORT=8062 \
  APP_HOST=0.0.0.0

RUN adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder /install /usr/local

COPY src/app.py ./

EXPOSE 8062

RUN apk add --no-cache wget \
 && echo 'hosts: files dns' > /etc/nsswitch.conf

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:${APP_PORT}/healthz || exit 1

USER 10001

ENTRYPOINT ["python", "app.py"]
