# Architecture

This document defines the architecture and technology decisions for LaundryOS.

AI must follow this document and must not introduce alternative technologies unless explicitly approved.

---

# Backend

Language

- Go 1.25

Responsibilities

- Customer
- Order
- Pickup
- Payment
- Workflow
- Notification
- AI

Technology

- HTTP Router: Chi
- Database Driver: pgx/v5
- Query Generator: sqlc
- Migration: golang-migrate
- Configuration: cleanenv
- Logging: slog
- Validation: go-playground/validator
- API Documentation: swaggo
- Testing: testing + testify

Architecture

- Lightweight Clean Architecture
- Repository Pattern
- Manual Dependency Injection
- REST API

Rules

- Go Backend is the source of truth.
- Business logic belongs only in the Go Backend.
- Do not introduce an ORM.
- Use sqlc-generated queries for database access.
- Keep the architecture simple.
- Avoid unnecessary abstractions.

---

# WhatsApp Gateway

Language

- Node.js 22 LTS

Responsibilities

- WhatsApp connection
- Baileys session management
- Receive messages
- Send messages
- Forward requests to the Backend

Technology

- Express
- Baileys

Rules

- Gateway must NEVER contain business logic.
- Gateway is only a transport layer.
- Validation should only cover transport-level concerns.
- Business validation belongs to the Backend.

---

# Database

Primary Database

- PostgreSQL

Future Components

- Redis (Future Sprint)

Rules

- PostgreSQL is the system of record.
- Database schema changes must be managed through migrations.

---

# Frontend

Framework

- Vue 3

Responsibilities

- Dashboard
- Administration
- Reporting

Rules

- Frontend communicates only with the Backend.
- Frontend never communicates directly with the Gateway.

---

# Communication Flow

Customer

↓

WhatsApp

↓

Baileys Gateway

↓

REST API

↓

Go Backend

↓

PostgreSQL

↓

REST Response

↓

Baileys Gateway

↓

Customer

---

# Architecture Principles

- Keep the implementation minimal.
- Keep the diff as small as possible.
- Reuse the existing architecture.
- Do not introduce new frameworks without approval.
- Do not introduce new layers unless required.
- Prefer standard library when sufficient.
- Follow the active sprint exactly.