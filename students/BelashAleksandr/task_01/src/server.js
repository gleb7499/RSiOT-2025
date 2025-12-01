const express = require("express");
const redis = require("redis");

const app = express();

// Переменные окружения
const PORT = process.env.PORT || 8031;
const REDIS_HOST = process.env.REDIS_HOST || "localhost";
const REDIS_PORT = process.env.REDIS_PORT || 6379;
const STU_ID = process.env.STU_ID || "220031";
const STU_GROUP = process.env.STU_GROUP || "АС-64";
const STU_VARIANT = process.env.STU_VARIANT || "25";

// Redis клиент
const redisClient = redis.createClient({
  socket: {
    host: REDIS_HOST,
    port: REDIS_PORT,
  },
});

redisClient.on("error", (err) => {
  console.error("[ERROR] Redis connection error:", err);
});

redisClient.on("connect", () => {
  console.log("[INFO] Connected to Redis");
});

// Подключение к Redis
(async () => {
  try {
    await redisClient.connect();
  } catch (err) {
    console.error("[ERROR] Failed to connect to Redis:", err);
  }
})();

// Middleware для логирования
app.use((req, res, next) => {
  console.log(`[${new Date().toISOString()}] ${req.method} ${req.path}`);
  next();
});

// Health check endpoint
app.get("/ping", async (req, res) => {
  try {
    const redisKey = `stu:${STU_ID}:v${STU_VARIANT}:ping`;
    const count = (await redisClient.get(redisKey)) || "0";
    const newCount = parseInt(count) + 1;
    await redisClient.set(redisKey, newCount.toString());

    res.status(200).json({
      status: "ok",
      message: "pong",
      student_id: STU_ID,
      group: STU_GROUP,
      variant: STU_VARIANT,
      ping_count: newCount,
    });
  } catch (err) {
    console.error("[ERROR] Health check failed:", err);
    res.status(500).json({ status: "error", message: "Redis unavailable" });
  }
});

// Базовый маршрут
app.get("/", (req, res) => {
  res.json({
    service: "RSiOT Lab01",
    student: "Белаш Александр Олегович",
    group: STU_GROUP,
    variant: STU_VARIANT,
  });
});

// Запуск сервера
const server = app.listen(PORT, () => {
  console.log("=".repeat(60));
  console.log("[INFO] Server started");
  console.log(`[INFO] Student ID: ${STU_ID}`);
  console.log(`[INFO] Group: ${STU_GROUP}`);
  console.log(`[INFO] Variant: ${STU_VARIANT}`);
  console.log(`[INFO] Listening on port ${PORT}`);
  console.log("=".repeat(60));
});

// Graceful shutdown
let isShuttingDown = false;

const gracefulShutdown = async (signal) => {
  if (isShuttingDown) {
    console.log("[WARN] Shutdown already in progress");
    return;
  }

  isShuttingDown = true;
  console.log("=".repeat(60));
  console.log(`[INFO] Received ${signal}, starting graceful shutdown`);

  // Закрываем HTTP сервер
  server.close(async () => {
    console.log("[INFO] HTTP server closed");

    // Закрываем Redis соединение
    try {
      await redisClient.quit();
      console.log("[INFO] Redis connection closed");
    } catch (err) {
      console.error("[ERROR] Error closing Redis connection:", err);
    }

    console.log("[INFO] Graceful shutdown completed");
    console.log("=".repeat(60));
    process.exit(0);
  });

  // Принудительное завершение через 10 секунд
  setTimeout(() => {
    console.error("[ERROR] Forced shutdown after timeout");
    process.exit(1);
  }, 10000);
};

process.on("SIGTERM", () => gracefulShutdown("SIGTERM"));
process.on("SIGINT", () => gracefulShutdown("SIGINT"));
