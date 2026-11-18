# Release Notes

## v0.4.0 — Sprint 04: Pickup Management

Status: **Completed** (Sprint 04 Finished)

Date:
2026-11-18

### Completed

- Implemented interactive multi-turn conversation state machine (`IDLE`, `AWAITING_ADDRESS`, `AWAITING_SERVICE`, `AWAITING_SCHEDULE`, `AWAITING_CONFIRMATION`).
- Built thread-safe session store (`MemorySessionStore`) with 15-minute TTL expiration and background cleaner loop.
- Added returning customer memory: auto-retrieves stored address and offers instant confirmation.
- Implemented customer auto-registration and address updates on confirmed pickup requests.
- Integrated order creation with status `WAITING_PICKUP` and deterministic sequence numbering.
- Implemented `NotificationService` broadcasting structured pickup ticket alerts to the Courier/Staff group.
- Added graceful cancellation flow (`BATAL` / `CANCEL`) that immediately clears active session states.
- Implemented Order Tracking & Completion Schedule feature (`STATUS` / `CEK STATUS` / explicit `ORD...` lookup) with humanized status progression and duration estimation.
- Added Gateway `GET /groups` endpoint to retrieve all participating WhatsApp groups and their JIDs for easy configuration.
- Expanded integration test suite (`scripts/test-integration.sh`) to 28 test cases with a 100% pass rate.
- Sprint 04 is officially concluded.

### Verification

- Multi-turn Turn 1 to Turn 5 verified live against Docker container ✔
- Customer `Bayu` automatically registered in `customers` table with `CST` code ✔
- Order created in PostgreSQL `orders` table with status `WAITING_PICKUP` and `pickup_required: true` ✔
- Returning customer query recognized saved address instantly ✔
- Check status query by customer phone returned active order status and estimation ✔
- Check status query by explicit order number (`ORD261118000001`) verified ✔
- Gateway `GET /groups` returned HTTP 200 ✔
- Sesi cancellation confirmed to clear memory store cleanly ✔
- Unit test suite: 88 passed in 13 packages ✔
- Integration test suite: 28 passed, 0 failed ✔

---

## v0.3.0 — Sprint 03: WhatsApp FAQ Bot

Status: **Completed** (Sprint 03 Finished)

Date:
2026-11-18

### Completed

- Implemented rule-based FAQ knowledge base supporting all 6 business scopes (Price, Service, Operating Hours, Address & Maps, Delivery Coverage, and Pickup Procedures).
- Built token-based weighted keyword matcher with text normalization (punctuation stripping, case folding, exact and phrase weighting).
- Implemented Dynamic Knowledge Base Provider querying live active laundry services from PostgreSQL `services` table, with seamless in-memory fallback.
- Added structured audit logging for unrecognized customer queries (`slog.Warn`) to support business intelligence and bot training.
- Created HTTP Bot Webhook endpoint `POST /bot/messages` with request validation and asynchronous WhatsApp Gateway dispatcher.
- Updated WhatsApp Gateway (`apps/gateway/src/whatsapp/handler.js`) to automatically forward non-group incoming text messages to the backend bot endpoint.
- Updated Swagger and OpenAPI specifications (`spec/api/openapi.yaml`).
- Expanded automated integration test suite (`scripts/test-integration.sh`) to 20 test cases, achieving 100% pass rate.
- Benchmark: Response time is ~50ms (well within the < 2s SLA requirement).
- Sprint 03 is officially concluded.

### Verification

- Live price query returned dynamic PostgreSQL active services list with IDR formatting ✔
- Live hours, address, pickup, and coverage queries matched categories accurately ✔
- Unrecognized query logged structured audit payload and returned friendly fallback menu ✔
- Empty body properly rejected with HTTP 400 validation error ✔
- Integration test suite: 20 passed, 0 failed ✔
- Unit test suite: 84 passed, 0 failed ✔

---

## v0.2.0 — Sprint 02 / Milestone 6: API Standardization & Test Suite

Status: **Completed** (Sprint 02 Finished)

Date:
2026-11-18

### Completed

