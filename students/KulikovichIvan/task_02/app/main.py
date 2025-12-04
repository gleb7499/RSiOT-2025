import os
import signal
import sys
import time
from flask import Flask
import redis
import logging

REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
APP_PORT = int(os.getenv('APP_PORT', 8073))
STU_ID = os.getenv('STU_ID', 'default_id')
STU_GROUP = os.getenv('STU_GROUP', 'default_group')
STU_VARIANT = os.getenv('STU_VARIANT', 'default_variant')

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

logger.info(f"Запуск приложения. StudentID: {STU_ID}, Group: {STU_GROUP}, Variant: {STU_VARIANT}")
logger.info(f"Конфигурация: Redis {REDIS_HOST}:{REDIS_PORT}, Port: {APP_PORT}")

app = Flask(__name__)
cache = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True, socket_connect_timeout=5, socket_timeout=5)
KEY_PREFIX = f"stu:{STU_ID}:v{STU_VARIANT}"

def get_hit_count():
    """Получает и увеличивает счетчик обращений в Redis."""
    key = f"{KEY_PREFIX}:hits"
    try:
        return cache.incr(key)
    except redis.exceptions.RedisError as e:
        logger.error(f"Redis error: {e}")
        return "Ошибка подключения к Redis"

@app.route('/')
def hello():
    """Основной маршрут, возвращающий счетчик посещений."""
    count = get_hit_count()
    return f'<h1>Hello, Docker!</h1><p>Просмотр страницы №: {count}</p><p>Student: {STU_GROUP} - {STU_ID} - v{STU_VARIANT}</p>'

@app.route('/health')
def health():
    """Маршрут для healthcheck."""
    try:
        if cache.ping():
            return "OK", 200
    except redis.exceptions.RedisError:
        pass
    return "Service Unavailable", 503

def graceful_shutdown(signum, frame):
    """Обработчик сигналов для graceful shutdown."""
    logger.info(f"Получен сигнал {signum}. Начинаем graceful shutdown...")
    logger.info("Приложение корректно завершило работу.")
    sys.exit(0)

signal.signal(signal.SIGTERM, graceful_shutdown)
signal.signal(signal.SIGINT, graceful_shutdown)

if __name__ == '__main__':
    logger.info(f"Запуск приложения. StudentID: {STU_ID}, Group: {STU_GROUP}, Variant: {STU_VARIANT}")
    app.run(host='0.0.0.0', port=APP_PORT, debug=False)