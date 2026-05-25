package chat

import (
	"context"
	"errors"
	"testing"

	"github.com/NeevBhandari13/leetcoach/internal/llm"
)

// mockLLMClient implements llm.Client without hitting the real API.
type mockLLMClient struct {
	sendFn func(ctx context.Context, system string, messages []llm.Message) (string, error)
}

func (m *mockLLMClient) Send(ctx context.Context, system string, messages []llm.Message) (string, error) {
	return m.sendFn(ctx, system, messages)
}

func TestNewChatService(t *testing.T) {
	tests := []struct {
		name   string
		client llm.Client
	}{
		{
			name:   "creates non-nil service with valid client",
			client: &mockLLMClient{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewChatService(tt.client)
			if service == nil {
				t.Fatal("expected non-nil service, got nil")
			}
		})
	}
}

func TestReply(t *testing.T) {
	tests := []struct {
		name       string
		system     string
		messages   []llm.Message
		mockReturn string
		mockErr    error
		wantResult string
		wantErr    bool
	}{
		{
			name:       "returns response from client",
			system:     "You are a technical interviewer.",
			messages:   []llm.Message{{Role: "user", Content: "I'm ready"}},
			mockReturn: "Tell me about linked lists.",
			wantResult: "Tell me about linked lists.",
		},
		{
			name:     "propagates client error",
			system:   "You are a technical interviewer.",
			messages: []llm.Message{{Role: "user", Content: "Hello"}},
			mockErr:  errors.New("API unavailable"),
			wantErr:  true,
		},
		{
			name:       "passes system prompt through unchanged",
			system:     "custom system prompt",
			messages:   []llm.Message{{Role: "user", Content: "hello"}},
			mockReturn: "ok",
			wantResult: "ok",
		},
		{
			name:   "handles multi-turn conversation",
			system: "You are a technical interviewer.",
			messages: []llm.Message{
				{Role: "user", Content: "hello"},
				{Role: "assistant", Content: "hi there"},
				{Role: "user", Content: "let's start"},
			},
			mockReturn: "Sure, what topic?",
			wantResult: "Sure, what topic?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotSystem string
			var gotMessages []llm.Message

			client := &mockLLMClient{
				sendFn: func(ctx context.Context, system string, messages []llm.Message) (string, error) {
					gotSystem = system
					gotMessages = messages
					return tt.mockReturn, tt.mockErr
				},
			}

			service := NewChatService(client)
			result, err := service.Reply(context.Background(), tt.system, tt.messages)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.wantResult {
				t.Errorf("expected result %q, got %q", tt.wantResult, result)
			}
			if gotSystem != tt.system {
				t.Errorf("expected system prompt %q forwarded to client, got %q", tt.system, gotSystem)
			}
			if len(gotMessages) != len(tt.messages) {
				t.Errorf("expected %d messages forwarded to client, got %d", len(tt.messages), len(gotMessages))
			}
		})
	}
}
