package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/NeevBhandari13/leetcoach/internal/session"
)

// noopStore is a Store where every method succeeds and returns zero values.
// Used for routing tests where only status codes matter, not business logic.
var noopStore Store = &mockStore{
	problemFn: func(_ context.Context) (string, error) { return "problem text", nil },
	createFn: func(_ context.Context, _, _ string) (*session.Session, error) {
		return &session.Session{SessionID: "test-id", State: session.IntroState}, nil
	},
	updateFn: func(_ context.Context, _ string, _ llm.Message) error { return nil },
	getFn: func(_ context.Context, _ string) (*session.Session, error) {
		return &session.Session{SessionID: "test-id", State: session.IntroState}, nil
	},
	replyFn: func(_ context.Context, _, _, _ string) (string, error) {
		return `{"reply":"ok","current_state":"intro"}`, nil
	},
	stateFn: func(_ context.Context, _ string, _ session.State) error { return nil },
}

func TestNewRouter(t *testing.T) {
	t.Run("returns non-nil engine", func(t *testing.T) {
		if NewRouter(noopStore) == nil {
			t.Fatal("expected non-nil router, got nil")
		}
	})
}

func TestSetupRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       any
		wantStatus int
	}{
		{
			name:       "POST /start is registered and returns 200",
			method:     http.MethodPost,
			path:       "/start",
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET /sessions/:id is registered and returns 200",
			method:     http.MethodGet,
			path:       "/sessions/test-id",
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST /sessions/:id/reply is registered and returns 200",
			method:     http.MethodPost,
			path:       "/sessions/test-id/reply",
			body:       replyRequest{Message: "hello"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET /start returns 405 method not allowed",
			method:     http.MethodGet,
			path:       "/start",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "unknown path returns 404",
			method:     http.MethodPost,
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter(noopStore)

			var bodyBytes []byte
			if tt.body != nil {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d (body: %s)", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}
