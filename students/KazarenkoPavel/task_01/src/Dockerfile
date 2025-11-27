# Build stage
FROM python:3.11-slim as builder

WORKDIR /app

# Копируем только файлы зависимостей для лучшего кэширования
COPY requirements.txt .

# Устанавливаем зависимости в виртуальное окружение
RUN pip install --user --no-cache-dir -r requirements.txt

# Final stage
FROM python:3.11-slim

# Устанавливаем curl для healthcheck
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Создаем не-root пользователя
RUN groupadd -r webuser && useradd -r -g webuser webuser

WORKDIR /app

# Копируем зависимости из builder stage
COPY --from=builder /root/.local /home/webuser/.local
COPY --chown=webuser:webuser . .

# Устанавливаем права для пользователя
RUN chown -R webuser:webuser /app

# Добавляем .local/bin в PATH
ENV PATH=/home/webuser/.local/bin:$PATH

# Переключаемся на non-root пользователя
USER webuser

# Экспортируем порт
EXPOSE 5000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:5000/health || exit 1

# Запускаем приложение
CMD ["python", "app.py"]

# Метаданные
LABEL org.bstu.student.fullname="Kazarenko Pavel Vladimirovich"
LABEL org.bstu.student.id="220008"
LABEL org.bstu.group="as-63"
LABEL org.bstu.variant="05"
LABEL org.bstu.course="RSIOT"
