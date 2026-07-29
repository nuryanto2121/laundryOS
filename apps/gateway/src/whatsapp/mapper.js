export function normalizeMessage(msg) {
  const id = msg?.key?.id || '';
  const remoteJid = msg?.key?.remoteJid || '';
  const from = remoteJid.replace('@s.whatsapp.net', '');
  const pushName = msg?.pushName || '';
  const timestamp = msg?.messageTimestamp;

  let type = 'unknown';
  let text = '';

  const message = msg?.message;
  if (message) {
    if (message.conversation) {
      type = 'text';
      text = message.conversation;
    } else if (message.extendedTextMessage) {
      type = 'text';
      text = message.extendedTextMessage.text || '';
    } else if (message.imageMessage) {
      type = 'image';
      text = message.imageMessage.caption || '';
    } else if (message.videoMessage) {
      type = 'video';
      text = message.videoMessage.caption || '';
    } else if (message.documentMessage) {
      type = 'document';
      text = message.documentMessage.caption || '';
    }
  }

  return {
    id,
    from,
    pushName,
    timestamp,
    type,
    text,
  };
}
