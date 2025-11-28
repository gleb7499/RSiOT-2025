FROM python:3.12-slim AS builder

WORKDIR /app

COPY requirements.txt .

RUN pip install --prefix=/install --no-cache-dir -r requirements.txt

FROM python:3.12-slim

WORKDIR /app

LABEL org.bstu.student.fullname="Иванюк Иван Александрович" \
      org.bstu.student.id="220041" \
      org.bstu.group="АС-64" \
      org.bstu.variant="8" \
      org.bstu.course="RSIOT"

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl=8.5.0-* && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

RUN groupadd -r appuser && useradd -r -g appuser appuser && \
    chown -R appuser:appuser /app

COPY --from=builder /install /usr/local

COPY app.py .

USER appuser

EXPOSE 8094

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8094/live || exit 1

CMD ["python", "app.py"]