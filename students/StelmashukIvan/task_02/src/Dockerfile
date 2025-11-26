FROM python:3.12-slim AS builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --prefix=/install --no-cache-dir -r requirements.txt

FROM python:3.12-slim
WORKDIR /app

LABEL org.bstu.student.fullname="Стельмашук Иван Александрович"
LABEL org.bstu.student.id="220025"
LABEL org.bstu.group="АС-63"
LABEL org.bstu.variant="20"
LABEL org.bstu.course="RSIOT"

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl=* && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

COPY --from=builder /install /usr/local
COPY app.py .

RUN groupadd -r appuser && useradd -r -g appuser -u 10001 appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8054

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8054/ping || exit 1

CMD ["python", "app.py"]