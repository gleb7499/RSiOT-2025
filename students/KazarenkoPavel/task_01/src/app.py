import os
import signal
import sys
import time
from flask import Flask, jsonify
import redis
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Новые переменные окружения
STU_ID = os.getenv('STU_ID', '220008')
STU_GROUP = os.getenv('STU_GROUP', 'as-63')
STU_VARIANT = os.getenv('STU_VARIANT', '05')

# Логируем при старте
logger.info(f"Student: {STU_ID}, Group: {STU_GROUP}, Variant: {STU_VARIANT}")

# Конфигурация
REDIS_HOST = os.getenv('REDIS_HOST', 'redis-as-63-220008-v05')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
APP_PORT = int(os.getenv('APP_PORT', 5000))
REDIS_PREFIX = f"stu:{STU_ID}:v{STU_VARIANT}"

redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)

def shutdown_handler(signum, frame):
    logger.info(f"Received signal {signum}, initiating graceful shutdown...")
    sys.exit(0)

signal.signal(signal.SIGTERM, shutdown_handler)
signal.signal(signal.SIGINT, shutdown_handler)

@app.route('/')
def hello():
    return jsonify({
        'message': 'Hello from Flask!',
        'status': 'success',
        'student': STU_ID,
        'variant': STU_VARIANT
    })

@app.route('/health')
def health():
    try:
        redis_client.ping()
        return jsonify({
            'status': 'healthy',
            'redis': 'connected',
            'student': STU_ID,
            'timestamp': time.time()
        }), 200
    except Exception as e:
        return jsonify({
            'status': 'unhealthy',
            'redis': 'disconnected',
            'error': str(e)
        }), 503

@app.route('/visit')
def visit_count():
    try:
        count = redis_client.incr(f'{REDIS_PREFIX}:visit_count')
        return jsonify({
            'visit_count': count,
            'message': f'This is visit number {count}',
            'redis_prefix': REDIS_PREFIX
        })
    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    logger.info("Starting Flask application...")
    app.run(host='0.0.0.0', port=APP_PORT, debug=False)
