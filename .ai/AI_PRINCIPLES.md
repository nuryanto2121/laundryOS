# AI Principles

This document defines the engineering principles that the AI must follow when contributing to the LaundryOS project.

These principles apply to every sprint and milestone unless explicitly overridden by the developer.

---

# 1. Stay Within the Current Milestone

Implement only what is required by the active milestone.

Do not add features that belong to future milestones.

Avoid scope creep.

---

# 2. Keep the Implementation Minimal

Prefer the simplest implementation that satisfies the acceptance criteria.

Avoid unnecessary complexity.

Do not over-engineer.

---

# 3. Keep Diffs Small

Modify existing files whenever reasonable.

Do not create new files unless they provide clear value or are explicitly required.

Keep pull requests and commits easy to review.

---

# 4. Reuse Existing Architecture

Respect the current project structure.

Reuse existing modules before introducing new ones.

Do not introduce new architectural layers unless explicitly requested.

Examples:

- Controllers
- Services
- Repositories
- DTOs
- Validators
- Middleware

are NOT added unless required by the milestone.

---

# 5. Clean Separation of Responsibilities

Each module should have a single responsibility.

Examples:

- client.js manages the Baileys socket.
- session.js manages session persistence.
- mapper.js translates external payloads.
- handler.js processes incoming events.

Business logic must not leak into infrastructure modules.

---

# 6. Prefer Defensive Programming

The Gateway must never crash because of malformed or unexpected external data.

Use:

- optional chaining
- null checks
- safe fallbacks

Avoid throwing exceptions for recoverable situations.

---

# 7. Do Not Optimize Prematurely

Do not introduce:

- caching
- queues
- retries
- concurrency improvements
- abstractions
- performance optimizations

unless the current milestone explicitly requires them.

---

# 8. Follow Existing Project Style

Match the existing codebase.

Maintain consistency in:

- naming
- formatting
- logging
- folder structure
- error handling

Avoid introducing a different coding style.

---

# 9. Logging Should Be Useful

Logs should help developers understand what the Gateway is doing.

Logs should be:

- concise
- readable
- structured

Avoid excessive debug output.

---

# 10. Preserve Internal Models

Internal models should preserve original data whenever practical.

Formatting should happen only at presentation boundaries (logs or HTTP responses), not inside mapping logic unless explicitly required.

---

# 11. Documentation by Request

Do not modify documentation unless:

- explicitly requested
- included in the milestone scope

Avoid changing README or other documents as part of unrelated implementation work.

---

# 12. Manual Verification Before Completion

A milestone is NOT considered completed after implementation alone.

Completion requires:

- implementation finished
- lint passes
- manual verification passes
- developer approval

Only then may CURRENT.md be updated.

---

# 13. Never Assume Future Requirements

If a feature is planned for a later milestone:

Do not implement it now.

Leave future enhancements for their designated milestone.

---

# 14. Prefer Explicit Code

Favor readability over cleverness.

Simple code is preferred over generic abstractions.

The project values maintainability more than flexibility.

---

# 15. Ask Before Expanding Scope

If an improvement is beneficial but outside the current milestone:

Do not implement it.

Instead:

- mention it during the implementation plan
- explain why it may be valuable
- wait for developer approval

---
# 16. Plan First, Code Second

The AI must never start implementing immediately.

For every milestone:

1. Read AGENTS.md.
2. Read AI_WORKFLOW.md.
3. Read AI_PRINCIPLES.md.
4. Read CURRENT.md.
5. Read the active sprint document.
6. Produce an implementation plan.
7. Wait for developer approval.
8. Only then begin implementation.

Implementation must always follow the approved plan.

# Guiding Philosophy

The LaundryOS project follows an incremental development approach.

Every milestone should:

- solve one problem well
- remain easy to understand
- remain easy to review
- remain easy to test

Small, correct, and maintainable changes are preferred over ambitious implementations.