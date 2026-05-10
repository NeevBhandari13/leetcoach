package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockLLMClient implements llm.Client so we can inject controlled responses.
// Defined here and reused by routers_test.go (same package).
type mockLLMClient struct {
	sendFn func(ctx context.Context, system string, messages []llm.Message) (string, error)
}

func (m *mockLLMClient) Send(ctx context.Context, system string, messages []llm.Message) (string, error) {
	return m.sendFn(ctx, system, messages)
}

// newServiceWithMock builds a ChatService backed by a controlled mock function.
func newServiceWithMock(fn func(ctx context.Context, system string, messages []llm.Message) (string, error)) *chat.ChatService {
	return chat.NewChatService(&mockLLMClient{sendFn: fn})
}

func TestChatHandler(t *testing.T) {
	successFn := func(ctx context.Context, system string, messages []llm.Message) (string, error) {
		return "Tell me about arrays.", nil
	}
	errorFn := func(ctx context.Context, system string, messages []llm.Message) (string, error) {
		return "", errors.New("LLM failure")
	}

	tests := []struct {
		name       string
		body       any
		mockFn     func(ctx context.Context, system string, messages []llm.Message) (string, error)
		wantStatus int
		wantMsg    string
	}{
		{
			name: "returns 200 with message on success",
			body: ChatRequest{
				System:   "You are a technical interviewer.",
				Messages: []llm.Message{{Role: "user", Content: "I am ready"}},
			},
			mockFn:     successFn,
			wantStatus: http.StatusOK,
			wantMsg:    "Tell me about arrays.",
		},
		{
			name:       "returns 400 on invalid JSON body",
			body:       "not valid json{{{",
			mockFn:     successFn,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "returns 500 when service returns error",
			body: ChatRequest{
				System:   "You are a technical interviewer.",
				Messages: []llm.Message{{Role: "user", Content: "Hi"}},
			},
			mockFn:     errorFn,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "passes system and messages to service",
			body: ChatRequest{
				System:   "custom system",
				Messages: []llm.Message{{Role: "user", Content: "question"}},
			},
			mockFn:     successFn,
			wantStatus: http.StatusOK,
			wantMsg:    "Tell me about arrays.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newServiceWithMock(tt.mockFn)

			r := gin.New()
			r.POST("/chat", ChatHandler(service))

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d (body: %s)", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantMsg != "" {
				var resp ChatResponse
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
