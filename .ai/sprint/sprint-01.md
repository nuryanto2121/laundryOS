# Sprint 01 - WhatsApp Gateway Foundation

## Sprint Objective

Build the initial WhatsApp Gateway capable of connecting to WhatsApp using Baileys.

At the end of Sprint 01, LaundryOS must be able to:

- Connect to WhatsApp
- Generate QR Code
- Persist login session
- Automatically reconnect
- Receive incoming private messages
- Forward incoming messages to Go Backend
- Receive response from Go Backend
- Send response back to WhatsApp

No business logic shall exist inside the Gateway.

The Gateway is only responsible for WhatsApp communication.

---

# Sprint Scope

## In Scope

- Baileys initialization
- QR Login
- Session persistence
- Auto reconnect
- Receive incoming messages
- Send outgoing messages
- HTTP communication with Go Backend
- Logging
- Configuration

## Out of Scope

- Customer Management
- Laundry Orders
- Database
- Redis
- Authentication
- AI
- POS
- Dashboard Features
- Payment
- Inventory

---

# Architecture Rules

Node.js Gateway MUST:

- communicate with WhatsApp only
- communicate with Go Backend only
- never contain business logic
- never access PostgreSQL directly
- never access Redis directly

Business decisions always belong to Go Backend.

---

# Milestone 01 - Baileys Initialization

## Status

Planned

## Goal

Create a minimal WhatsApp Gateway capable of connecting to WhatsApp.

## Tasks

- Install Baileys
- Initialize socket
- Load configuration
- Display QR Code in terminal
- Connect successfully
- Display connection state in console

## Deliverables

- Baileys initialized
- QR displayed
- Connected event handled
- Basic logger

## Out of Scope

- Session persistence
- Message handler
- API integration

## Acceptance Criteria

Developer runs:

npm run dev

Expected result:

- QR Code appears
- QR can be scanned
- WhatsApp becomes connected
- Gateway does not crash

Definition of Done

- Code builds successfully
- No lint errors
- Docker image builds successfully

---

# Milestone 2 — Session Persistence

## Goal

Implement persistent WhatsApp authentication so QR Code is required only once.

## In Scope

- Persist Baileys auth state using useMultiFileAuthState
- Store session under:
  apps/gateway/storage/session/
- Automatically reuse existing session
- Do not display QR when valid session exists
- Separate session management from client initialization
- Keep current reconnect logic for restartRequired only
- Improve startup logs

## Out of Scope

- Receive messages
- Send messages
- Webhook/API integration
- Business logic
- Database
- Redis
- Order processing

## Deliverables

- storage/session/
- session.js
- Updated client.js
- Updated README

## Acceptance Criteria

- First run requires QR scan
- Session files are generated
- Restart gateway
- Gateway reconnects automatically
- QR is NOT displayed
- Health endpoint still works
- Docker build passes
- Lint passes

## Definition of Done

✓ npm run dev
✓ Login only once
✓ Restart without scanning
✓ Session survives restart
✓ Documentation updated

---

# Milestone 03 - Receive Incoming Message

## Status

Planned

## Goal

Receive customer messages.

## Tasks

Subscribe to:

messages.upsert

Extract:

- Phone Number
- Push Name
- Message
- Timestamp
- Message Type

Ignore:

- Group messages
- Broadcast messages
- Status updates

Log received messages.

## Deliverables

Incoming messages appear in console log.

## Out of Scope

- Reply
- Database
- API

## Acceptance Criteria

Customer sends:

Hello

Console shows:

Phone

Message

Timestamp

Gateway continues running.

Definition of Done

No panic.

No unhandled error.

---

# Milestone 04 - Send Message

## Status

Planned

## Goal

Expose a minimal HTTP endpoint that sends a WhatsApp text message through the existing Baileys connection.

The Gateway acts only as a transport layer and contains no business logic.

---

## Scope

### In Scope

- Create HTTP endpoint:
  ```
  POST /send-message
  ```

- Accept JSON request:
  ```json
  {
    "phone": "6281234567890",
    "message": "Hello"
  }
  ```

- Validate required fields:
  - `phone`
  - `message`

- Convert phone number to WhatsApp JID:
  ```
  6281234567890
      ↓
  6281234567890@s.whatsapp.net
  ```

- Send message using the active Baileys socket.

- Return appropriate HTTP response.

- Log outgoing messages.

---

### Out of Scope

- Go Backend integration
- Authentication / Authorization
- Rate limiting
- Retry mechanism
- Scheduled messages
- Media messages
- Broadcast messages
- Database
- Delivery status tracking
- Queue / Worker

---

## Proposed Changes

### MODIFY

#### `apps/gateway/src/index.js`

- Register:
  ```
  POST /send-message
  ```
- Delegate request handling to a dedicated HTTP handler.

---

#### `apps/gateway/src/whatsapp/client.js`

Expose a function:

```js
sendMessage(phone, message)
```

Responsibilities:

- Verify WhatsApp connection is active.
- Convert phone number into WhatsApp JID.
- Send text message via Baileys.
- Return send result.
- No HTTP-specific logic.

---

### NEW

#### `apps/gateway/src/http/sendMessage.js`

Responsibilities:

- Parse HTTP request.
- Validate request body.
- Invoke `client.sendMessage()`.
- Return HTTP response.
- No WhatsApp-specific implementation.

---

## API Contract

### Request

**POST** `/send-message`

```json
{
  "phone": "6281234567890",
  "message": "Hello"
}
```

---

### Success Response

HTTP 200

```json
{
  "success": true
}
```

---

### Validation Error

HTTP 400

