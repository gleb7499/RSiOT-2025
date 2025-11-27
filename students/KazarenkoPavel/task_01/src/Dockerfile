FROM python:3.11-slim as builder

WORKDIR /app

COPY requirements.txt .

RUN pip install --user --no-cache-dir -r requirements.txt

FROM python:3.11-slim

RUN apt-get update && apt-get install -y --no-install-recommends curl \
    && groupadd -r webuser && useradd -r -g webuser webuser \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /root/.local /home/webuser/.local
COPY --chown=webuser:webuser . .

RUN chown -R webuser:webuser /app

ENV PATH=/home/webuser/.local/bin:$PATH

USER webuser

EXPOSE 5000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:5000/health')" || exit 1

CMD ["python", "app.py"]

LABEL org.bstu.student.fullname="Kazarenko Pavel Vladimirovich"
LABEL org.bstu.student.id="220008"
LABEL org.bstu.group="as-63"
LABEL org.bstu.variant="05"
LABEL org.bstu.course="RSIOT"
