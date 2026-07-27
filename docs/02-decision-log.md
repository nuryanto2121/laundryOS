# Decision Log

## 2026-07

### Architecture

Decision

Use Monorepo architecture.

Reason

Simpler deployment and easier maintenance.

---

Decision

Use Golang for Backend.

Reason

High performance, low memory usage, simple deployment, suitable for long-running services.

---

Decision

Use Node.js + Baileys only as WhatsApp Gateway.

Reason

Baileys has the best ecosystem for unofficial WhatsApp integration.

Business logic remains in Go Backend.

---

Decision

Use PostgreSQL.

Reason

Reliable relational database with excellent ecosystem.

---

Decision

Use Redis.

Reason

Support caching, session management, and future event processing.

---

Decision

Keep Edoby POS during Phase 1.

Reason

Reduce migration risk.

Focus on WhatsApp automation first.