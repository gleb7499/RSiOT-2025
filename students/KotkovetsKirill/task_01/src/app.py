import os
import signal
import sys
import time
from flask import Flask, jsonify
import redis

app = Flask(__name__)

# Конфигурация через переменные окружения
REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
REDIS_PREFIX = os.getenv('REDIS_PREFIX', 'stu:220044:v35')
STU_ID = os.getenv('STU_ID', '220044')
STU_GROUP = os.getenv('STU_GROUP', 'АС-64')
STU_VARIANT = os.getenv('STU_VARIANT', '35')
PORT = int(os.getenv('PORT', 8013))

# Подключение к Redis
try:
    redis_client = redis.Redis(
        host=REDIS_HOST,
        port=REDIS_PORT,
        decode_responses=True
    )
    redis_client.ping()
    print(f"[INFO] Connected to Redis at {REDIS_HOST}:{REDIS_PORT}")
except Exception as e:
    print(f"[ERROR] Failed to connect to Redis: {e}")
    redis_client = None

# Логирование метаданных при старте
print(f"[INFO] Starting application...")
print(f"[INFO] Student ID: {STU_ID}")
print(f"[INFO] Group: {STU_GROUP}")
print(f"[INFO] Variant: {STU_VARIANT}")
print(f"[INFO] Redis Prefix: {REDIS_PREFIX}")

@app.route('/ping', methods=['GET'])
def ping():
    """Health check endpoint"""
    try:
        if redis_client:
            redis_client.ping()
            status = "ok"
        else:
            status = "degraded"
        
        return jsonify({
            "status": status,
            "student_id": STU_ID,
            "variant": STU_VARIANT
        }), 200
    except Exception as e:
        return jsonify({
            "status": "error",
            "message": str(e)
        }), 503

@app.route('/', methods=['GET'])
def index():
    """Root endpoint"""
    # Инкремент счетчика в Redis
    if redis_client:
        key = f"{REDIS_PREFIX}:visits"
        redis_client.incr(key)
        visits = redis_client.get(key)
    else:
        visits = "N/A"
    
    return jsonify({
        "message": "Flask app is running",
        "student": STU_ID,
        "group": STU_GROUP,
        "variant": STU_VARIANT,
        "visits": visits
    })

# Graceful shutdown handler
def signal_handler(sig, frame):
    print(f"\n[INFO] Received signal {sig}")
    print("[INFO] Shutting down gracefully...")
    
    if redis_client:
        redis_client.close()
        print("[INFO] Redis connection closed")
    
    print("[INFO] Application stopped")
    sys.exit(0)

# Регистрация обработчиков сигналов
signal.signal(signal.SIGTERM, signal_handler)
signal.signal(signal.SIGINT, signal_handler)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=PORT)
