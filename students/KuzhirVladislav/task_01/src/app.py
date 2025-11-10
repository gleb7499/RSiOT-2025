import os
import signal
import sys
from flask import Flask, jsonify
from threading import Event

app = Flask(__name__)
shutdown_event = Event()

STU_ID = os.getenv('STU_ID', 'unknown')
STU_GROUP = os.getenv('STU_GROUP', 'unknown')
STU_VARIANT = os.getenv('STU_VARIANT', 'unknown')

@app.route('/')
def index():
    return jsonify({
        'message': 'Hello from Flask!',
        'student_id': STU_ID,
        'group': STU_GROUP,
        'variant': STU_VARIANT
    })

@app.route('/health')
def healthz():
    return jsonify({'status': 'ok'}), 200

def handle_sigterm(signum, frame):
    print('Received SIGTERM, shutting down gracefully...')
    shutdown_event.set()

signal.signal(signal.SIGTERM, handle_sigterm)

if __name__ == '__main__':
    print(f"Starting Flask app | STU_ID={STU_ID} | GROUP={STU_GROUP} | VARIANT={STU_VARIANT}")
    from threading import Thread
    def run_app():
        app.run(host='0.0.0.0', port=8032)
    t = Thread(target=run_app)
    t.start()
    try:
        while not shutdown_event.is_set():
            shutdown_event.wait(1)
    except KeyboardInterrupt:
        print('Received KeyboardInterrupt, shutting down...')
    sys.exit(0)
