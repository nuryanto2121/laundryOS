/**
 * WhatsApp Client
 *
 * Initializes the Baileys socket and manages the connection lifecycle.
 * Session management is delegated to session.js.
 *
 * Milestone 2 scope:
 * - Persist session using useMultiFileAuthState via session.load()
 * - Reuse existing session on startup (no QR if session is valid)
 * - Show QR Code only when no valid session exists
 * - Reconnect automatically on restartRequired
 * - Stop gracefully on loggedOut (no auto-delete, no restart)
 *
 * Out of scope:
 * - Message handling (Milestone 3)
 * - Go Backend integration (Milestone 5)
 */

import fs from 'node:fs';
import path from 'node:path';
import makeWASocket, {
  DisconnectReason,
  fetchLatestBaileysVersion,
  jidEncode,
  jidNormalizedUser,
} from '@whiskeysockets/baileys';
import qrcode from 'qrcode-terminal';
import logger from '../logger/index.js';
import * as session from './session.js';
import { handleMessagesUpsert } from './handler.js';

const SESSION_PATH = process.env.SESSION_PATH || session.DEFAULT_SESSION_PATH;

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
let activeSocket = null;

/**
 * Returns the current WhatsApp connection state.
 * @returns {string}
 */
export function getState() {
  return currentState;
}

/**
 * Starts the WhatsApp client.
 *
 * Loads auth state from disk via session.load().
 * If no session exists, Baileys will emit a QR code.
 * If a valid session exists, Baileys will connect silently.
 */
export async function startWhatsApp() {
  logger.info('Starting WhatsApp client...');
  logger.info(`Session path resolved to: ${SESSION_PATH}`);

  const credsExist = fs.existsSync(path.join(SESSION_PATH, 'creds.json'));
  if (credsExist) {
    logger.info('Existing session detected. Resuming connection...');
  } else {
    logger.info('No existing session detected. A new QR code will be generated.');
  }

  const { state, saveCreds } = await session.load(SESSION_PATH);

  let version = [2, 3000, 1015901307];
  try {
    const fetched = await fetchLatestBaileysVersion();
    version = fetched.version;
  } catch {
    logger.warn(`Failed to fetch WA version (network error). Using fallback v${version.join('.')}`);
  }

  logger.info(`Using WA v${version.join('.')}`);

  const socket = makeWASocket({
    auth: state,
    version,
    browser: ['LaundryOS Gateway', 'Chrome', '1.0.0'],
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

  activeSocket = socket;

  socket.ev.on('creds.update', saveCreds);

  socket.ev.on('messages.upsert', (upsert) => {
    handleMessagesUpsert(upsert);
  });

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

      const id = socket.user?.id || 'Unknown';
      const name = socket.user?.name || socket.user?.verifiedName || 'Unknown';

      logger.info(`Connected\n\nID      : ${id}\nName    : ${name}\nPlatform: WhatsApp`);
    }

    if (connection === 'close') {
      currentState = WAState.DISCONNECTED;
      logger.info(JSON.stringify(lastDisconnect, null, 2));
      // console.dir(lastDisconnect, { depth: null });

      const statusCode = lastDisconnect?.error?.output?.statusCode;
      const reason = DisconnectReason[statusCode] ?? `code ${statusCode}`;

      logger.warn(`WhatsApp disconnected. Reason: ${reason}`);

      if (statusCode === DisconnectReason.loggedOut) {
        // Session is no longer valid. Stop the WhatsApp lifecycle gracefully.
        // The process remains running so health checks continue to respond.
        // To re-login: delete the session directory and restart the gateway.
        logger.error(
          `WhatsApp logged out. Gateway is stopped. ` +
            `Delete ${SESSION_PATH} and restart to scan a new QR.`
        );
        return;
      }

      // For all other disconnect reasons (including restartRequired), reconnect.
      logger.info('Reconnecting in 2 seconds...');
      setTimeout(startWhatsApp, 2000);
    }
  });

  return socket;
}

/**
 * Sends a text message through the active WhatsApp connection.
 * Supports both standard phone numbers and group JIDs (e.g. 120363xxx@g.us).
 * @param {string} destination The destination phone number or group JID.
 * @param {string} message The text message to send.
 * @returns {Promise<boolean>} True if sent successfully.
 */
export async function sendMessage(destination, message) {
  if (currentState !== WAState.CONNECTED || !activeSocket) {
    throw new Error('WhatsApp is not connected');
  }

  let jid;
  if (destination.includes('@')) {
    jid = destination;
  } else {
    let clean = destination.replace(/\D/g, '');
    if (clean.startsWith('0')) {
      clean = '62' + clean.slice(1);
    }
    jid = jidNormalizedUser(jidEncode(clean, 's.whatsapp.net'));
  }

  try {
    await activeSocket.sendMessage(jid, { text: message });
    logger.info(
      `\nOutgoing Message\n\n` +
        `To        : ${destination}\n` +
        `Type      : text\n` +
        `Text      : ${message}\n` +
        `Status    : Sent\n` +
        `────────────────────────────`
    );
    return true;
  } catch (error) {
    logger.error(`Failed to send message. ${error.message || error}`);
    throw error;
  }
}

/**
 * Retrieves all WhatsApp groups the connected account is currently participating in.
 * @returns {Promise<Array<{id: string, name: string, participants_count: number}>>}
 */
export async function getParticipatingGroups() {
  if (currentState !== WAState.CONNECTED || !activeSocket) {
    return [];
  }
  try {
    const groups = await activeSocket.groupFetchAllParticipating();
    return Object.values(groups).map((g) => ({
      id: g.id,
      name: g.subject,
      participants_count: g.participants ? g.participants.length : 0,
    }));
  } catch (error) {
    logger.error(`Failed to fetch participating groups: ${error.message || error}`);
    return [];
  }
}
