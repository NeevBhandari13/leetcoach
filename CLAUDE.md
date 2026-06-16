# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Start here

Read `.claude/progress.md` at the beginning of every session — it tracks what has been built, what is not yet done, and a log of what happened in past sessions. Update it as work is completed or new gaps are discovered.

## What this project is

LeetCoach is an AI-powered mock coding interview app. The backend runs a stateful interview session via an LLM (Anthropic Claude), progressing the candidate through a defined state machine. The frontend is a split-pane UI: a code editor on the left and a chat window on the right.

## Commands

```bash
# Start both services (Ctrl-C kills both)
make dev

# Run backend only
make backend          # go run ./cmd/main.go from backend/

# Run frontend only
make frontend         # npm run dev from frontend/

# Run all backend tests
make test             # go test ./... from backend/

# Run a single backend test
cd backend && go test ./internal/api/... -run TestReplyHandler

# Run frontend type-check
cd frontend && npx tsc --noEmit
```

## Environment

The root `.env` file is loaded by the backend at startup (one level up from `backend/`):

```
ANTHROPIC_API_KEY=...
LLM_PROVIDER=anthropic
LLM_MODEL=claude-sonnet-4-6       # used for interview replies
REASONING_MODEL=claude-opus-4-7   # used for post-interview review
DSN=postgres://...
```

The frontend reads `frontend/.env.local`:
```
NEXT_PUBLIC_BACKEND_API_BASE_URL=http://localhost:8080
```

## Architecture

### Backend (`backend/`)

Entry point: `cmd/main.go` — opens the DB, runs migrations, seeds problems, constructs the store and router, listens on `:8080`.

**Startup sequence:** `db.Open` → `db.Migrate` (golang-migrate, embedded SQL files) → `db.Seed` (idempotent upsert from `seed/problems.json`) → `session.NewSessionStore` → `api.NewRouter`.

**Packages:**

- `internal/db` — DB open/migrate/seed. Migrations and seed data are embedded at compile time via `//go:embed`.
- `internal/llm` — `Client` and `ReviewClient` interfaces, Anthropic implementation. Interview replies use `claude-sonnet-4-6`; end-of-session review uses `claude-opus-4-7` (extended thinking).
- `internal/session` — `SessionStore` backed by Postgres. Owns all SQL. `Reply` inserts user message, fetches full history, calls `client.Send`, inserts assistant reply. `GenerateReview` fetches full session and calls `reviewClient.SendReview`.
- `internal/prompts` — `GetSystemPrompt(state, problemText, code)` builds the system prompt: fixed base + state-specific developer prompt + optional candidate code appended at the end.
- `internal/api` — Gin handlers and router. Defines a `Store` interface (subset of `*session.SessionStore`) so handlers are testable without a real DB.

**API endpoints:**

| Method | Path | Handler |
|--------|------|---------|
| POST | `/start` | `StartInterviewHandler` — picks random problem, creates session, returns hardcoded welcome message |
| GET | `/sessions/:id` | `GetSessionHandler` — returns full session (state, problem_text, chat_history) |
| POST | `/sessions/:id/reply` | `ReplyHandler` — runs one LLM turn, advances state machine |
| POST | `/sessions/:id/review` | `ReviewHandler` — generates end-of-interview review |

CORS is configured to allow `http://localhost:3000`.

### State machine

States (in `internal/session/session.go`): `intro` → `present_problem` → `clarify` → `initial_solution` → `optimisation` → `wrap_up`.

The LLM is instructed to return JSON `{"reply": "...", "current_state": "..."}` on every turn. `ReplyHandler` parses this and calls `store.SetState` to persist the transition. The system prompt for each turn is rebuilt from the current state so the LLM knows valid transitions.

### Frontend (`frontend/`)

