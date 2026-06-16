# LeetCoach — Progress

> **For Claude Code sessions:** Read this file at the start of every session to understand where the project stands. Update the relevant sections as work is completed or new decisions are made. Keep entries concise — one or two lines per item is enough.

---

## What has been built

### Backend (Go + Gin + PostgreSQL)

- **Database layer** (`internal/db`)
  - `db.Open` / `db.Migrate` / `db.Seed` wired into startup
  - 4 migrations: `sessions`, `messages`, `problems` tables + `state`/`problem_text` columns on sessions
  - Problems seeded from embedded `seed/problems.json` (idempotent on conflict)

- **LLM layer** (`internal/llm`)
  - `Client` interface with Anthropic implementation
  - `Send` — standard message turn (claude-sonnet-4-6, max 1024 tokens)
  - TLS 1.2 forced on the HTTP client to work around a Go TLS 1.3 / Anthropic API incompatibility

- **Session store** (`internal/session`)
  - Full PostgreSQL-backed `SessionStore` (no in-memory state)
  - `GetRandomProblemText` — picks a random problem from DB and formats it for the system prompt
  - `CreateSession` — inserts a new session row; pass `""` as sessionID to let Postgres generate the UUID
  - `GetSession` — loads session + full ordered message history
  - `GetState` / `SetState` / `GetProblemText` / `UpdateChatHistory`
  - `Reply` — full LLM turn: insert user msg → fetch history → call LLM → insert assistant reply → return raw string

- **Prompts** (`internal/prompts`)
  - `GetSystemPrompt(state, problemText, code)` — base instructions + state-specific developer prompt + optional candidate code appended at the end
  - State prompts built once at package init time (`var statePrompts = ...` using `fmt.Sprintf`)

- **API** (`internal/api`)
  - `Store` interface decouples handlers from concrete `*session.SessionStore` (enables mock testing)
  - `POST /start` — picks random problem, creates session, returns hardcoded welcome + session ID
  - `GET /sessions/:id` — returns full session (state, problem_text, chat_history)
  - `POST /sessions/:id/reply` — one LLM turn, parses `{"reply","current_state"}` JSON, advances state machine
  - CORS configured for `http://localhost:3000`

- **State machine**
  ```
  intro → present_problem → clarify → initial_solution → optimisation → wrap_up
  ```
  The LLM drives transitions by returning `current_state` in its JSON response. `ReplyHandler` persists it via `SetState`.

- **Tests**
  - Handler unit tests with `mockStore` (injectable function fields, no DB)
  - Router smoke tests with `noopStore`
  - Prompt content tests in `internal/prompts`

### Frontend (Next.js, Pages Router, TypeScript)

- **`pages/index.tsx`** — landing page with "Start Interview" button
- **`pages/chat.tsx`** — main interview UI:
  - Split-pane layout: resizable code editor (left) + chat panel (right)
  - Draggable divider (mouse drag, clamped 20–80%)
  - Code editor content sent with every reply so the LLM can see candidate's code
- **`types/types.tsx`** — shared TypeScript interfaces for all API shapes
- Session flow: `POST /start` → redirect to `/chat?sessionID=...&initialText=...` (welcome message passed in URL to avoid extra round-trip)

### Infrastructure

- **`Makefile`** — `make dev` starts both services in parallel; Ctrl-C kills both
- **`CLAUDE.md`** — architecture reference and commands for Claude Code sessions

---

## What is not yet done / known gaps

### High priority

- [ ] **Frontend UI redesign** — the user started describing a new UI design in the last session but the message was cut off before the design was shown. This is the most immediate pending task. Ask the user to share the design if they haven't already.
- [ ] **Error handling in the frontend** — network errors and non-2xx responses are only `console.error`'d; no user-facing error state exists
- [ ] **Loading state during LLM reply** — no spinner or disabled input while waiting for the backend to respond
- [ ] **`internal/chat` cleanup** — `chat/chat.go` is a legacy wrapper around `llm.Client.Send` that is no longer wired into anything. It can be deleted.

### Medium priority

- [ ] **LLM JSON retry / fallback** — if the LLM returns malformed JSON, `ReplyHandler` immediately returns 500. A retry or graceful error message to the user would make the app more robust.
- [ ] **State validation in `SetState`** — any string the LLM returns as `current_state` is persisted verbatim. Add an allowlist check.
- [ ] **Session expiry / cleanup** — sessions accumulate in the DB indefinitely; no TTL or delete mechanism
- [ ] **Frontend ADRs** — user requested an `frontend/decisions/` folder with ADRs (was in progress when this file was created)
- [ ] **Backend ADRs** — user requested a `backend/decisions/` folder with ADRs (was in progress when this file was created)
- [ ] **More seed problems** — currently only a small set of problems in `seed/problems.json`

### Low priority / nice to have

- [ ] **Syntax highlighting in the code editor** — currently a plain `<textarea>`; a library like CodeMirror or Monaco would improve the experience
- [ ] **Markdown rendering in chat** — LLM replies are rendered as plain text; code blocks and formatting are lost
- [ ] **Session history / resume** — users can't return to a previous session; `GET /sessions/:id` exists but nothing in the UI uses it
- [ ] **Auth** — no authentication; any caller can start sessions and read any session by ID

---

## Architecture diagram

```
┌─────────────────────────────────────────────────────────────────┐
│  Browser (localhost:3000)                                        │
│                                                                  │
│  pages/index.tsx  ──POST /start──►  pages/chat.tsx              │
│                                          │                       │
│                    ◄──sessionID,msg──────┘                       │
│                                          │                       │
│                    POST /sessions/:id/reply (message + code)     │
│                    GET  /sessions/:id                            │
└──────────────────────────────┬──────────────────────────────────┘
                               │ HTTP (localhost:8080)
┌──────────────────────────────▼──────────────────────────────────┐
│  Go backend (Gin)                                                │
│                                                                  │
│  api/handlers.go                                                 │
│    StartInterviewHandler  ──► session.GetRandomProblemText       │
│                               session.CreateSession              │
│                               session.UpdateChatHistory          │
│                                                                  │
│    ReplyHandler           ──► session.GetSession                 │
│                               prompts.GetSystemPrompt            │
│                               session.Reply ──► llm.Client.Send  │
│                               session.SetState                   │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│  PostgreSQL                                                      │
│  sessions │ messages │ problems                                  │
└─────────────────────────────────────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│  Anthropic API                                                   │
│  claude-sonnet-4-6  (interview turns)                            │
└─────────────────────────────────────────────────────────────────┘
```

---

## Session log

| Date | What happened |
|------|---------------|
| 2026-06-16 | Initial full-stack build: DB + migrations + session store + LLM layer + API handlers + state machine + frontend + Makefile. E2E test passed. CLAUDE.md written. |
| 2026-06-16 | `progress.md` created. ADRs requested by user (frontend/decisions + backend/decisions) — not yet written. Frontend UI redesign pending (user's message was cut off). |
