# Decision Log

## 2026-07-29

Decision:
Do not use Controller layer.

Reason:
Gateway is intentionally minimal.

---

Decision:
Inline validation is preferred.

Reason:
Current architecture doesn't require Validator abstraction.

---

Decision:
Unknown message type maps to "unknown".

Reason:
Future milestones will expand support.