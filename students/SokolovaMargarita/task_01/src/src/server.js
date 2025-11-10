// src/server.js
const express = require("express");
const redis = require("redis");

const app = express();

// Используем переменную PORT, по-умолчанию 8053 (вариант 19)
const port = parseInt(process.env.PORT || "8053", 10);

// Подключение к Redis, читаем хост/порт из окружения
const redisUrl = `redis://${process.env.REDIS_HOST || 'localhost'}:${process.env.REDIS_PORT || 6379}`;
const redisClient = redis.createClient({ url: redisUrl });

redisClient.on('error', (err) => console.error("Redis error:", err));

(async () => {
  try {
    await redisClient.connect();
    console.log("Redis подключен:", redisUrl);
  } catch (err) {
    console.error("Не удалось подключиться к Redis:", err.message);
  }
})();

// /healthz — liveness
app.get("/healthz", (req, res) => res.status(200).send("OK"));

// /ready — readiness: проверяем Redis
app.get("/ready", async (req, res) => {
  try {
    // ping вернёт "PONG" или сработает исключение
    await redisClient.ping();
    res.status(200).send("READY");
  } catch (err) {
    res.status(500).send("NOT READY");
  }
});

// Корневой и дополнительные эндпоинты
app.get("/", (req, res) => {
  res.send("<h1>Добро пожаловать!</h1><p>Сервер работает через Docker Compose ✅</p>");
});
app.get("/about", (req, res) => res.send("<h1>О проекте</h1><p>Это вторая страница.</p>"));
app.get("/contact", (req, res) => res.send("<h1>Контакты</h1><p>Это третья страница.</p>"));

// Запуск сервера
const server = app.listen(port, '0.0.0.0', () => {
  console.log(`Server listening on http://0.0.0.0:${port}`);
});

// Graceful shutdown: ловим SIGTERM и закрываем сервер и Redis
const shutdown = async () => {
  console.log("SIGTERM/SIGINT получен — начинаем корректный shutdown...");
  server.close(async () => {
    try {
      await redisClient.quit();
      console.log("Redis корректно завершён");
    } catch (err) {
      console.error("Ошибка при завершении Redis:", err);
    }
    console.log("HTTP сервер остановлен");
    process.exit(0);
  });

  // На случай, если close "зависнет" — аварийный выход через 10s
  setTimeout(() => {
    console.warn("Shutdown занял слишком много времени — форсируем exit");
    process.exit(1);
  }, 10000).unref();
};

process.on('SIGTERM', shutdown);
process.on('SIGINT', shutdown); // чтобы CTRL+C локально тоже корректно работало
