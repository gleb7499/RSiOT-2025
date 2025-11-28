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
    apt-get install -y --no-install-recommends curl && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

COPY --from=builder /install /usr/local

COPY app.py .

RUN groupadd -r appuser && useradd -r -g appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8094

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD ["/usr/bin/curl", "-fss", "http://localhost:8094/live"]

CMD ["python", "app.py"]