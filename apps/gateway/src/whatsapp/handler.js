import logger from '../logger/index.js';
import { normalizeMessage } from './mapper.js';

export function handleMessagesUpsert(upsert) {
  if (upsert.type !== 'notify') {
    return;
  }

  const messages = upsert.messages || [];

  for (const msg of messages) {
    if (!msg || !msg.key) {
      continue;
    }

    if (msg.key.fromMe) {
      continue;
    }

    const remoteJid = msg.key.remoteJid || '';
    if (remoteJid.endsWith('@g.us') || remoteJid.endsWith('@broadcast')) {
      continue;
    }

    const normalized = normalizeMessage(msg);

    let displayTime = '-';
    let ts = normalized.timestamp;

    if (typeof ts === 'object' && ts !== null && ts.low !== undefined) {
      ts = ts.low;
    }

    ts = Number(ts) * 1000;
    if (!isNaN(ts) && ts > 0) {
      displayTime = new Date(ts).toISOString();
    }

    logger.info(
      `\nIncoming Message\n\n` +
        `From      : ${normalized.from}\n` +
        `Push Name : ${normalized.pushName}\n` +
        `Type      : ${normalized.type}\n` +
        `Text      : ${normalized.text}\n` +
        `Time      : ${displayTime}\n` +
        `────────────────────────────`
    );
  }
}
