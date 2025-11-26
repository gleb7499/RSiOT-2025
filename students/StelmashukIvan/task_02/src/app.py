import signal
import sys
import time
from flask import Flask, jsonify
import os

app = Flask(__name__)

HOST = os.getenv('FLASK_HOST', '0.0.0.0')
PORT = int(os.getenv('FLASK_PORT', '8054'))
DB_HOST = os.getenv('DB_HOST', 'localhost')
DB_NAME = os.getenv('DB_NAME', 'mydb')

is_shutting_down = False

def graceful_shutdown(signum, frame):
    global is_shutting_down
    print(f"Received signal {signum}, initiating graceful shutdown...")
    is_shutting_down = True
    time.sleep(2)  # Simulate cleanup
    print("Shutdown complete")
    sys.exit(0)

signal.signal(signal.SIGTERM, graceful_shutdown)
signal.signal(signal.SIGINT, graceful_shutdown)

@app.route('/ping')
def health_check():
    if is_shutting_down:
        return jsonify({"status": "shutting down"}), 503
    return jsonify({
        "status": "healthy",
        "timestamp": time.time(),
        "db_host": DB_HOST,
        "db_name": DB_NAME
    })

@app.route('/')
def hello():
    if is_shutting_down:
        return jsonify({"message": "Service is shutting down"}), 503
    return jsonify({"message": "Hello from Flask API"})

@app.route('/slow')
def slow_endpoint():
    """Эндпоинт для тестирования graceful shutdown"""
    if is_shutting_down:
        return jsonify({"error": "Service shutting down"}), 503
    time.sleep(5)
    return jsonify({"message": "Slow operation completed"})

if __name__ == '__main__':
    # Log student ENV vars at startup
    print(f"STU_ID: {os.getenv('STU_ID', 'N/A')}")
    print(f"STU_GROUP: {os.getenv('STU_GROUP', 'N/A')}")
    print(f"STU_VARIANT: {os.getenv('STU_VARIANT', 'N/A')}")
    print(f"Starting Flask application on {HOST}:{PORT}...")
    app.run(host=HOST, port=PORT, debug=False)