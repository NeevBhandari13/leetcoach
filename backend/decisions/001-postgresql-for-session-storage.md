# ADR 001: PostgreSQL for Session Storage

**Status:** Accepted

## Context

Each interview is a session that accumulates state across multiple HTTP requests: the problem text, the current interview stage, and the full message history. This data needs to be accessible on every request.

## Decision

Use PostgreSQL as the single persistence layer for sessions and messages, with `database/sql` and `lib/pq`.

## Alternatives Considered

| Option | Why rejected |
|--------|-------------|
| In-memory map (Go struct) | Lost on server restart; can't scale to multiple processes |
| Redis | Extra infrastructure dependency; no relational joins for ordered message history |
| SQLite | Simpler setup, but file-locking issues under concurrent writes |

## Data Model

```
sessions                          messages
─────────────────────────         ──────────────────────────────────
id          UUID  PK              id           SERIAL PK
state       TEXT                  session_id   UUID  FK → sessions.id
problem_text TEXT                 role         TEXT  (user | assistant)
created_at  TIMESTAMPTZ           content      TEXT
                                  created_at   TIMESTAMPTZ
```

Messages are ordered by `(created_at, id)` to reconstruct conversation history in the correct sequence for every LLM call.

## Request Flow

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant SessionStore
    participant Postgres

    Client->>Handler: POST /sessions/:id/reply
    Handler->>SessionStore: Reply(sessionID, system, userMsg)
    SessionStore->>Postgres: INSERT user message
    SessionStore->>Postgres: SELECT all messages (ordered)
    SessionStore->>SessionStore: Call LLM with full history
    SessionStore->>Postgres: INSERT assistant reply
    SessionStore->>Handler: return raw LLM response string
    Handler->>Postgres: UPDATE sessions SET state = ...
    Handler->>Client: {message, state}
```

## Consequences

- Migrations run automatically at startup via `db.Migrate()` (golang-migrate, embedded SQL).
- `db.Seed()` is idempotent — safe to run on every restart.
- Full message history is fetched on every `Reply` call. This is fine at interview scale (tens of messages) but would need pagination at larger volumes.