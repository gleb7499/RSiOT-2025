import express from "express";
import process from "process";
import { createClient } from "redis";

const app = express();

// Конфигурация через ENV
const PORT = process.env.PORT || 8081;
const STU_ID = process.env.STU_ID || "1";
const STU_GROUP = process.env.STU_GROUP || "feis";
const STU_VARIANT = process.env.STU_VARIANT || "v01";
const REDIS_HOST = process.env.REDIS_HOST || "redis";
const REDIS_PORT = process.env.REDIS_PORT || 6379;

// Префикс ключей
const REDIS_PREFIX = `stu:${STU_ID}:v${STU_VARIANT}:`;

// Создаём клиента Redis
const redisClient = createClient({
  socket: {
    host: REDIS_HOST,
    port: REDIS_PORT,
  },
});

redisClient.on("error", (err) => console.error("Redis Client Error:", err));

// Обёртка для запуска сервера
async function main() {
  try {
    await redisClient.connect();
    console.log("✅ Connected to Redis");

    // Логирование всех запросов
    app.use((req, res, next) => {
      if (req.url !== '/ready' && req.url !== '/health') {
        console.log(`[${new Date().toISOString()}] ${req.method} ${req.url}`);
      }
      next();
    });

    // Health endpoint
    app.get("/health", (req, res) => res.status(200).send("OK"));

    // Readiness endpoint
    app.get("/ready", (req, res) => res.status(200).send("READY"));

    // Пример обращения к Redis
    app.get("/visit", async (req, res) => {
      try {
        const key = `${REDIS_PREFIX}visits`;
        const visits = await redisClient.incr(key);
        res.send(`Количество визитов: ${visits}`);
      } catch (err) {
        console.error("Redis operation failed:", err);
        res.status(500).send("Ошибка Redis");
      }
    });

    // Запуск сервера
    const server = app.listen(PORT, () => {
      console.log(`🚀 Server started on port ${PORT}`);
      console.log(`Student ID: ${STU_ID}, Group: ${STU_GROUP}, Variant: ${STU_VARIANT}`);
    });

    // Graceful shutdown
    const shutdown = async () => {
      console.log("⚡ SIGTERM received. Shutting down gracefully...");
      try {
        await redisClient.disconnect();
        console.log("✅ Redis client disconnected");
      } catch (err) {
        console.error("❌ Error disconnecting Redis:", err);
      }
      server.close(() => {
        console.log("✅ Server closed");
        process.exit(0);
      });

      // Форсированный выход через 5 секунд
      setTimeout(() => {
        console.error("❌ Force exit");
        process.exit(1);
      }, 5000);
    };

    process.on("SIGTERM", shutdown);
    process.on("SIGINT", shutdown);

  } catch (err) {
    console.error("❌ Failed to start server:", err);
    process.exit(1);
  }
}

// Запускаем main
main();