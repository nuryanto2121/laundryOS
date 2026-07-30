/**
 * LaundryOS Gateway — Entry Point
 *
 * Startup order:
 * 1. Start WhatsApp client (primary responsibility)
 * 2. Start HTTP server (secondary — liveness health endpoint)
 */

import express from 'express';
import logger from './logger/index.js';
import { startWhatsApp, sendMessage, getState } from './whatsapp/client.js';

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
app.get('/health', (req, res) => {
  logger.info(`GET /health`);
  res.status(200).json({
    status: 'UP',
    whatsapp: getState(),
  });
  logger.info(`GET /health -> 200`);
});

/**
 * POST /send-message
 * Sends a WhatsApp text message to the specified phone number.
 */
app.post('/send-message', async (req, res) => {
  logger.info(`POST /send-message`);
  try {
    const { phone, message } = req.body;

    if (typeof phone !== 'string' || !phone.trim()) {
      logger.info(`POST /send-message -> 400`);
      return res.status(400).json({ success: false, error: 'phone is required' });
    }
    if (typeof message !== 'string' || !message.trim()) {
      logger.info(`POST /send-message -> 400`);
      return res.status(400).json({ success: false, error: 'message is required' });
    }

    await sendMessage(phone.trim(), message.trim());

    logger.info(`POST /send-message -> 200`);
    return res.status(200).json({ success: true });
  } catch (error) {
    if (error.message === 'WhatsApp is not connected') {
      logger.info(`POST /send-message -> 503`);
      return res.status(503).json({ success: false, error: error.message });
    }
    logger.info(`POST /send-message -> 500`);
    return res.status(500).json({ success: false, error: 'Internal server error' });
  }
});

app.listen(PORT, () => {
  logger.info(`laundryos-gateway v${APP_VERSION} started on port ${PORT}`);
});