#### Milestone 6 — API Standardization
- Implemented declarative request validation using `go-playground/validator/v10`
- Integrated validation rules on Customer (`full_name`, `phone`, `email`, `status`) and Order (`customer_id`, `service_id`, `total_weight`, `total_price`) models
- Standardized error response format across all endpoints
- Implemented structured JSON request logging middleware via `slog` (method, path, status, latency, remote address)
- Built automated End-to-End Integration & Regression Test Suite (`scripts/test-integration.sh` and `make test-integration`)
- Updated README.md with comprehensive guides and Swagger references
- Unit tests pass with 100% success (39/39 tests)
- Integration & Regression test suite passes with 100% success (17/17 tests)
- `golangci-lint` passes with 0 issues
- Sprint 02 (Backend Foundation) is officially concluded

### Verification

- Validation rejected empty customer body with HTTP 400 ✔
- Validation rejected invalid UUID order creation with HTTP 400 ✔
- Duplicate phone registration rejected with HTTP 400 ✔
- Full customer lifecycle validated end-to-end ✔
- Full order lifecycle with sequence generator validated end-to-end ✔
- Structured request logging confirmed in Docker container logs ✔
- All 17 automated regression assertions passed ✔

---

## v0.2.0 — Sprint 02 / Milestone 5: Swagger Documentation

Status: **Completed**

Date:
2026-11-18

### Completed

#### Milestone 5 — Swagger Documentation
- Integrated `swaggo` (`swag` v1.16.4 and `github.com/swaggo/http-swagger/v2`)
- Added declarative Swaggo annotations for all endpoints across Customer, Order, and Health modules
- Added `model.ErrorResponse` to document standard JSON error response schemas
- Exposed interactive Swagger UI at `/swagger/index.html` and raw JSON spec at `/swagger/doc.json`
- Added `make swagger` target to `Makefile`
- Synced `spec/api/openapi.yaml` with generated Swagger OpenAPI specification
- `golangci-lint` passes with 0 issues
- Unit tests pass with 100% success (39/39 tests)
- Docker container build and healthcheck pass

### Verification

- GET /swagger/index.html returns HTTP 200 OK with interactive Swagger UI ✔
- GET /swagger/doc.json returns HTTP 200 OK with valid OpenAPI 2.0 schema ✔
- All Customer endpoints (`POST`, `GET`, `PUT`, `PATCH /status`, `DELETE`) documented ✔
- All Order endpoints (`POST`, `GET`, `PUT`, `DELETE`) documented ✔
- Health check endpoint (`GET /health`) documented ✔

---

## v0.2.0 — Sprint 02 / Milestone 4: Order Module

Status: **Completed**

Date:
2026-11-18

### Completed

#### Milestone 4 — Order Module
- Implemented Order CRUD REST API: Create, Get, List, Update, Delete
- Configured sqlc queries for order persistence (`CreateOrder`, `GetOrder`, `ListOrders`, `UpdateOrder`, `DeleteOrder`, `GetLatestOrderNumberToday`)
- Order domain model with deterministic order numbering (`ORD` + `YYMMDD` + 6-digit daily sequence)
- Order status state enum (`NEW`)
- Soft/hard deletion aligned with database schema
- Repository layer maps all pgx/pgtype database types safely
- Unit tests for Order service with mock repository (`service_test.TestOrderService_*`)
- `golangci-lint` passes with 0 issues
- Docker build and container verification passes

### Verification

- POST /orders creates order with deterministic sequence `ORD261118000001` ✔
- GET /orders lists active orders ✔
- GET /orders/{id} retrieves order detail ✔
- PUT /orders/{id} updates mutable order properties while preserving status and order number ✔
- DELETE /orders/{id} removes order record (returns HTTP 204, subsequent GET returns HTTP 404) ✔
- Unit test suite passes (39/39 tests) ✔
- golangci-lint run — 0 issues ✔

---

## v0.2.0 — Sprint 02 / Milestone 3: Customer Module

Status: **Completed**

Date:
2026-07-31

### Completed

#### Milestone 3 — Customer Module
- Implemented Customer CRUD REST API: Create, Get, List, Update, Delete
- Configured `sqlc` with `sqlc.yaml` and generated type-safe database access code
- Customer SQL queries: `CreateCustomer`, `GetCustomer`, `ListCustomers`, `UpdateCustomer`, `DeleteCustomer` (soft delete)
- Customer Model separated into `customer_request.go` and `customer_response.go`
- Repository layer owns all `pgtype` / `sqlc` mapping — no database types leak to Service or Handler
- Service layer operates entirely on domain models (pure Go types)
- Handler layer: HTTP-only — decode, call service, encode
- Fixed `pgtype.Numeric` conversion (`Scan` requires string, not float64)
- UUID formatted using `github.com/google/uuid` library
- Timestamps formatted using `time.RFC3339` stdlib constant
- Added `make sqlc` target to Makefile
- `golangci-lint` passes with 0 issues

