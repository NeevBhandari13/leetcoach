package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/NeevBhandari13/leetcoach/internal/session"
	"github.com/gin-gonic/gin"
)

// mockStore implements the Store interface with injectable functions so each
// test case can control exactly what the store returns.
type mockStore struct {
	problemFn func(ctx context.Context) (string, error)
	createFn  func(ctx context.Context, sessionID, problemText string) (*session.Session, error)
	updateFn  func(ctx context.Context, sessionID string, message llm.Message) error
	getFn     func(ctx context.Context, sessionID string) (*session.Session, error)
	replyFn   func(ctx context.Context, sessionID, system, userMessage string) (string, error)
	stateFn   func(ctx context.Context, sessionID string, state session.State) error
}

func (m *mockStore) GetRandomProblemText(ctx context.Context) (string, error) {
	return m.problemFn(ctx)
}
func (m *mockStore) CreateSession(ctx context.Context, sessionID, problemText string) (*session.Session, error) {
	return m.createFn(ctx, sessionID, problemText)
}
func (m *mockStore) UpdateChatHistory(ctx context.Context, sessionID string, msg llm.Message) error {
	return m.updateFn(ctx, sessionID, msg)
}
func (m *mockStore) GetSession(ctx context.Context, sessionID string) (*session.Session, error) {
	return m.getFn(ctx, sessionID)
}
func (m *mockStore) Reply(ctx context.Context, sessionID, system, userMessage string) (string, error) {
	return m.replyFn(ctx, sessionID, system, userMessage)
}
func (m *mockStore) SetState(ctx context.Context, sessionID string, state session.State) error {
	return m.stateFn(ctx, sessionID, state)
}

// successStore returns a store where every method succeeds with sensible defaults.
func successStore() *mockStore {
	return &mockStore{
		problemFn: func(_ context.Context) (string, error) { return "Given an array...", nil },
		createFn: func(_ context.Context, _, _ string) (*session.Session, error) {
			return &session.Session{SessionID: "test-id", State: session.IntroState}, nil
		},
		updateFn: func(_ context.Context, _ string, _ llm.Message) error { return nil },
		getFn: func(_ context.Context, _ string) (*session.Session, error) {
			return &session.Session{
				SessionID:   "test-id",
				State:       session.IntroState,
				ProblemText: "Given an array...",
			}, nil
		},
		replyFn: func(_ context.Context, _, _, _ string) (string, error) {
			return `{"reply":"Great question!","current_state":"present_problem"}`, nil
		},
		stateFn: func(_ context.Context, _ string, _ session.State) error { return nil },
	}
}

// --- StartInterviewHandler ---

func TestStartInterviewHandler(t *testing.T) {
	tests := []struct {
		name       string
		store      Store
		wantStatus int
		wantID     string
		wantMsg    string
	}{
		{
			name:       "returns 200 with session ID and welcome message",
			store:      successStore(),
			wantStatus: http.StatusOK,
			wantID:     "test-id",
			wantMsg:    "Hello! Welcome to LeetCoach!",
		},
		{
			name: "returns 500 when CreateSession fails",
			store: &mockStore{
				problemFn: func(_ context.Context) (string, error) { return "problem", nil },
				createFn: func(_ context.Context, _, _ string) (*session.Session, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "returns 500 when UpdateChatHistory fails",
			store: &mockStore{
				problemFn: func(_ context.Context) (string, error) { return "problem", nil },
				createFn: func(_ context.Context, _, _ string) (*session.Session, error) {
					return &session.Session{SessionID: "test-id"}, nil
				},
				updateFn: func(_ context.Context, _ string, _ llm.Message) error {
					return errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.POST("/start", StartInterviewHandler(tt.store))

			req := httptest.NewRequest(http.MethodPost, "/start", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d (body: %s)", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantID != "" {
				var resp startInterviewResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.SessionID != tt.wantID {
					t.Errorf("expected session_id %q, got %q", tt.wantID, resp.SessionID)
				}
				if len(resp.Message) == 0 {
					t.Error("expected non-empty message")
				}
			}
		})
	}
}

// --- ReplyHandler ---

func TestReplyHandler(t *testing.T) {
	validBody := replyRequest{Message: "I am ready"}

	tests := []struct {
		name       string
		store      Store
		body       any
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "returns 200 with parsed reply on success",
			store:      successStore(),
			body:       validBody,
			wantStatus: http.StatusOK,
			wantMsg:    "Great question!",
		},
		{
			name:       "returns 400 on invalid JSON body",
			store:      successStore(),
			body:       "not-json{{{",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "returns 404 when session not found",
			store: &mockStore{
				getFn: func(_ context.Context, _ string) (*session.Session, error) {
					return nil, session.ErrNotFound
				},
			},
			body:       validBody,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "returns 500 when Reply fails",
			store: &mockStore{
				getFn: func(_ context.Context, _ string) (*session.Session, error) {
					return &session.Session{State: session.IntroState}, nil
				},
				replyFn: func(_ context.Context, _, _, _ string) (string, error) {
					return "", errors.New("llm error")
				},
			},
			body:       validBody,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "returns 500 when LLM returns malformed JSON",
			store: &mockStore{
				getFn: func(_ context.Context, _ string) (*session.Session, error) {
					return &session.Session{State: session.IntroState}, nil
				},
				replyFn: func(_ context.Context, _, _, _ string) (string, error) {
					return "not valid json", nil
				},
			},
			body:       validBody,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.POST("/sessions/:id/reply", ReplyHandler(tt.store))

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/sessions/test-id/reply", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d (body: %s)", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantMsg != "" {
				var resp replyResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.Message != tt.wantMsg {
					t.Errorf("expected message %q, got %q", tt.wantMsg, resp.Message)
				}
			}
		})
	}
}
