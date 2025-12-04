# Build stage
FROM python:3.11-slim AS builder

WORKDIR /app

COPY requirements.txt .

RUN pip install --user --no-cache-dir --no-warn-script-location -r requirements.txt

# Final stage
FROM python:3.11-slim

RUN adduser -D -u 10001 webuser && \
    apt-get update && \
    apt-get install -y --no-install-recommends wget=1.21.4-1+deb12u1 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder --chown=webuser:webuser /root/.local /home/webuser/.local
COPY --chown=webuser:webuser . .

ENV PATH=/home/webuser/.local/bin:$PATH
ENV PORT=5000

USER webuser

EXPOSE $PORT

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:$PORT/health || exit 1

CMD ["python", "app.py"]

LABEL org.bstu.student.fullname="Kazarenko Pavel Vladimirovich" \
      org.bstu.student.id="220008" \
      org.bstu.group="as-63" \
      org.bstu.variant="05" \
      org.bstu.course="RSIOT"