Next.js with Pages Router. Two pages:
- `pages/index.tsx` — landing page with a start button. `POST /start` → redirects to `/chat?sessionID=...&initialText=...`
- `pages/chat.tsx` — split-pane interview UI. Left: `<textarea>` code editor (value sent with every reply). Right: `ChatWindow` + `ChatInput`. When `currentState === 'wrap_up'`, automatically `POST /sessions/:id/review` and shows `ReviewModal`.

### Database schema

Four migrations (auto-applied on startup):
1. `sessions` table — UUID primary key, `state TEXT DEFAULT 'intro'`, `problem_text TEXT`
2. `messages` table — `session_id` FK, `role`, `content`, ordered by `created_at, id`
3. `problems` table — `slug` (unique), `title`, `difficulty`, `description`, `examples` (JSONB), `constraints`, `topics` (TEXT[])
4. Adds `state` and `problem_text` columns to `sessions`

### Testing

Handler tests use `mockStore` (in `session_handlers_test.go`) with injectable function fields — no DB needed. Router tests use `noopStore`. The `Store` interface in `internal/api/handlers.go` is the seam that makes this possible.

`handlers_test.go` contains only an `init()` that sets Gin to test mode — it is not a test file itself, just shared setup for the package.

**Extending the `Store` interface:** if you add a new method to `session.SessionStore` and call it from a handler, you must add it in three places: (1) `Store` interface in `handlers.go`, (2) `mockStore` in `session_handlers_test.go`, (3) `noopStore` in `routers_test.go`. Missing any one causes a compile error.

## Key design decisions

**Problem text is hidden from the frontend.** It lives only in the DB and is injected into the LLM system prompt each turn. The frontend never receives it — the LLM reveals the problem naturally during the interview flow.

**`UpdateChatHistory` vs `Reply`:** `UpdateChatHistory` is a raw INSERT used exactly once in `StartInterviewHandler` to persist the hardcoded welcome message without calling the LLM. All subsequent turns go through `Reply`, which does the full round-trip: insert user message → fetch full history → call LLM → insert assistant reply → return raw LLM string.

**Code editor content is not stored in `messages`.** The `replyRequest` body carries both `message` (user chat text, stored in DB) and `code` (full editor content, not stored). The code is appended to the system prompt each turn by `GetSystemPrompt` so the LLM always sees the latest version, but it is never written to the database.

**Initial message passed via URL query param.** After `POST /start`, the frontend passes the welcome message to `/chat` as `?initialText=...` rather than fetching it from `GET /sessions/:id`. This avoids an extra round-trip and keeps `problem_text` from ever reaching the client (it would be visible in the session response).

**`statePrompts` map is built at package init time.** In `internal/prompts/system.go`, `var statePrompts = map[session.State]string{...}` uses `fmt.Sprintf` calls that reference `stateInstructions`. Go initialises package-level vars in declaration order within a file, so `stateInstructions` must be declared before `statePrompts` in that file — keep this order.

**Review uses a separate LLM client (`ReviewClient`).** The review call (`GenerateReview`) uses `REASONING_MODEL` (claude-opus-4-7) rather than the interview model. Both are Anthropic clients but the interface is split (`Client`/`ReviewClient`) to allow different models or providers per use case.

**`reviewFetchedRef` in `chat.tsx`** is a `useRef` (not state) that prevents the review from being fetched twice when `currentState` becomes `wrap_up`. Because the `useEffect` that watches `currentState` may fire during re-renders, using a ref instead of state avoids triggering another render cycle.

## Known fragilities

- **LLM JSON contract has no retry.** `ReplyHandler` immediately returns 500 if the LLM response is not valid JSON matching `{"reply": "...", "current_state": "..."}`. There is no fallback or retry logic.
- **No state validation in `SetState`.** Any string the LLM returns as `current_state` will be persisted verbatim. Invalid state names are not rejected at the store layer.
- **`internal/chat` is a legacy package** (`chat/chat.go`) that wraps `llm.Client.Send` but is no longer wired into `main.go`. It can be deleted if it causes confusion.