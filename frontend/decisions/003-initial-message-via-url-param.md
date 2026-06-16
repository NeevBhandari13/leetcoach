# ADR 003: Initial Message Passed via URL Query Parameter

**Status:** Accepted

## Context

After `POST /start`, the backend returns `{ session_id, message }`. The chat page needs both the session ID (to make subsequent API calls) and the welcome message (to display as the first chat message). We needed to decide how to get this data from the landing page to the chat page.

## Decision

Pass both values as URL query parameters when navigating from `/` to `/chat`:

```
router.push({
    pathname: '/chat',
    query: {
        sessionID: data.session_id,
        initialText: data.message,
    },
})
```

The chat page reads them via `router.query` on mount.

## Alternatives Considered

| Option | Why rejected |
|--------|-------------|
| `GET /sessions/:id` on chat page mount | Extra round-trip; response also contains `problem_text`, which must stay hidden from the client |
| `localStorage` | Survives the navigation but requires cleanup; breaks if the user opens two tabs |
| React Context / global state | Doesn't survive a full page navigation in Pages Router |
| Redux / Zustand | Significant overhead for passing two strings |

## Why Not Fetch from the API?

```mermaid
flowchart TD
    subgraph Chosen approach
        A[POST /start] -->|{session_id, message}| B[router.push with query params]
        B --> C[chat page reads router.query]
    end

    subgraph Rejected approach
        D[POST /start] -->|{session_id}| E[router.push]
        E --> F[chat page mounts]
        F -->|GET /sessions/:id| G[backend]
        G -->|{session_id, state, problem_text, chat_history}| F
        F -->|problem_text exposed to browser!| H[Security issue]
    end
```

`GET /sessions/:id` returns `problem_text` in the response body. Since the problem should only be revealed by the LLM during the interview, fetching the session on the client would expose it in the network tab before the interview begins.

## Consequences

- The welcome message is visible in the browser's address bar. This is acceptable — it is not sensitive information.
- If the user navigates directly to `/chat` without a `sessionID` query param, `handleSend` will call `POST /sessions/undefined/reply`, which the backend returns 404 for. A future improvement would be to redirect to `/` if `sessionID` is missing.
- Long welcome messages may produce a long URL, but in practice the message is a short fixed string.