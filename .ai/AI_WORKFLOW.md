# AI Workflow

This document defines the standard development workflow for this project.

The AI must always follow this workflow unless explicitly instructed otherwise.

---

## Development Flow

For every milestone:

1. Read AGENTS.md.
2. Read CURRENT.md.
3. Read the active sprint document.
4. Produce an implementation plan.
5. Wait for approval.
6. Implement the approved plan.
7. Wait for manual verification by the developer.
8. Revise the implementation if issues are found.
9. Produce a walkthrough document summarizing the implementation and verification.
10. Wait for confirmation that the milestone is complete.
11. Only after confirmation:
    - update CURRENT.md
    - mark the milestone as completed

---

## Important Rules

- Never mark a milestone as completed before manual verification.
- Never update CURRENT.md before approval.
- Never generate git commits or tags unless requested.
- Keep the implementation within the current milestone.
- Do not introduce unnecessary abstractions.
- Keep the diff as small as possible.

---

## Definition of Done

A milestone is considered completed only after:

- Implementation finished.
- Lint passes.
- Manual verification passes.
- Developer confirms completion.
- CURRENT.md updated.