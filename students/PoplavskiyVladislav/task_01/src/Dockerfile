# ---- build stage ----
FROM python:3.11-slim AS builder
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY app/requirements.txt ./
RUN python -m venv /venv && /venv/bin/pip install --upgrade pip
RUN /venv/bin/pip wheel --wheel-dir=/wheels -r requirements.txt

# ---- final stage ----
FROM python:3.11-slim
ARG APP_UID=65532
RUN groupadd -g ${APP_UID} appgroup || true
RUN useradd -r -u ${APP_UID} -g appgroup -m -d /home/appuser -s /usr/sbin/nologin appuser || true
WORKDIR /app
COPY --from=builder /wheels /wheels
COPY app/requirements.txt ./
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    pip install --no-index --find-links=/wheels --no-cache-dir -r requirements.txt && \
    rm -rf /wheels /var/lib/apt/lists/* /root/.cache
COPY app /app
ENV FLASK_APP=app.py \
    FLASK_ENV=production \
    PORT=8051 \
    STU_ID=220021 \
    STU_GROUP="ASOI 63" \
    STU_VARIANT=17
EXPOSE 8051
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8051/ready')" || exit 1
USER ${APP_UID}
LABEL org.bstu.student.fullname="Vladislav Vladimirovich Poplavsky" \
      org.bstu.student.id="220021" \
      org.bstu.group="ASOI 63" \
      org.bstu.variant="17" \
      org.bstu.course="RSIOT"
CMD ["/usr/local/bin/python", "-u", "app.py"]
