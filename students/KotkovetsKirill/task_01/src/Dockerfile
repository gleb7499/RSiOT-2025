# Multi-stage build для уменьшения размера образа
FROM python:3.11-slim as builder

WORKDIR /app

# Копируем requirements и устанавливаем зависимости
COPY requirements.txt .
RUN pip install --no-cache-dir --user -r requirements.txt

# Финальный образ
FROM python:3.11-slim

# Метаданные студента
LABEL org.bstu.student.fullname="Котковец Кирилл Викторович"
LABEL org.bstu.student.id="220044"
LABEL org.bstu.group="АС-64"
LABEL org.bstu.variant="35"
LABEL org.bstu.course="RSiOT"

WORKDIR /app

# Копируем установленные пакеты из builder
COPY --from=builder /root/.local /root/.local
ENV PATH=/root/.local/bin:$PATH

# Копируем исходный код
COPY app.py .

# Переменные окружения
ENV PORT=8013
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379
ENV REDIS_PREFIX=stu:220044:v35
ENV STU_ID=220044
ENV STU_GROUP=АС-64
ENV STU_VARIANT=35

# Непривилегированный пользователь
USER 65532

# Expose порт
EXPOSE 8013

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8013/ping').read()" || exit 1

# Запуск приложения
CMD ["python", "app.py"]
