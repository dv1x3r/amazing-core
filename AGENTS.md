# AGENTS.md

Guidelines for AI agents working on the Amazing Core repository.

## Project overview

- Amazing Core is an open-source server emulator for **Amazing World** (MMO, shut down 2018)
- Original client: Unity, communicates with server via a custom TCP-based GSF protocol
- Amazing Core reimplements server-side behavior
- Many handlers and gameplay systems are **incomplete or placeholders** and must not be assumed to exist

## Client reference

- `Assembly-CSharp/` (if present) contains C# client code for:
  - Protocol message structures and field types
  - Request/response workflows
  - Verifying server implementation matches client expectations
- Message implementation:
  - Verify request/response types in `Assembly-CSharp/*.cs` first
  - Do not create handlers for unused or unimplemented messages

## Development environment

- Language: **Go 1.25** (pinned, do not modify go.mod or go.sum without permission)
- Database: **SQLite** (migrations via goose)
  - `CGO_ENABLED=0` → `modernc.org/sqlite` (default)
  - `CGO_ENABLED=1` → `github.com/mattn/go-sqlite3`
- Driver-specific logic:
  - `internal/lib/db/sqlite_nocgo.go` — `CGO_ENABLED=0`
  - `internal/lib/db/sqlite_cgo.go` — `CGO_ENABLED=1`

## Build & test

Server binary:

- `make build` — compile binary to `./build/server`
- `make run` — build & run server
- `make test` — run tests (verbose)
- `make generate` — run Go code generators (e.g., `stringer` for enums in `internal/gsf/types/`)

Database migrations (`./data_db/core.db`):

- `make db-status` — show migration status
- `make db-up` — apply all migrations
- `make db-up-by-one` — apply next migration
- `make db-down` — rollback last migration
- `make db-reset` — rollback all migrations
- `make db-create` — create new migration (interactive)

## Repository structure

### cmd/

Entry points, primary server: `cmd/server/main.go`

- Initialize the database and apply migrations
- Start HTTP API (default port 3000)
- Start TCP game server (default port 8182)
- **Agents**: avoid alternative entry points unless required

### data/

SQL schemas and migrations:

- `data/sql/base/` — base database schemas.
- `data/sql/updates/` — incremental migrations.
- **Append-only** migrations; never modify applied files

### internal/api/

HTTP API & admin dashboard backend:

- **Keep API handlers thin**; delegate logic to services

### internal/game/

Game server message routing.

### internal/network/

TCP game server protocol handling:

- Key subfolders and files:
  - `gsf/` — networking protocol implementation
    - `server.go` — TCP server bootstrap
    - `router.go` — message routing
    - `protocol.go` — protocol interface and slice/map helpers
    - `types/` — DTO types and enums, used within messages
    - `messages/` — DTO message types that implement client-server interaction, may reference `types/`
    - `bitprotocol/` — bit-level protocol encoding and decoding
- Use `ReadSlice`/`WriteSlice`, `ReadMap`/`WriteMap`, `ReadObject`/`WriteObject` helpers
- **Do not break wire compatibility**
- **Do not modify existing message structures** unless verified in client code or explicitly instructed.

### internal/services/

Business logic, database access, reusable domain operations.

### internal/config/

Centralized config from `config.json` (ports, DB paths, logger)

### internal/lib/

- Utilities: `db/`, `logger/`, `wrap/` etc.
- Reusable across the codebase

### tools/

- Developer utilities not part of server runtime
- **Do not import into production code**

### web/

- Embedded admin dashboard assets
- Plain JavaScript; **no build tools or package managers**

## Engineering guidelines

These rules apply across the entire codebase unless explicitly stated otherwise.

### Agent behavior

Agents may:

- Implement requested features or message handlers
- Fill placeholders only if behavior is confirmed by client code or instructions
- Refactor locally for clarity without changing behavior

Agents must NOT:

- Invent gameplay rules, formulas, progression, or economy logic
- Guess protocol semantics or message usage
- Run any build/test commands without explicit permission or instruction

For any non-trivial change:

- State necessity
- Cite references (client code, pattern, instruction)
- List assumptions

If unclear:

- Prefer TODOs/stubs over speculation
- Ask a single, specific clarification

Agents are **implementation assistants**, not designers.

### Code principles

- Prefer **clarity, correctness, minimal scope**
- Follow **existing architecture**; extend rather than duplicate
- Keep **Go idiomatic**, small, and incremental
- When in doubt, do nothing and ask.
