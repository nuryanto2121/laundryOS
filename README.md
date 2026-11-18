# LaundryOS

> Internal Laundry Management System for FAN Laundry Station.

LaundryOS automates daily laundry operations — starting with WhatsApp communication — and will evolve into a complete Laundry Management System.

---

## Architecture

```
Customer
  └─▶ WhatsApp
        └─▶ Gateway (Node.js + Baileys)
              └─▶ API (Go)
                    └─▶ PostgreSQL
                    └─▶ Redis
              ◀─────────────────────
        ◀─────────────────────
  ◀─────────────────────
```

| Service | Runtime | Role | Port |
|---------|---------|------|------|
| `api` | Go 1.24 | Business logic, REST API | `8080` |
| `gateway` | Node.js 22 | WhatsApp Gateway (Baileys) | `3000` |
| `dashboard` | Vue 3 + Vite | Internal ops UI | `5173` |
| `postgres` | PostgreSQL 16 | Primary database | `5432` |
| `redis` | Redis 7 | Cache & sessions | `6379` |

---

## Prerequisites

| Tool | Version |
|------|---------|
| Docker | 24+ |
| Docker Compose | v2 |
| Go | 1.24+ |
| Node.js | 22 LTS |
| Make | Any |
| golangci-lint | Latest (`brew install golangci-lint`) |

---

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/nuryanto2121/laundryos.git
cd laundryos

# 2. Create environment file
cp .env.example .env
# Edit .env and set POSTGRES_PASSWORD and other secrets
# NOTE: .env must exist before running make up.
# Docker Compose requires the file to be present (env_file directive).

# 3. Start all services
make up

# 4. Verify services are running
make ps
```

---

## Health Checks

```bash
# Go API Health
curl http://localhost:8080/health

# Swagger Interactive API Documentation
open http://localhost:8080/swagger/index.html

# Node Gateway
curl http://localhost:3000/health

# Vue Dashboard
open http://localhost:5173
```

Expected response:
```json
{
  "status": "ok",
  "service": "laundryos-api",
  "version": "0.1.0",
  "timestamp": "2026-07-27T09:00:00Z"
}
```

---

## Developer Workflow

### All services via Docker
```bash
make up       # Start everything
make logs     # Follow all logs
make down     # Stop everything
make restart  # Rebuild and restart
make ps       # Check container status
```

### Build & Test
```bash
make build             # Build all Docker images
make test              # Run all unit tests
make test-integration  # Run automated E2E integration & regression tests
make lint              # Run all linters
make fmt               # Format all code
make swagger           # Regenerate Swagger docs from Go code
```

### Local development (without Docker)
```bash
# Run each service locally — requires local Go/Node env
make api        # http://localhost:8080
make gateway    # http://localhost:3000
make dashboard  # http://localhost:5173
```

### Clean up
```bash
make clean    # Remove containers, volumes, and build artifacts
```

---

## Project Structure

```
laundryOS/
├── apps/
│   ├── api/          # Go backend (business logic)
│   ├── gateway/      # Node.js WhatsApp gateway
│   └── dashboard/    # Vue 3 internal dashboard
├── deployments/      # Reserved for k8s / CI configs
├── docs/             # Architecture & business documentation
├── spec/             # API, database, and workflow specifications
├── scripts/          # Utility scripts
├── docker-compose.yml
├── Makefile
└── .env.example
```

---

## Current Phase

**Phase 1 — Core Foundation**

- Sprint 00 — Project Foundation (Completed)
- Sprint 01 — WhatsApp Gateway Foundation (Completed)
- Sprint 02 — Backend Foundation (Completed)
- Sprint 03 — FAQ Bot (Next)

---

## License

Private. Internal use only — FAN Laundry Station.
