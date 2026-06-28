# ADR 003: GCP Secret Manager for Credentials

**Status:** Accepted

## Context

The backend requires two sensitive values at runtime:

- `ANTHROPIC_API_KEY` — authenticates calls to the Anthropic API
- DB password — used in the Cloud SQL DSN

These values need to be available to the Cloud Run container without being committed to source control or visible to everyone with GCP console access.

GCP offers two ways to pass configuration to a Cloud Run container:
1. **Environment variables** set directly in the Cloud Run service configuration (visible in the console to anyone with project access)
2. **Secret Manager** — values are stored encrypted, access is controlled by IAM, and they are mounted into the container as environment variables or volume files at runtime

## Decision

Store both secrets in **GCP Secret Manager** and reference them in the Cloud Run service configuration. Cloud Run fetches and injects them as environment variables when the container starts. The Go code reads them with `os.Getenv()` — no Secret Manager SDK needed in the application.

```
Secret Manager
  └── ANTHROPIC_API_KEY  (version 1)
  └── DB_PASSWORD        (version 1)
         │
         │  injected at container startup
         ▼
Cloud Run environment
  └── ANTHROPIC_API_KEY=sk-ant-...
  └── DSN=postgres://user:PASSWORD@localhost/db?host=/cloudsql/...
```

## Why Secret Manager Over Plain Env Vars

| Concern | Plain Cloud Run env vars | Secret Manager |
|---------|--------------------------|----------------|
| Visibility in GCP console | Visible to anyone with Run Viewer role | Requires explicit Secret Accessor IAM grant |
| Audit trail | None | Every access is logged in Cloud Audit Logs |
| Rotation | Must redeploy to change | Update secret version; optionally trigger redeploy |
| Source control risk | Easy to accidentally commit | Values never leave Secret Manager unless explicitly fetched |

The Anthropic API key in particular has a direct cost impact if leaked — a compromised key could run up large API bills. Storing it in Secret Manager with a narrow IAM grant (only the Cloud Run service account can read it) limits the blast radius of any credential exposure.

## Access Control

The Cloud Run service account is granted the `Secret Manager Secret Accessor` role **scoped to individual secrets**, not project-wide. This means the service account can only read the two secrets it needs — not any other secrets that might be added to the project later.

## Local Development

Local development uses a `.env` file at the repo root (gitignored). `godotenv.Load` reads it at startup. The Secret Manager is not involved locally — developers set their own values in `.env`.

This means there are two credential paths:

| Environment | Source |
|-------------|--------|
| Local | `.env` file (gitignored) |
| Cloud Run | GCP Secret Manager → injected as env vars |

The Go code (`os.Getenv`) is identical in both cases.

## Consequences

- Rotating a secret (e.g. rolling the Anthropic API key) requires adding a new version in Secret Manager and redeploying Cloud Run to pick it up. This can be automated with a Cloud Run trigger on secret rotation.
- The Cloud Run service account must have `Secret Manager Secret Accessor` on each referenced secret.
- Anyone who needs to run the app locally must obtain the secrets out-of-band and add them to their `.env` file — there is no automatic local sync from Secret Manager.
