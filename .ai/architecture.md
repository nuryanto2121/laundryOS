# Architecture

NodeJS

Responsibilities

- WhatsApp Gateway
- Baileys
- Receive Message
- Send Message

NodeJS must NEVER contain business logic.

Business logic belongs to Go.

---

Go Backend

Responsibilities

- Customer
- Order
- Pickup
- Payment
- Workflow
- Notification
- AI

Go Backend is the source of truth.

---

Database

PostgreSQL

Redis

---

Frontend

Vue 3

---

Communication

Customer

↓

WhatsApp

↓

Baileys

↓

REST API

↓

Go Backend

↓

Database

↓

REST Response

↓

Baileys

↓

Customer