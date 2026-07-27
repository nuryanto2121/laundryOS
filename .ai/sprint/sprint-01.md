# Sprint 01 - WhatsApp Gateway

## Goal

Implement a stable WhatsApp Gateway using Baileys.

## Business Problem

The business currently replies to customers manually using WhatsApp.

The system must be able to receive and send messages automatically.

## Scope

- QR Login
- Session Persistence
- Auto Reconnect
- Receive Private Message
- Receive Group Message
- Send Message

## User Stories

As the laundry owner,

I want the system to connect to WhatsApp,

so customer conversations can be automated.

## Acceptance Criteria

- QR Login works
- Session persists after restart
- Auto reconnect works
- Private messages received
- Group messages received
- Messages can be sent

## Technical Tasks

- Install Baileys
- Configure authentication
- Store session
- Implement reconnect logic
- Receive events
- Send messages
- Logging

## Deliverables

- WhatsApp Gateway

## Out of Scope

- AI
- Customer Database
- Orders
- Dashboard
- POS

## Definition of Done

- Gateway runs for 24 hours without manual intervention