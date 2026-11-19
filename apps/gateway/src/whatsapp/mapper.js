export function normalizeMessage(msg) {
  const id = msg?.key?.id || '';
  const remoteJid = msg?.key?.remoteJid || '';
  const altJid = msg?.key?.remoteJidAlt || msg?.key?.participantAlt || '';
  const isGroup = remoteJid.endsWith('@g.us');

  // Prefer phone number (@s.whatsapp.net) if remoteJid is a privacy/companion @lid
  let from = remoteJid.replace('@s.whatsapp.net', '');
  if (remoteJid.endsWith('@lid') && altJid.endsWith('@s.whatsapp.net')) {
    from = altJid.replace('@s.whatsapp.net', '');
  }

  let participant = '';
  if (isGroup) {
    const rawPart = msg?.key?.participant || '';
    const rawAlt = msg?.key?.participantAlt || '';
    participant = rawPart.replace('@s.whatsapp.net', '');
    if (rawPart.endsWith('@lid') && rawAlt.endsWith('@s.whatsapp.net')) {
      participant = rawAlt.replace('@s.whatsapp.net', '');
    }
  }

  const pushName = msg?.pushName || '';
  const timestamp = msg?.messageTimestamp;

  let type = 'unknown';
  let text = '';
  let quotedText = '';

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

    const contextInfo = message.extendedTextMessage?.contextInfo ||
      message.imageMessage?.contextInfo ||
      message.videoMessage?.contextInfo;

    if (contextInfo?.quotedMessage) {
      const qm = contextInfo.quotedMessage;
      quotedText = qm.conversation ||
        qm.extendedTextMessage?.text ||
        qm.imageMessage?.caption ||
        qm.videoMessage?.caption ||
        '';
    }
  }

  return {
    id,
    from,
    pushName,
    timestamp,
    type,
    text,
    isGroup,
    groupJid: isGroup ? remoteJid : '',
    participant,
    quotedText,
  };
}
