# Coding Standard

Language

Backend

Go

Gateway

NodeJS

Never mix responsibilities.

Follow Clean Architecture.

Preferred structure

internal/

domain/

repository/

service/

handler/

middleware/

config/

Never access database directly from HTTP handlers.

Repository only accesses database.

Service contains business logic.

Use UUID for IDs.

Use human-readable order numbers.

Example

FLS2507270001

Always write readable code.

Prefer explicit code over magic.

Avoid unnecessary abstraction.

Avoid premature optimization.

Write unit tests whenever possible.