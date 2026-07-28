/**
 * LaundryOS Gateway — Entry Point
 *
 * Startup order:
 * 1. Start WhatsApp client (primary responsibility)
 * 2. Start HTTP server (secondary — health endpoint only for Milestone 1)
 */

import express from 'express';
import logger from './logger/index.js';
import { startWhatsApp, getState } from './whatsapp/client.js';

const PORT = process.env.PORT || 3000;
const APP_VERSION = process.env.APP_VERSION || '0.1.0';

// ---------------------------------------------------------------------------
// WhatsApp — primary
// ---------------------------------------------------------------------------

startWhatsApp().catch((err) => {
  logger.error('WhatsApp client failed to start.', { error: err.message });
  process.exit(1);
});

// ---------------------------------------------------------------------------
// HTTP server — secondary
// ---------------------------------------------------------------------------

const app = express();

app.use(express.json());

/**
 * GET /health
 * Standard health check. Includes current WhatsApp connection state.
 */
app.get('/health', (_req, res) => {
  res.status(200).json({
    status: 'ok',
    service: 'laundryos-gateway',
    version: APP_VERSION,
    timestamp: new Date().toISOString(),
    whatsapp: getState(),
  });
});

app.listen(PORT, () => {
  logger.info(`laundryos-gateway v${APP_VERSION} started on port ${PORT}`);
});
