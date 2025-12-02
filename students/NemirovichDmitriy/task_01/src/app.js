const express = require("express");
const redis = require("redis");

const app = express();
const PORT = process.env.PORT || 8001;

// Метаданные студента из переменных окружения
const STU_ID = process.env.STU_ID || "220050";
const STU_GROUP = process.env.STU_GROUP || "АС-64";
const STU_VARIANT = process.env.STU_VARIANT || "37";

// Redis конфигурация
const REDIS_HOST = process.env.REDIS_HOST || "localhost";
const REDIS_PORT = process.env.REDIS_PORT || 6379;
const REDIS_PREFIX = `stu:${STU_ID}:v${STU_VARIANT}:`;

let redisClient;
let isShuttingDown = false;

// Подключение к Redis
async function connectRedis() {
  redisClient = redis.createClient({
    socket: {
      host: REDIS_HOST,
      port: REDIS_PORT,
    },
  });

  redisClient.on("error", (err) => {
    console.error("[ERROR] Redis Client Error:", err);
  });

  redisClient.on("connect", () => {
    console.log("[INFO] Redis connected");
  });

  await redisClient.connect();
}

// Логирование метаданных при старте
console.log("[INFO] Starting application...");
console.log(`[INFO] StudentID: ${STU_ID}`);
console.log(`[INFO] Group: ${STU_GROUP}`);
console.log(`[INFO] Variant: ${STU_VARIANT}`);
console.log(`[INFO] Port: ${PORT}`);
console.log(`[INFO] Redis Host: ${REDIS_HOST}:${REDIS_PORT}`);
console.log(`[INFO] Redis Prefix: ${REDIS_PREFIX}`);

// Базовый endpoint
app.get("/", async (req, res) => {
  try {
    const key = `${REDIS_PREFIX}counter`;
    const counter = await redisClient.incr(key);

    console.log(`[INFO] Request processed, counter: ${counter}`);

    res.json({
      message: "Hello from RSiOT Lab 01",
      student: {
        id: STU_ID,
        group: STU_GROUP,
        variant: STU_VARIANT,
      },
      counter: counter,
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    console.error("[ERROR] Request processing error:", error);
    res.status(500).json({ error: "Internal server error" });
  }
});

// Health check endpoint
app.get("/ready", async (req, res) => {
  if (isShuttingDown) {
    return res.status(503).json({ status: "shutting down" });
  }

  try {
    await redisClient.ping();
    res.status(200).json({ status: "ready" });
  } catch (error) {
    console.error("[ERROR] Health check failed:", error);
    res.status(503).json({ status: "not ready", error: error.message });
  }
});

// Graceful shutdown
async function gracefulShutdown(signal) {
  console.log(`[INFO] ${signal} received, starting graceful shutdown...`);
  isShuttingDown = true;

  // Закрываем HTTP сервер
  server.close(() => {
    console.log("[INFO] HTTP server closed");
  });

  // Закрываем Redis подключение
  try {
    await redisClient.quit();
    console.log("[INFO] Redis connection closed");
  } catch (error) {
    console.error("[ERROR] Error closing Redis:", error);
  }

  console.log("[INFO] Graceful shutdown completed");
  process.exit(0);
}

// Запуск приложения
let server;

(async () => {
  try {
    await connectRedis();

    server = app.listen(PORT, () => {
      console.log(`[INFO] Server started on port ${PORT}`);
      console.log("[INFO] Application ready to handle requests");
    });

    // Обработчики сигналов для graceful shutdown
    process.on("SIGTERM", () => gracefulShutdown("SIGTERM"));
    process.on("SIGINT", () => gracefulShutdown("SIGINT"));
  } catch (error) {
    console.error("[ERROR] Failed to start application:", error);
    process.exit(1);
  }
})();
