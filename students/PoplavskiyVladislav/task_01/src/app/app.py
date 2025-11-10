import os, signal, logging, time
from flask import Flask, jsonify
import redis

logging.basicConfig(level=logging.INFO, format='[%(asctime)s] %(levelname)s %(message)s')
logger = logging.getLogger(__name__)
STU_ID = os.getenv('STU_ID', '220021')
STU_GROUP = os.getenv('STU_GROUP', 'ASOI 63')
STU_VARIANT = os.getenv('STU_VARIANT', '17')
REDIS_HOST = os.getenv('REDIS_HOST', 'redis')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
REDIS_PREFIX = f"stu:{STU_ID}:v{STU_VARIANT}:entity:"
app = Flask(__name__)
redis_client = None

def get_redis():
    global redis_client
    if redis_client is None:
        redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)
    return redis_client

def log_startup():
    logger.info(f"Starting app — STU_ID={STU_ID}, STU_GROUP={STU_GROUP}, STU_VARIANT={STU_VARIANT}")
    logger.info(f"Redis -> {REDIS_HOST}:{REDIS_PORT}, key prefix: {REDIS_PREFIX}")

log_startup()  # вызываем сразу при запуске

@app.route('/ready')
def ready():
    return jsonify({"status": "ok"}), 200

@app.route('/')
def index():
    logger.info('Received request to /')
    r = get_redis()
    key = REDIS_PREFIX + 'counter'
    try:
        count = r.incr(key)
    except Exception as e:
        logger.exception('Redis error')
        count = None
    return jsonify({"message": "Hello from student service", "hits": count})

def shutdown_handler(signum, frame):
    logger.info(f'Received signal {signum}, starting graceful shutdown')
    try:
        if redis_client:
            redis_client.close()
            logger.info('Redis connection closed')
    except Exception:
        logger.exception('Error while closing redis')
    logger.info('Exiting process')
    time.sleep(0.2)
    os._exit(0)

signal.signal(signal.SIGTERM, shutdown_handler)
signal.signal(signal.SIGINT, shutdown_handler)

if __name__ == '__main__':
    port = int(os.getenv('PORT', 8051))
    app.run(host='0.0.0.0', port=port)
