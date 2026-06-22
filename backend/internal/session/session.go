package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/NeevBhandari13/leetcoach/internal/llm"
)

type State string

const (
	NilState            State = ""
	IntroState          State = "intro"
	PresentProblemState State = "present_problem"
	ClarifyState        State = "clarify"
	InitialSolutionState State = "initial_solution"
	OptimisationState   State = "optimisation"
	WrapUpState         State = "wrap_up"
)

var ErrNotFound = errors.New("session not found")

type Session struct {
	SessionID   string
	State       State
	ChatHistory []llm.Message
	ProblemText string
}

type SessionStore struct {
	db     *sql.DB
	client llm.Client
}

func NewSessionStore(db *sql.DB, client llm.Client) *SessionStore {
	return &SessionStore{db: db, client: client}
}

// GetRandomProblemText picks a random problem from the problems table and
// returns it formatted for inclusion in the LLM system prompt.
func (s *SessionStore) GetRandomProblemText(ctx context.Context) (string, error) {
	var title, difficulty, description, constraints string
	err := s.db.QueryRowContext(ctx,
		`SELECT title, difficulty, description, COALESCE(constraints, '')
		   FROM problems
		  ORDER BY RANDOM()
		  LIMIT 1`,
	).Scan(&title, &difficulty, &description, &constraints)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Title: %s\nDifficulty: %s\n\n%s\n\nConstraints: %s",
		title, difficulty, description, constraints), nil
}

// CreateSession inserts a new session row with the given problem text and
// returns the created session. Pass an empty sessionID to let Postgres generate
// the UUID.
func (s *SessionStore) CreateSession(ctx context.Context, sessionID, problemText string) (*Session, error) {
	var row *sql.Row
	if sessionID == "" {
		row = s.db.QueryRowContext(ctx,
			`INSERT INTO sessions (problem_text) VALUES ($1) RETURNING id, state, problem_text`,
			problemText,
		)
	} else {
		row = s.db.QueryRowContext(ctx,
			`INSERT INTO sessions (id, problem_text) VALUES ($1, $2) RETURNING id, state, problem_text`,
			sessionID, problemText,
		)
	}

	var sess Session
	if err := row.Scan(&sess.SessionID, &sess.State, &sess.ProblemText); err != nil {
		return nil, err
	}
	sess.ChatHistory = []llm.Message{}
	return &sess, nil
}

func (s *SessionStore) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	var sess Session
	err := s.db.QueryRowContext(ctx,
		`SELECT id, state, problem_text FROM sessions WHERE id = $1`,
		sessionID,
	).Scan(&sess.SessionID, &sess.State, &sess.ProblemText)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT role, content FROM messages WHERE session_id = $1 ORDER BY created_at, id`,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m llm.Message
		if err := rows.Scan(&m.Role, &m.Content); err != nil {
			return nil, err
		}
		sess.ChatHistory = append(sess.ChatHistory, m)
	}
	return &sess, rows.Err()
}

func (s *SessionStore) GetState(ctx context.Context, sessionID string) (State, error) {
	var state State
	err := s.db.QueryRowContext(ctx,
		`SELECT state FROM sessions WHERE id = $1`,
		sessionID,
	).Scan(&state)
	if errors.Is(err, sql.ErrNoRows) {
		return NilState, ErrNotFound
	}
	return state, err
}

func (s *SessionStore) SetState(ctx context.Context, sessionID string, state State) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE sessions SET state = $1 WHERE id = $2`,
		state, sessionID,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("session with ID %s not found", sessionID)
	}
	return nil
}

func (s *SessionStore) GetProblemText(ctx context.Context, sessionID string) (string, error) {
	var text string
	err := s.db.QueryRowContext(ctx,
		`SELECT problem_text FROM sessions WHERE id = $1`,
		sessionID,
	).Scan(&text)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return text, err
}

func (s *SessionStore) UpdateChatHistory(ctx context.Context, sessionID string, message llm.Message) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO messages (session_id, role, content) VALUES ($1, $2, $3)`,
		sessionID, message.Role, message.Content,
	)
	return err
}

func (s *SessionStore) Reply(ctx context.Context, sessionID, system, userMessage string) (string, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM sessions WHERE id = $1)`, sessionID,
	).Scan(&exists)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrNotFound
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO messages (session_id, role, content) VALUES ($1, $2, $3)`,
		sessionID, "user", userMessage,
	)
	if err != nil {
		return "", err
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT role, content FROM messages WHERE session_id = $1 ORDER BY created_at, id`,
		sessionID,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var messages []llm.Message
	for rows.Next() {
		var m llm.Message
		if err := rows.Scan(&m.Role, &m.Content); err != nil {
			return "", err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	return s.client.Send(ctx, system, messages)
}
