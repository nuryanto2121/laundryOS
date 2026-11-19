# Sprint 04 - Pickup Management & Staff Operations

## Goal

Allow customers to request laundry pickup through WhatsApp and enable staff to manage pickups, weighing, and order completion directly via WhatsApp staff group.

## Business Problem

Pickup requests, courier dispatches, and in-store weighing were handled manually without systematic tracking or automated notifications.

## Scope

- Customer Pickup Request & Multi-Turn Conversation
- Customer Registration & Geocoding/Address Collection
- Dynamic Service Selection & SLA Completion Estimation (3 hours for Self Service, Daily for regular)
- Order Creation with single-tap copy format (`#261119xxxxxx`)
- Courier Claim & Staff Group Broadcast Alerts
- Staff Commands: Weighing (`TIMBANG #... <weight>`) and Completion (`SELESAI #...`)
- Automatic Customer Status Updates via WhatsApp

## User Stories

1. As a customer, I want to request laundry pickup via WhatsApp so I don't have to visit the store.
2. As a staff / courier, I want to claim pickup jobs in the staff group, record laundry weight, and mark orders complete via simple chat commands.

## Acceptance Criteria

- [x] Multi-turn conversation state machine manages pickup workflow cleanly.
- [x] Orders created with status `WAITING_PICKUP` and format `#261119xxxxxx`.
- [x] Alerts dispatched to courier staff group with order details and 1-click WhatsApp customer link.
- [x] Courier claim replies update status and alert customer.
- [x] Staff commands `TIMBANG` and `SELESAI` update order weight, total price, and lifecycle status.
- [x] All 28 integration test cases passing in `scripts/test-integration.sh`.
- [x] 109 Go unit tests passing.

## Deliverables

- Pickup Module & Session Store
- Courier Group Relay & Command Handlers
- Integration Test Suite (28/28 verified)
