/**
 * LaundryOS Gateway — Entry Point
 *
 * Startup order:
 * 1. Start WhatsApp client (primary responsibility)
 * 2. Start HTTP server (secondary — liveness health endpoint)
 */

import express from 'express';
import logger from './logger/index.js';
import { startWhatsApp } from './whatsapp/client.js';

const PORT = process.env.PORT || 3000;
const APP_VERSION = process.env.APP_VERSION || '0.1.0';

// ---------------------------------------------------------------------------
// WhatsApp — primary
// ---------------------------------------------------------------------------

startWhatsApp().catch((err) => {
  logger.error('WhatsApp client failed to start.', { error: err.message });
});

// ---------------------------------------------------------------------------
// HTTP server — secondary
// ---------------------------------------------------------------------------

const app = express();

app.use(express.json());

/**
 * GET /health
 * Liveness probe. Responds 200 while the process is running.
 * WhatsApp connectivity is an application-level concern, not a liveness indicator.
 */
app.get('/health', (_req, res) => {
  res.status(200).json({
    status: 'ok',
    service: 'laundryos-gateway',
    version: APP_VERSION,
    timestamp: new Date().toISOString(),
  });
});

app.listen(PORT, () => {
  logger.info(`laundryos-gateway v${APP_VERSION} started on port ${PORT}`);
});
