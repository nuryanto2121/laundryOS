/**
 * Logger
 *
 * Lightweight structured logger wrapping the built-in console.
 * API: logger.info(), logger.warn(), logger.error()
 *
 * The implementation can be swapped (e.g. to pino) without changing callers,
 * as long as the same method signatures are preserved.
 */

const LEVELS = { debug: 0, info: 1, warn: 2, error: 3 };

const configuredLevel = (process.env.LOG_LEVEL || 'info').toLowerCase();
const minLevel = LEVELS[configuredLevel] ?? LEVELS.info;

/**
 * Format a log entry as a single line.
 * @param {string} level
 * @param {string} message
 * @param {object} [meta]
 */
function format(level, message, meta) {
  const timestamp = new Date().toISOString();
  const base = `${timestamp} [${level.toUpperCase()}] ${message}`;
  return meta ? `${base} ${JSON.stringify(meta)}` : base;
}

const logger = {
  debug(message, meta) {
    if (LEVELS.debug >= minLevel) {
      console.debug(format('debug', message, meta));
    }
  },

  info(message, meta) {
    if (LEVELS.info >= minLevel) {
      console.info(format('info', message, meta));
    }
  },

  warn(message, meta) {
    if (LEVELS.warn >= minLevel) {
      console.warn(format('warn', message, meta));
    }
  },

  error(message, meta) {
    if (LEVELS.error >= minLevel) {
      console.error(format('error', message, meta));
    }
  },
};

export default logger;
