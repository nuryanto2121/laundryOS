import logger from '../logger/index.js';
import { normalizeMessage } from './mapper.js';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

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
    if (remoteJid.endsWith('@broadcast')) {
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
        `From        : ${normalized.from}\n` +
        `Is Group    : ${normalized.isGroup}\n` +
        `Participant : ${normalized.participant || '-'}\n` +
        `Push Name   : ${normalized.pushName}\n` +
        `Type        : ${normalized.type}\n` +
        `Text        : ${normalized.text}\n` +
        `Quoted Text : ${normalized.quotedText || '-'}\n` +
        `Time        : ${displayTime}\n` +
        `────────────────────────────`
    );

    // Forward text or reply messages to the backend bot endpoint (transport only)
    if ((normalized.text || normalized.quotedText) && normalized.from) {
      forwardToBackend(normalized);
    }
  }
}

async function forwardToBackend(normalized) {
  const url = `${API_BASE_URL}/bot/messages`;
  try {
    let ts = normalized.timestamp;
    if (typeof ts === 'object' && ts !== null && ts.low !== undefined) {
      ts = ts.low;
    }
    ts = Number(ts);
    if (isNaN(ts) || ts <= 0) {
      ts = Math.floor(Date.now() / 1000);
    }

    const payload = {
      from: normalized.from,
      text: normalized.text || '',
      push_name: normalized.pushName || '',
      timestamp: ts,
      is_group: !!normalized.isGroup,
      group_jid: normalized.groupJid || '',
      participant: normalized.participant || '',
      quoted_text: normalized.quotedText || '',
    };

    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      logger.warn(`Failed forwarding message to backend: ${res.status} ${res.statusText}`);
    } else {
      logger.info(`Message forwarded to backend successfully: ${normalized.from}`);
    }
  } catch (err) {
    logger.error('Error forwarding message to backend:', { error: err.message });
  }
}
