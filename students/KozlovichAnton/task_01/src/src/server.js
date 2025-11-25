// Импорт необходимых модулей для работы сервера
const express = require('express'); // Фреймворк для создания веб-сервера
const redis = require('redis'); // Клиент для взаимодействия с Redis
const process = require('process'); // Модуль для работы с процессом Node.js

// Создание экземпляра Express-приложения
const myApp = express();
const serverPort = process.env.PORT || 8093; // Порт сервера, по умолчанию 8093

// Получение данных студента из переменных окружения для логирования
const studentIdentifier = process.env.STU_ID;
const studentGroupName = process.env.STU_GROUP;
const studentTaskVariant = process.env.STU_VARIANT;
console.log(`Запуск сервера с параметрами: STU_ID=${studentIdentifier}, STU_GROUP=${studentGroupName}, STU_VARIANT=${studentTaskVariant}`);

// Настройка подключения к Redis из переменных окружения
const redisServerHost = process.env.REDIS_HOST;
const redisServerPort = process.env.REDIS_PORT;
const redisServerPassword = process.env.REDIS_PASSWORD;

// Создание клиента Redis с заданными параметрами
const { createClient } = require('redis');
const redisConnection = redis.createClient({
  url: `redis://${process.env.REDIS_HOST || '127.0.0.1'}:${process.env.REDIS_PORT || 6379}`,
  password: process.env.REDIS_PASSWORD || undefined,
});
redisConnection.on('error', (err) => {
  console.error('Ошибка подключения к Redis:', err);
  console.log('Сервер продолжает работу без подключения к Redis');
});
redisConnection.connect().then(() => {
  console.log('Успешное подключение к Redis');
}).catch((err) => {
  console.error('Не удалось подключиться к Redis:', err);
  console.log('Сервер продолжает работу без подключения к Redis');
});

// Определение маршрутов API
myApp.get('/ready', (request, response) => {
  response.send('OK'); // Ответ на проверку готовности
  console.log('Маршрут /ready вызван. Сервис функционирует нормально');
});

// Новый маршрут для проверки здоровья сервиса
myApp.get('/health', async (request, response) => {
  try {
    // Проверка подключения к Redis
    await redisConnection.ping();
    response.json({ status: 'healthy', redis: 'connected' });
    console.log('Маршрут /health: Сервис здоров, Redis подключен');
  } catch (error) {
    response.status(500).json({ status: 'unhealthy', redis: 'disconnected' });
    console.error('Маршрут /health: Ошибка проверки здоровья:', error.message);
  }
});

// Реализация корректного завершения работы сервера
process.on('SIGTERM', async () => {
  console.log('Начинается завершение работы сервера...');
  try {
    await redisConnection.quit(); // Закрытие соединения с Redis
    console.log('Соединение с Redis закрыто');
    process.exit(0); // Успешное завершение
  } catch (shutdownError) {
    console.error('Ошибка при закрытии соединения с Redis:', shutdownError);
    process.exit(1); // Завершение с ошибкой
  }
});

process.on('SIGINT', async () => {
  console.log('Начинается завершение работы сервера...');
  try {
    await redisConnection.quit(); // Закрытие соединения с Redis
    console.log('Соединение с Redis закрыто');
    process.exit(0); // Успешное завершение
  } catch (shutdownError) {
    console.error('Ошибка при закрытии соединения с Redis:', shutdownError);
    process.exit(1); // Завершение с ошибкой
  }
});

// Запуск сервера на указанном порту
myApp.listen(serverPort, () => {
  console.log(`Сервер запущен и слушает порт ${serverPort}`);
});
