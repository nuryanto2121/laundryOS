/**
 * WhatsApp Client
 *
 * Initializes the Baileys socket and manages the connection lifecycle.
 * This module is the primary responsibility of the Gateway.
 *
 * Milestone 1 scope:
 * - Connect to WhatsApp
 * - Display QR code in terminal
 * - Log connection state transitions
 *
 * Out of scope (future milestones):
 * - Session persistence
 * - Message handling
 * - Go Backend integration
 */

import makeWASocket, { DisconnectReason, useMultiFileAuthState } from '@whiskeysockets/baileys';
import qrcode from 'qrcode-terminal';
import logger from '../logger/index.js';

/**
 * WhatsApp connection state.
 * @readonly
 * @enum {string}
 */
export const WAState = {
  DISCONNECTED: 'DISCONNECTED',
  CONNECTING: 'CONNECTING',
  CONNECTED: 'CONNECTED',
};

/**
 * Holds the current connection state.
 * Updated on every connection.update event.
 * @type {string}
 */
let currentState = WAState.DISCONNECTED;

/**
 * Returns the current WhatsApp connection state.
 * Used by the health endpoint to expose state externally.
 * @returns {string}
 */
export function getState() {
  return currentState;
}

/**
 * Starts the WhatsApp client.
 *
 * NOTE: Auth state is in-memory for Milestone 1.
 * Session persistence will be added in Milestone 2.
 */
export async function startWhatsApp() {
  logger.info('Starting WhatsApp client...');

  // In-memory auth state — session is lost on restart (Milestone 1 only).
  // Milestone 2 will replace this with useMultiFileAuthState persisted to disk.
  const { state, saveCreds } = await useMultiFileAuthState('session_temp');

  const socket = makeWASocket({
    auth: state,
    // Disable Baileys' built-in QR printing; we handle it manually below
    // for clean, explicit terminal output.
    printQRInTerminal: false,
    // Suppress Baileys' internal logger output to keep terminal readable.
    logger: {
      level: 'silent',
      log: () => {},
      info: () => {},
      warn: () => {},
      error: () => {},
      debug: () => {},
      trace: () => {},
      child: () => ({
        level: 'silent',
        log: () => {},
        info: () => {},
        warn: () => {},
        error: () => {},
        debug: () => {},
        trace: () => {},
        child: () => ({}),
      }),
    },
  });

  socket.ev.on('creds.update', saveCreds);

  socket.ev.on('connection.update', (update) => {
    const { connection, lastDisconnect, qr } = update;

    if (qr) {
      logger.info('QR Code generated — scan with WhatsApp to connect.');
      qrcode.generate(qr, { small: true });
    }

    if (connection === 'connecting') {
      currentState = WAState.CONNECTING;
      logger.info('WhatsApp connecting...');
    }

    if (connection === 'open') {
      currentState = WAState.CONNECTED;
      logger.info('WhatsApp connected.');
    }

    if (connection === 'close') {
      currentState = WAState.DISCONNECTED;

      const statusCode = lastDisconnect?.error?.output?.statusCode;
      const reason = DisconnectReason[statusCode] ?? `code ${statusCode}`;

      logger.warn(`WhatsApp disconnected. Reason: ${reason}`);

      // Reconnect automatically unless logged out
      if (statusCode !== DisconnectReason.loggedOut) {
        logger.info('Reconnecting in 2 seconds...');
        setTimeout(startWhatsApp, 2000);
      } else {
        logger.error('WhatsApp logged out. Please remove session_temp folder and restart to scan new QR.');
      }
    }
  });

  return socket;
}
