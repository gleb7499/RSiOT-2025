import express from 'express'
import process from 'node:process'
import { Pool } from 'pg'

const {
  PORT = 8064,
  DATABASE_URL,
  STU_ID = '220020',
  STU_GROUP = 'AS-63',
  STU_VARIANT = 'v16',
} = process.env
const pool = new Pool({ connectionString: DATABASE_URL })

const app = express()

app.get('/', (req, res) => {
  res.json({ ok: true, ts: new Date().toISOString() })
})

app.get('/health', async (_req, res) => {
  try {
    const r = await pool.query('SELECT 1 as ok;')
    res.status(200).json({ status: 'up', db: r.rows[0], time: new Date().toISOString() })
  } catch (e) {
    res.status(503).json({ status: 'down', error: String(e) })
  }
})

function logStudentMeta() {
  console.log(`[BOOT] STU_ID=${STU_ID} STU_GROUP=${STU_GROUP} STU_VARIANT=${STU_VARIANT}`)
}

const server = app.listen(PORT, () => {
  logStudentMeta()
  console.log(`[BOOT] HTTP server listening on :${PORT}`)
})

async function shutdown(signal) {
  console.log(`[SHUTDOWN] Received ${signal}. Closing HTTP server...`)
  await new Promise((resolve) => server.close(resolve))
  console.log('[SHUTDOWN] HTTP server closed. Closing PG pool...')
  await pool.end().catch((e) => console.error('[SHUTDOWN] PG pool close error:', e))
  console.log('[SHUTDOWN] Bye!')
  process.exit(0)
}

process.on('SIGTERM', () => shutdown('SIGTERM'))
process.on('SIGINT', () => shutdown('SIGINT'))
