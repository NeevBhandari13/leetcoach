# ADR 002: Cloud SQL Auth Proxy for Database Connectivity

**Status:** Accepted

## Context

The Go backend running on Cloud Run needs to connect to the Cloud SQL (PostgreSQL) instance. There are two ways to reach a Cloud SQL instance from Cloud Run:

1. **Public IP with an authorised network or SSL** — the database is reachable over the internet; callers authenticate with a username/password and optionally an SSL certificate.
2. **Cloud SQL Auth Proxy** — GCP runs a sidecar process alongside the Cloud Run container that opens a local Unix socket. The backend connects to that socket as if it were a local Postgres server; the proxy handles authentication via IAM and encrypts the connection transparently.

## Decision

Use the **Cloud SQL Auth Proxy via Unix socket**. Cloud Run manages the proxy automatically when the instance connection name is specified in the Cloud Run service configuration.

The DSN passed to `db.Open` in production takes the form:

```
postgres://USER:PASSWORD@localhost/DATABASE?host=/cloudsql/PROJECT:REGION:INSTANCE
```

The `?host=` parameter tells the `lib/pq` driver to connect via the Unix socket at that path rather than TCP.

## How It Works

```mermaid
sequenceDiagram
    participant App as Go backend
    participant Proxy as Cloud SQL Auth Proxy (sidecar)
    participant IAM as GCP IAM
    participant DB as Cloud SQL (PostgreSQL)

    App->>Proxy: connect via Unix socket\n(/cloudsql/project:region:instance)
    Proxy->>IAM: authenticate using\nCloud Run service account
    IAM-->>Proxy: short-lived token
    Proxy->>DB: encrypted TCP connection\nusing token
    DB-->>Proxy: query results
    Proxy-->>App: query results
```

The application code is unaware of the proxy. From Go's perspective it is connecting to a local Postgres instance.

## Why This Approach

| Concern | Public IP approach | Auth Proxy approach |
|---------|--------------------|---------------------|
| Authentication | Username + password + IP allowlist | IAM-based (no IP allowlist) |
| Network exposure | DB reachable from the internet | No public DB exposure |
| Encryption | Requires manual SSL cert management | Handled automatically by the proxy |
| Secrets to manage | DB password + SSL certs | DB password only (IAM handles the rest) |
| Code changes | None | None (transparent to `lib/pq`) |

The proxy eliminates the need to manage IP allowlists and SSL certificates. It also means the database has no public IP exposure — even if the DB password were leaked, there is no network path to reach it without the IAM-authenticated proxy.

## Local Development

In local development the backend uses a standard `postgres://` DSN from the `.env` file, connecting directly to a local Postgres instance. No proxy is involved. This means:

- The DSN format differs between local and production, which is expected.
- The `godotenv.Load("../.env")` call in `cmd/main.go` is a no-op on Cloud Run (the file doesn't exist in the container) — the runtime DSN comes from a Cloud Run environment variable populated from Secret Manager.

## Consequences

- The Cloud Run service account must be granted the `Cloud SQL Client` IAM role.
- The instance connection name (`PROJECT:REGION:INSTANCE`) must be specified in the Cloud Run service configuration (`--add-cloudsql-instances` flag or equivalent).
- Switching to a non-GCP Postgres host (e.g. Supabase, self-managed VM) requires only a DSN change — the Go code itself is unaffected.