### Verification

- Create Customer ✔
- Get Customer Detail ✔
- List Customers ✔
- Update Customer ✔
- Delete Customer (soft delete — excluded from List and Get) ✔
- Latitude / Longitude persisted correctly ✔
- `golangci-lint run` — 0 issues ✔

---

## v0.2.0 — Sprint 02 / Milestone 2: Database Migration

Status: **Completed**

Date:
2026-07-30

### Completed

#### Milestone 2 — Database Migration
- Configured database schema migration using `golang-migrate`
- Created the initial migration version `000001_init_schema`
- Ensured migration tables align with the single source of truth (`spec/database/schema.dbml`)
- Created the `customers` table with soft delete support
- Created the `services` table for laundry services catalogue
- Created the `orders` table with foreign keys referencing customers and services
- Updated `Makefile` with `migrate-up` and `migrate-down` targets for local development convenience

### Verification

- Migration `up` creates `customers`, `services`, and `orders` successfully ✔
- Migration `down` cleanly drops all tables in reverse dependency order ✔
- Migration tables strictly reflect column precisions from the spec ✔

---

## v0.2.0 — Sprint 02 / Milestone 1: Backend Initialization

Status: **Completed**

Date:
2026-07-30

### Completed

#### Milestone 1 — Backend Initialization
- Initialized Go Backend using Chi HTTP router
- Configuration loading via cleanenv (env vars + `.env` file)
- PostgreSQL connection pool via pgx/v5
- Database connection verified on startup with ping
- Graceful shutdown on SIGINT / SIGTERM
- Structured JSON logging via slog
- GET /health endpoint returning HTTP 200
- Docker image builds successfully

### Verification

- Backend starts successfully ✔
- GET /health returns HTTP 200 ✔
- Database connected log confirmed ✔
- Graceful shutdown (Ctrl+C) confirmed ✔
- golangci-lint passes ✔
- Docker build passes ✔

---

## v0.1.0 — Sprint 01: WhatsApp Gateway Foundation

Status: **Completed**

Date:
2026-07-30

### Completed

#### Milestone 1 — Baileys Initialization
- Initialize WhatsApp Gateway using Baileys
- QR Code login support
- WhatsApp connection lifecycle
- Connection state logging

#### Milestone 2 — Session Persistence
- Persist WhatsApp session using `useMultiFileAuthState`
- Configurable session path via `SESSION_PATH`
- Reuse existing session after restart
- QR displayed only for first login
- Improved startup logging
- Removed legacy session references

#### Milestone 3 — Message Receiver
- Subscribe to `messages.upsert`
- Ignore non-notify events
- Ignore group, broadcast, status, and self messages
- Normalize incoming WhatsApp messages
- Log incoming messages in a human-readable format
- Preserve native Baileys timestamp internally
- Minimal supported message types (`text`, `image`, `video`, `document`, `unknown`)

#### Milestone 4 — Send Message
- Added `POST /send-message`
- Validate request payload (`phone`, `message`)
- Send text messages through Baileys
- Log outgoing messages
- Return HTTP success response after message delivery

#### Milestone 5 — Gateway Integration
- Enhanced `GET /health`
- Expose WhatsApp connection state
- Standardized API error responses
- Added endpoint-level request logging
- Improved request validation
- Gateway API ready for backend integration

### Verification

- WhatsApp connection established successfully ✔
- Session persistence verified ✔
- Incoming message logging verified ✔
- Outgoing message delivery verified ✔
- Health endpoint verified ✔
- API validation verified ✔
- ESLint passed ✔

### Next Version

**v0.2.0 — Sprint 02: Go Backend Integration**

---

## v0.1.0 — Sprint 00: Project Foundation

Status: **Completed**

### Completed

- Repository created
- AI documentation
- Project structure
- Docker Compose
- Go API scaffold
- Node.js Gateway scaffold
- Vue Dashboard scaffold
- Development tooling

---

## v0.0.1

Status: **Completed**

### Completed

- Initial repository