#!/usr/bin/env bash
# Analyze the last git commit and create or update Architecture Decision Records.
# Called automatically by .git/hooks/post-commit. Safe to run manually too.

set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

# Nothing to do on the very first commit (no parent)
if ! git rev-parse HEAD~1 >/dev/null 2>&1; then
  exit 0
fi

LOG="$REPO_ROOT/.claude/adr-update.log"

run_analysis() {
  DIFF=$(git diff HEAD~1 HEAD | head -c 10000)
  COMMIT_MSG=$(git log -1 --pretty=format:"%s%n%n%b" | head -c 500)
  COMMIT_HASH=$(git log -1 --pretty=format:"%h")
  COMMIT_DATE=$(git log -1 --pretty=format:"%as")

  # Highest existing ADR number in a folder, returns next as zero-padded string
  next_num() {
    local dir="$1"
    local max=0
    for f in "$dir"/[0-9][0-9][0-9]-*.md; do
      [[ -f "$f" ]] || continue
      n=$(basename "$f" | grep -o '^[0-9]*' | sed 's/^0*//')
      [[ -n "$n" ]] && (( n > max )) && max=$n
    done
    printf "%03d" $((max + 1))
  }

  list_adrs() {
    local dir="$1"
    local found=0
    for f in "$dir"/[0-9][0-9][0-9]-*.md; do
      [[ -f "$f" ]] || continue
      found=1
      title=$(head -1 "$f" | sed 's/^#[[:space:]]*//')
      status=$(grep '^\*\*Status' "$f" 2>/dev/null | head -1 | sed 's/\*\*Status:\*\* //' || echo "Unknown")
      echo "  - $(basename "$f"): $title [$status]"
    done
    [[ $found -eq 0 ]] && echo "  (none yet)"
  }

  NEXT_BACKEND=$(next_num "backend/decisions")
  NEXT_FRONTEND=$(next_num "frontend/decisions")
  EXISTING_BACKEND=$(list_adrs "backend/decisions")
  EXISTING_FRONTEND=$(list_adrs "frontend/decisions")

  PROMPT="You maintain Architecture Decision Records (ADRs) for LeetCoach — an AI-powered mock coding interview app (Go/Gin backend, Next.js frontend, PostgreSQL, Anthropic API).

## Task

Analyze the commit below. Decide whether it introduces or meaningfully changes an architectural decision worth recording.

**Create or update an ADR only for genuine architectural choices**: a new library, data-storage decision, API contract, design pattern, or significant tradeoff. Skip cosmetic changes, doc updates, test tweaks, comment edits, and minor refactors with no architectural significance.

If nothing warrants an ADR, output exactly the text: No ADRs needed.

## Commit

Hash: ${COMMIT_HASH}
Date: ${COMMIT_DATE}
Message: ${COMMIT_MSG}

## Diff

\`\`\`diff
${DIFF}
\`\`\`

## Existing ADRs

backend/decisions/:
${EXISTING_BACKEND}

frontend/decisions/:
${EXISTING_FRONTEND}

## ADR file format

\`\`\`markdown
# ADR <NNN>: <Title>

**Status:** Accepted

## Context

<What situation or constraint forced this decision?>

## Decision

<What was decided?>

## Alternatives Considered

| Option | Why rejected |
|--------|-------------|
| <alt> | <reason> |

## Consequences

<Trade-offs and follow-on effects.>
\`\`\`

## File locations

- New backend ADR → backend/decisions/<NNN>-<kebab-title>.md  (next available number: ${NEXT_BACKEND})
- New frontend ADR → frontend/decisions/<NNN>-<kebab-title>.md (next available number: ${NEXT_FRONTEND})
- Updating existing → Read it first, then Edit it in place.

## Index maintenance

After writing any ADR, rebuild the index for that folder. Read all existing \`.md\` files in the folder (excluding index.md) to get current titles, then overwrite the index file:

\`\`\`markdown
# <Backend|Frontend> Architecture Decision Records

| # | Decision | Status |
|---|----------|--------|
| [001](001-filename.md) | Title | Accepted |
\`\`\`

## Steps

1. Read and understand the diff.
2. Identify any architectural decisions.
3. For each decision:
   - New: create the ADR file at the correct path.
   - Significant change to existing: Read the file, then Edit it.
   - No decisions: output 'No ADRs needed.' and stop.
4. For every folder you touched, read all existing ADR files, then overwrite that folder's index.md."

  claude -p "$PROMPT" --allowedTools "Read,Write,Edit"
}

(
  echo ""
  echo "=== $(date -Iseconds) | commit $(git log -1 --pretty=format:"%h") ==="
  run_analysis
) >> "$LOG" 2>&1 &

disown $!
echo "[adr] analyzing commit in background → .claude/adr-update.log"
