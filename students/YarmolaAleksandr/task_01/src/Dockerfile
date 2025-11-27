FROM python:3.12-alpine AS builder

WORKDIR /app

COPY app/requirements.txt .

RUN pip install --prefix=/install --no-cache-dir -r requirements.txt


FROM python:3.12-alpine

WORKDIR /app

COPY --from=builder /install /usr/local

COPY app/ .

RUN adduser -u 65532 -D appuser

USER appuser

EXPOSE 8043

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8043/live').read()" || exit 1

LABEL org.bstu.student.fullname="Yarmola Aleksandr Olegovich" \
      org.bstu.student.id="220028" \
      org.bstu.group="as-63" \
      org.bstu.variant="23" \
      org.bstu.course="RSIOT"

CMD ["gunicorn", "-b", "0.0.0.0:8043", "--graceful-timeout", "30", "--timeout", "60", "-w", "2", "main:app"]
