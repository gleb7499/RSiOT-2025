import os
import signal
import sys
import redis
from flask import Flask, jsonify

print("=== Flask app init start ===", flush=True)

REDIS_HOST = os.environ.get('REDIS_HOST', 'localhost')
REDIS_PORT = int(os.environ.get('REDIS_PORT', 6379))
STU_ID = os.environ.get('STU_ID', 'UNKNOWN')
STU_GROUP = os.environ.get('STU_GROUP', 'UNKNOWN')
STU_VARIANT = os.environ.get('STU_VARIANT', 'UNKNOWN')

print(f"Start app with STU_ID={STU_ID}, STU_GROUP={STU_GROUP}, STU_VARIANT={STU_VARIANT}", flush=True)
print(f"REDIS_HOST={REDIS_HOST}, REDIS_PORT={REDIS_PORT}", flush=True)

app = Flask(__name__)

print("=== Before redis.Redis() ===", flush=True)
r = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)
print("=== After redis.Redis() ===", flush=True)

@app.route('/')
def index():
    key = f"stu:{STU_ID}:v{STU_VARIANT}:counter"
    try:
        value = r.incr(key)
    except Exception as e:
        print(f"Redis error: {e}", flush=True)
        value = -1
    return jsonify({"hello": "world", "counter": value})

@app.route('/live')
def live():
    return jsonify({"status": "ok"}), 200

@app.route('/shutdown_info')
def shutdown_info():
    return jsonify({"message": "App will shutdown gracefully"}), 200

def handle_sigterm(*args):
    print("Received SIGTERM, shutting down gracefully...", flush=True)
    sys.exit(0)

if __name__ == '__main__':
    print("=== Flask native run start ===", flush=True)
    signal.signal(signal.SIGTERM, handle_sigterm)
    app.run(host='0.0.0.0', port=8043)