# app.py
import os, signal, logging, time, sys
from flask import Flask, jsonify
import redis

# Настройка логирования для Kubernetes
logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout)
    ]
)
logger = logging.getLogger(__name__)

STU_ID = os.getenv('STU_ID', '220021')
STU_GROUP = os.getenv('STU_GROUP', 'ASOI 63')
STU_VARIANT = os.getenv('STU_VARIANT', '17')
REDIS_HOST = os.getenv('REDIS_HOST', 'web17-redis')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
REDIS_PREFIX = f"stu:{STU_ID}:v{STU_VARIANT}:entity:"

app = Flask(__name__)
redis_client = None
shutting_down = False

def get_redis():
    global redis_client
    if redis_client is None:
        try:
            redis_client = redis.Redis(
                host=REDIS_HOST,
                port=REDIS_PORT,
                decode_responses=True,
                socket_connect_timeout=5,
                retry_on_timeout=True
            )
            # Проверка соединения
            redis_client.ping()
            logger.info(f"Redis подключен: {REDIS_HOST}:{REDIS_PORT}")
        except Exception as e:
            logger.error(f"Ошибка подключения к Redis: {e}")
            redis_client = None
    return redis_client

def log_startup():
    logger.info("=" * 50)
    logger.info(f"Запуск приложения")
    logger.info(f"Студент: {STU_ID}, Группа: {STU_GROUP}, Вариант: {STU_VARIANT}")
    logger.info(f"Redis: {REDIS_HOST}:{REDIS_PORT}")
    logger.info(f"Префикс ключей: {REDIS_PREFIX}")
    logger.info("=" * 50)

@app.route('/ready')
def ready():
    """Readiness probe endpoint"""
    global shutting_down
    if shutting_down:
        return jsonify({"status": "shutting down"}), 503
    
    try:
        # Проверяем Redis только для readiness
        r = get_redis()
        if r:
            r.ping()
        return jsonify({
            "status": "ok",
            "student": STU_ID,
            "group": STU_GROUP,
            "variant": STU_VARIANT
        }), 200
    except Exception as e:
        logger.error(f"Readiness check failed: {e}")
        return jsonify({"status": "error", "message": str(e)}), 500

@app.route('/health')
def health():
    """Liveness probe endpoint"""
    if shutting_down:
        return jsonify({"status": "shutting down"}), 503
    return jsonify({"status": "healthy"}), 200

@app.route('/')
def index():
    logger.info('Запрос к главной странице')
    r = get_redis()
    key = REDIS_PREFIX + 'counter'
    try:
        count = r.incr(key) if r else None
        logger.info(f'Счетчик посещений: {count}')
        return jsonify({
            "message": "Веб-приложение студента",
            "student": STU_ID,
            "group": STU_GROUP,
            "variant": STU_VARIANT,
            "visits": count,
            "pod": os.getenv('HOSTNAME', 'unknown')
        })
    except Exception as e:
        logger.exception('Ошибка Redis')
        return jsonify({
            "message": "Приложение работает, но Redis недоступен",
            "student": STU_ID,
            "visits": None,
            "error": str(e)
        }), 200

@app.route('/metrics')
def metrics():
    """Простые метрики для мониторинга"""
    return jsonify({
        "status": "running",
        "student_id": STU_ID,
        "variant": STU_VARIANT
    }), 200

def shutdown_handler(signum, frame):
    global shutting_down, redis_client
    shutting_down = True
    logger.info(f'Получен сигнал {signum}, начало graceful shutdown')
    
    # Закрываем соединения
    try:
        if redis_client:
            redis_client.close()
            logger.info('Соединение с Redis закрыто')
    except Exception as e:
        logger.error(f'Ошибка при закрытии Redis: {e}')
    
    logger.info('Приложение завершено')
    time.sleep(1)
    sys.exit(0)

# Регистрация обработчиков сигналов
signal.signal(signal.SIGTERM, shutdown_handler)
signal.signal(signal.SIGINT, shutdown_handler)

if __name__ == '__main__':
    log_startup()
    port = int(os.getenv('PORT', 8051))
    app.run(host='0.0.0.0', port=port, debug=False)
