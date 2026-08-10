# Release Notes

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