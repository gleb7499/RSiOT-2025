const express = require('express');
const { Pool } = require('pg');

const app = express();
const port = process.env.PORT || 8072;

// Environment variables for DB config
const dbConfig = {
  user: process.env.DB_USER || 'postgres',
  host: process.env.DB_HOST || 'db-AC-63-220046-v10',
  database: process.env.DB_NAME || 'app_220046_v10',
  password: process.env.DB_PASSWORD || 'password',
  port: process.env.DB_PORT || 5432,
};

// Create DB pool
const pool = new Pool(dbConfig);

// Log student env vars on start
console.log(`STU_ID: ${process.env.STU_ID}, STU_GROUP: ${process.env.STU_GROUP}, STU_VARIANT: ${process.env.STU_VARIANT}`);

// Health endpoint for liveness/readiness
app.get('/ping', async (req, res) => {
  try {
    // Test DB connection
    await pool.query('SELECT 1');
    res.status(200).send('OK');
  } catch (err) {
    console.error('DB health check failed:', err);
    res.status(500).send('DB Error');
  }
});

// Simple root endpoint
app.get('/', (req, res) => {
  res.send('Hello from Node/Express service!');
});

// Graceful shutdown
const server = app.listen(port, async () => {
  console.log(`Server running on port ${port}`);
  try {
    // Connect to DB on startup
    await pool.connect();
    console.log('Connected to Postgres');
  } catch (err) {
    console.error('Failed to connect to Postgres:', err);
    process.exit(1);
  }
});

process.on('SIGTERM', () => {
  console.log('SIGTERM received. Shutting down gracefully...');
  server.close(() => {
    console.log('HTTP server closed.');
    pool.end(() => {
      console.log('DB pool closed.');
      process.exit(0);
    });
  });
});