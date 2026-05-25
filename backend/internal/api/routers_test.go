package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/NeevBhandari13/leetcoach/internal/llm"
)

// noopFn is a mock send function that always succeeds — used when the LLM
// response doesn't matter for the test, only the routing behaviour does.
func noopFn(_ context.Context, _ string, _ []llm.Message) (string, error) {
	return "ok", nil
}

func newNoopService() *chat.ChatService {
	return chat.NewChatService(&mockLLMClient{sendFn: noopFn})
}

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "returns non-nil engine"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter(newNoopService(), nil)
			if router == nil {
				t.Fatal("expected non-nil router, got nil")
			}
		})
	}
}

func TestSetupRoutes(t *testing.T) {
	validBody := ChatRequest{
		System:   "You are a technical interviewer.",
		Messages: []llm.Message{{Role: "user", Content: "I am ready"}},
	}

	tests := []struct {
		name       string
		method     string
		path       string
		body       any
		wantStatus int
	}{
		{
			name:       "POST /chat is registered and returns 200",
			method:     http.MethodPost,
			path:       "/chat",
			body:       validBody,
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET /chat returns 405 method not allowed",
			method:     http.MethodGet,
			path:       "/chat",
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
			router := NewRouter(newNoopService(), nil)

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
