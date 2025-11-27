FROM python:3.11-alpine AS builder

WORKDIR /app
COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt

FROM python:3.11-alpine

RUN adduser -D appuser && \
    apk add --no-cache curl=8.14.1-r2

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /root/.local /home/appuser/.local
COPY --chown=appuser:appuser main.py .

ENV PATH=/home/appuser/.local/bin:$PATH
ENV PYTHONPATH=/home/appuser/.local/lib/python3.11/site-packages

ENV APP_PORT=8073
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379
ENV STU_ID=220015
ENV STU_GROUP=AS-63
ENV STU_VARIANT=11

EXPOSE ${APP_PORT}
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${APP_PORT}/health || exit 1

LABEL org.bstu.student.fullname="Куликович Иван Сергеевич"
LABEL org.bstu.student.id="220015"
LABEL org.bstu.group="AS-63"
LABEL org.bstu.variant="11"
LABEL org.bstu.course="RSIOT"

USER appuser

CMD ["python", "main.py"]