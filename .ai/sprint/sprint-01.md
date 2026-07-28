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

Send WhatsApp message programmatically.

## Tasks

Create internal endpoint.

POST

/send-message

Request

{
    "phone":"628xxxxxxxxxx",
    "message":"Hello"
}

Gateway sends message through Baileys.

## Deliverables

Message successfully delivered.

## Out of Scope

Go Backend.

## Acceptance Criteria

POST endpoint.

Customer receives message.

Definition of Done

HTTP returns success.

Message delivered.

---

# Milestone 05 - Gateway Integration

## Status

Planned

## Goal

Forward incoming WhatsApp messages to Go Backend.

Gateway becomes a dumb transport layer.

## Flow

Customer

↓

Gateway

↓

Go Backend

↓

Gateway

↓

Customer

## Tasks

Forward incoming message.

POST

/api/v1/messages

Receive backend response.

Reply to customer.

Backend response is temporary static text.

## Deliverables

Successful round trip.

## Out of Scope

Database.

Customer Management.

Laundry Orders.

AI.

## Acceptance Criteria

Customer sends:

Hello

Backend responds:

Welcome to FAN Laundry.

Customer receives reply.

Definition of Done

Round trip completed.

No business logic exists inside Gateway.

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