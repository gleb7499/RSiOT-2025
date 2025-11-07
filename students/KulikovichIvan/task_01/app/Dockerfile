FROM python:3.11-slim as builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY app/requirements.txt .
RUN pip install --no-cache-dir --target=/app/dependencies -r requirements.txt

FROM python:3.11-slim

RUN groupadd -r appuser && useradd -r -g appuser appuser

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/dependencies /app/dependencies
COPY app/ .

ENV PYTHONPATH=/app/dependencies:$PYTHONPATH
ENV PATH=/app/dependencies/bin:$PATH

ENV APP_PORT=8000
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379

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