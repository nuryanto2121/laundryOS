/**
 * WhatsApp Session
 *
 * Encapsulates all session-related filesystem operations.
 * client.js delegates session management here and stays focused on the socket lifecycle.
 */

import { useMultiFileAuthState } from '@whiskeysockets/baileys';
import { mkdir, rm } from 'node:fs/promises';

/**
 * Default session directory, relative to process.cwd().
 * Override via SESSION_PATH environment variable.
 */
export const DEFAULT_SESSION_PATH = 'storage/session';

/**
 * Loads WhatsApp auth state from disk.
 *
 * Creates the session directory if it does not exist.
 * When no credentials are found, Baileys will emit a QR code event.
 * When valid credentials exist, Baileys connects silently.
 *
 * @param {string} sessionPath - Absolute or relative path to the session directory.
 * @returns {Promise<{ state: object, saveCreds: Function }>}
 */
export async function load(sessionPath) {
  await mkdir(sessionPath, { recursive: true });
  return useMultiFileAuthState(sessionPath);
}

/**
 * Removes all files under the session directory.
 *
 * Reserved for future CLI/admin utilities.
 * Intentionally unused in Milestone 2.
 *
 * @param {string} sessionPath - Path to the session directory.
 * @returns {Promise<void>}
 */
export async function clear(sessionPath) {
  await rm(sessionPath, { recursive: true, force: true });
}
