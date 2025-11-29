# Stage 1: Builder
FROM python:3.11-slim as builder

WORKDIR /build

COPY requirements.txt .

RUN pip install --no-cache-dir --user -r requirements.txt


# Stage 2: Runtime
FROM python:3.11-slim

# Метаданные студента
LABEL org.bstu.student.fullname="Попов Алексей Сергеевич"
LABEL org.bstu.student.id="220051"
LABEL org.bstu.group="АС-64"
LABEL org.bstu.variant="38"
LABEL org.bstu.course="RSiOT"

WORKDIR /app

# Копируем зависимости из builder stage
COPY --from=builder /root/.local /home/appuser/.local

# Копируем исходный код
COPY app.py .

# Создаем непривилегированного пользователя
RUN useradd -m -u 10001 appuser && \
    chown -R appuser:appuser /app

# Переключаемся на непривилегированного пользователя
USER 10001

# Добавляем .local/bin в PATH
ENV PATH=/home/appuser/.local/bin:$PATH

# Переменные окружения
ENV STU_ID=220051
ENV STU_GROUP=АС-64
ENV STU_VARIANT=38

# Порт приложения
EXPOSE 8002

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8002/live')" || exit 1

# Запуск приложения
CMD ["python", "app.py"]