```json
{
  "success": false,
  "message": "phone is required"
}
```

or

```json
{
  "success": false,
  "message": "message is required"
}
```

---

### WhatsApp Not Connected

HTTP 503

```json
{
  "success": false,
  "message": "WhatsApp is not connected"
}
```

---

### Internal Error

HTTP 500

```json
{
  "success": false,
  "message": "Failed to send message"
}
```

---

## Logging

### Success

```
Outgoing Message

To        : 6281234567890
Type      : text
Text      : Hello
Status    : Sent
────────────────────────────
```

### Failure

```
Failed Sending Message

To        : 6281234567890
Reason    : <error>
────────────────────────────
```

---

## Deliverables

- HTTP endpoint implemented.
- WhatsApp text message successfully sent.
- Outgoing message logged.
- Request validation implemented.

---

## Acceptance Criteria

- POST `/send-message` endpoint exists.
- Accepts JSON request body.
- Validates required fields.
- Converts phone number to WhatsApp JID.
- Uses the existing Baileys connection.
- Customer receives the message.
- Returns HTTP 200 on success.
- Returns HTTP 400 for invalid requests.
- Returns HTTP 503 when WhatsApp is disconnected.
- Returns HTTP 500 for unexpected errors.
- Gateway remains running after sending messages.
- Lint passes.
- Docker build passes.

---

## Definition of Done

- HTTP endpoint implemented.
- Message successfully delivered.
- Proper HTTP response returned.
- Manual verification completed.
- Code reviewed.
- Lint passes.
- Docker build passes.

---

## Verification Plan

### Positive Test

Request:

```json
{
  "phone": "6281234567890",
  "message": "Hello"
}
```

Expected:

- HTTP 200
- Customer receives the message
- Gateway logs the outgoing message

---

### Validation Test

Missing `phone`:

```json
{
  "message": "Hello"
}
```

Expected:

- HTTP 400

---

Missing `message`:

```json
{
  "phone": "6281234567890"
}
```

Expected:

- HTTP 400

---

### Connection Test

Disconnect WhatsApp.

Call:

```
POST /send-message
```

Expected:

- HTTP 503
- No crash

---

### Error Handling Test

Force `sock.sendMessage()` to throw an exception.

Expected:

- HTTP 500
- Error logged
- Gateway remains running

---

# Milestone 05 - Gateway Integration

## Status

Planned

## Goal

Prepare the WhatsApp Gateway to be consumed by external services.

This milestone focuses on stabilizing the Gateway API, improving API responses, and exposing a minimal interface that can later be called by the Go Backend.

No Go Backend implementation is included in this milestone.

---

## Tasks

### HTTP API

Review and finalize the Gateway HTTP interface.

Endpoints:

GET /health

POST /send-message

Standardize HTTP responses.

Example success:

```json
{
  "success": true
}
```

Example validation error:

```json
{
  "success": false,
  "error": "phone is required"
}
```

Example gateway unavailable:

```json
{
  "success": false,
  "error": "WhatsApp is not connected"
}
```

---

### Health Endpoint

Enhance the existing `/health` endpoint.

Return:

```json
{
  "status": "UP",
  "whatsapp": "CONNECTED"
}
```

or

```json
{
  "status": "UP",
  "whatsapp": "DISCONNECTED"
}
```

The WhatsApp status must come from `client.getState()`.

---

### API Validation

Validate:

- Content-Type
- phone
- message

Return proper HTTP status codes.

400

Invalid request.

503

Gateway not connected.

500

Unexpected internal error.

---

### Logging

Log every API request.

Example:

```
POST /send-message
```

Log response status.

Example:

```
POST /send-message -> 200
```

No request body should be logged.

---

### Documentation

Document:

- endpoint
- request
- response
- status code

inside the project README.

---

## Deliverables

Stable Gateway HTTP API.

Consistent response structure.

Improved health endpoint.

Documented API.

---

## Out of Scope

Go Backend.

Authentication.

Authorization.

Webhook.

Retry mechanism.

Rate limiting.

Queue.

Redis.

PostgreSQL.

Media messages.

Message delivery status.

---

## Acceptance Criteria

- GET /health returns Gateway status and WhatsApp status.
- POST /send-message keeps working.
- HTTP responses are standardized.
- Validation errors return HTTP 400.
- Disconnected Gateway returns HTTP 503.
- Unexpected errors return HTTP 500.
- Every request is logged.
- API is documented.
- Lint passes.
- Docker build passes.

---

## Definition of Done

Gateway exposes a stable HTTP API that is ready to be consumed by external services.

No business logic exists inside the Gateway.

The Gateway remains responsible only for transporting WhatsApp messages.

---

# Sprint Definition of Done

Sprint 01 is considered complete only if:

- QR Login works
- Session persists
- Auto reconnect works
- Incoming messages received
- Outgoing messages sent
- Gateway communicates with Go Backend
- All lint passes
- All tests pass
- Docker builds successfully

---

# Files Expected To Be Modified

apps/gateway/

apps/api/

docker-compose.yml

README.md

CURRENT.md

---

# Documentation Update

When a milestone is completed, update:

.ai/CURRENT.md

Update:

Current Milestone

Completed Milestones

Next Milestone

If Sprint 01 is completed, also update:

docs/03-release-note.md

docs/02-decision-log.md

---

# Implementation Rules

For every milestone:

1. Validate architecture.
2. Generate implementation plan.
3. Wait for approval.
4. Generate code.
5. Wait for review.
6. Apply approved fixes only.
7. Update CURRENT.md.
8. Stop.

Never continue to the next milestone automatically.

The developer explicitly controls the progression of the sprint.