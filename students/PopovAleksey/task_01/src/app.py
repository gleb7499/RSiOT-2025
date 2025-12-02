import os
import signal
import sys
import logging
from flask import Flask, jsonify
import psycopg2
from psycopg2.extras import RealDictCursor

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Переменные окружения
STU_ID = os.getenv('STU_ID', '220051')
STU_GROUP = os.getenv('STU_GROUP', 'АС-64')
STU_VARIANT = os.getenv('STU_VARIANT', '38')
DB_HOST = os.getenv('DB_HOST', 'db')
DB_PORT = os.getenv('DB_PORT', '5432')
DB_NAME = os.getenv('DB_NAME', f'app_{STU_ID}_v{STU_VARIANT}')
DB_USER = os.getenv('DB_USER', 'postgres')
DB_PASSWORD = os.getenv('DB_PASSWORD', 'postgres')

# Флаг для graceful shutdown
shutdown_flag = False


def get_db_connection():
    """Создание подключения к БД"""
    try:
        conn = psycopg2.connect(
            host=DB_HOST,
            port=DB_PORT,
            database=DB_NAME,
            user=DB_USER,
            password=DB_PASSWORD
        )
        return conn
    except Exception as e:
        logger.error(f"Ошибка подключения к БД: {e}")
        return None


def init_db():
    """Инициализация БД"""
    conn = get_db_connection()
    if conn:
        try:
            cursor = conn.cursor()
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS requests (
                    id SERIAL PRIMARY KEY,
                    endpoint VARCHAR(255),
                    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
            """)
            conn.commit()
            cursor.close()
            conn.close()
            logger.info("БД инициализирована")
        except Exception as e:
            logger.error(f"Ошибка инициализации БД: {e}")


def signal_handler(signum, frame):
    """Обработчик сигнала SIGTERM для graceful shutdown"""
    global shutdown_flag
    logger.info(f"Получен сигнал {signum}. Начинаем graceful shutdown...")
    shutdown_flag = True
    logger.info("Приложение корректно завершено")
    sys.exit(0)


# Регистрация обработчика сигнала
signal.signal(signal.SIGTERM, signal_handler)


@app.route('/live')
def health():
    """Health check endpoint"""
    conn = get_db_connection()
    if conn:
        conn.close()
        return jsonify({"status": "healthy", "database": "connected"}), 200
    return jsonify({"status": "unhealthy", "database": "disconnected"}), 503


@app.route('/')
def index():
    """Главная страница"""
    conn = get_db_connection()
    if conn:
        try:
            cursor = conn.cursor()
            cursor.execute("INSERT INTO requests (endpoint) VALUES (%s)", ('/',))
            conn.commit()
            cursor.close()
            conn.close()
        except Exception as e:
            logger.error(f"Ошибка записи в БД: {e}")
    
    return jsonify({
        "message": "Flask App для ЛР01",
        "student_id": STU_ID,
        "group": STU_GROUP,
        "variant": STU_VARIANT
    }), 200


@app.route('/requests')
def get_requests():
    """Получение всех запросов из БД"""
    conn = get_db_connection()
    if conn:
        try:
            cursor = conn.cursor(cursor_factory=RealDictCursor)
            cursor.execute("SELECT * FROM requests ORDER BY timestamp DESC LIMIT 10")
            requests = cursor.fetchall()
            cursor.close()
            conn.close()
            return jsonify(requests), 200
        except Exception as e:
            logger.error(f"Ошибка чтения из БД: {e}")
            return jsonify({"error": str(e)}), 500
    return jsonify({"error": "Database connection failed"}), 503


if __name__ == '__main__':
    logger.info(f"Запуск приложения. STU_ID={STU_ID}, STU_GROUP={STU_GROUP}, STU_VARIANT={STU_VARIANT}")
    init_db()
    app.run(host='0.0.0.0', port=8002)
