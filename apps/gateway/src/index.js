'use strict';

const express = require('express');

const PORT = process.env.PORT || 3000;
const APP_VERSION = process.env.APP_VERSION || '0.1.0';

const app = express();

app.use(express.json());

/**
 * GET /health
 * Standard health check endpoint.
 */
app.get('/health', (req, res) => {
  res.status(200).json({
    status: 'ok',
    service: 'laundryos-gateway',
    version: APP_VERSION,
    timestamp: new Date().toISOString(),
  });
});

app.listen(PORT, () => {
  console.log(`laundryos-gateway v${APP_VERSION} started on port ${PORT}`);
});
