FROM python:3.11-alpine as builder

WORKDIR /app
COPY requirements.txt .
RUN --mount=type=cache,target=/root/.cache/pip pip install --no-cache-dir --prefix=/install -r requirements.txt

FROM python:3.11-alpine

LABEL org.bstu.student.fullname="Кужир Владислав Витальевич"
LABEL org.bstu.student.id="220047"
LABEL org.bstu.group="AS-64"
LABEL org.bstu.variant="26"
LABEL org.bstu.course="RSIOT"

RUN adduser --disabled-password --gecos '' --uid 10001 appuser

WORKDIR /app

COPY --from=builder /install /usr/local
COPY app.py .

RUN apk add --no-cache curl

RUN chown -R appuser:appuser /app

ENV FLASK_APP=app.py
ENV FLASK_ENV=production
ENV PORT=8032
ENV STU_ID="220047"
ENV STU_GROUP="AC-64"
ENV STU_VARIANT="26"

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8032/health || exit 1

EXPOSE 8032

USER appuser

CMD ["python", "app.py"]