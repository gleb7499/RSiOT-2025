# ---- Builder stage ----
FROM python:3.11-alpine AS builder

ARG PIP_NO_CACHE_DIR=1
ENV PYTHONDONTWRITEBYTECODE=1 \
  PYTHONUNBUFFERED=1

RUN apk add --no-cache build-base=0.5-r3

WORKDIR /app

COPY requirements.txt ./
# --no-cache-dir для pip исключает кэш, сокращая размер промежуточного слоя
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

# ---- Final stage ----
FROM python:3.11-alpine AS final

LABEL org.bstu.student.fullname="Козловская Анна Геннадьевна" \
      org.bstu.student.id="220012" \
      org.bstu.group="АС-63" \
      org.bstu.variant="8" \
      org.bstu.course="RSiOT"

ENV STU_ID=220012 \
  STU_GROUP=АС-63 \
  STU_VARIANT=8 \
  APP_PORT=8094 \
  APP_HOST=0.0.0.0 \
  PYTHONDONTWRITEBYTECODE=1 \
  PYTHONUNBUFFERED=1

RUN adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder /install /usr/local

COPY src/app.py ./

EXPOSE 8094

# Устанавливаем wget с фиксацией версии (совместимой с Alpine 3.22) и настраиваем nsswitch
RUN apk add --no-cache wget=1.25.0-r1 \
 && echo 'hosts: files dns' > /etc/nsswitch.conf

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-hsts -qO- http://127.0.0.1:${APP_PORT}/healthz || exit 1

USER 10001

ENTRYPOINT ["python", "app.py"]
