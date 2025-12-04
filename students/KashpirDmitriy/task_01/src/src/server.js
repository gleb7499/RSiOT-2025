const express = require('express');
const { Client } = require('pg');

const app = express();
const PORT = process.env.PORT || 8012;

const STU_ID = process.env.STU_ID || 'unknown';
const STU_GROUP = process.env.STU_GROUP || 'unknown';
const STU_VARIANT = process.env.STU_VARIANT || '34';

console.log(`Starting service | StudentID: ${STU_ID} | Group: ${STU_GROUP} | Variant: ${STU_VARIANT}`);

async function checkPostgres() {
  const client = new Client({
    host: process.env.POSTGRES_HOST || 'postgres',
    port: 5432,
    database: process.env.POSTGRES_DB,
    user: process.env.POSTGRES_USER,
    password: process.env.POSTGRES_PASSWORD,
  });
  try {
    await client.connect();
    await client.query('SELECT 1');
    await client.end();
    return true;
  } catch (err) {
    return false;
  }
}

app.get('/healthz', (req, res) => {
  res.status(200).send('OK');
});

app.get('/ready', async (req, res) => {
  const ready = await checkPostgres();
  res.status(ready ? 200 : 503).send(ready ? 'READY' : 'NOT READY');
});

app.get('/', (req, res) => {
  res.json({
    message: 'Hello from Docker Lab 01!',
    student: STU_ID,
    group: STU_GROUP,
    variant: STU_VARIANT,
    timestamp: new Date().toISOString()
  });
});

const server = app.listen(PORT, '0.0.0.0', () => {
  console.log(`Server running on http://0.0.0.0:${PORT}`);
  console.log(`Health check: http://localhost:${PORT}/healthz`);
});

process.on('SIGTERM', () => {
  console.log('SIGTERM received. Shutting down gracefully...');
  server.close(() => {
    console.log('HTTP server closed.');
    process.exit(0);
  });

  setTimeout(() => {
    console.error('Force shutdown after timeout');
    process.exit(1);
  }, 10000);
});

process.on('SIGINT', () => {
  console.log('SIGINT received.');
  process.emit('SIGTERM');
});