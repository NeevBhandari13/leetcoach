package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestToAnthropicMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    Message
		wantRole string
		wantText string
	}{
		{
			name:     "user role",
			input:    Message{Role: "user", Content: "what is a linked list?"},
			wantRole: "user",
			wantText: "what is a linked list?",
		},
		{
			name:     "assistant role",
			input:    Message{Role: "assistant", Content: "a linked list is a chain of nodes"},
			wantRole: "assistant",
			wantText: "a linked list is a chain of nodes",
		},
		{
			name:     "unknown role defaults to user",
			input:    Message{Role: "system", Content: "some content"},
			wantRole: "user",
			wantText: "some content",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			param := toAnthropicMessage(test.input)

			data, err := json.Marshal(param)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}

			var result map[string]any
			json.Unmarshal(data, &result)

			if result["role"] != test.wantRole {
				t.Errorf("expected role %q, got %v", test.wantRole, result["role"])
			}

			content := result["content"].([]any)
			block := content[0].(map[string]any)
			if block["text"] != test.wantText {
				t.Errorf("expected text %q, got %v", test.wantText, block["text"])
			}
		})
	}
}

func TestToAnthropicMessages(t *testing.T) {
	tests := []struct {
		name      string
		input     []Message
		wantLen   int
		wantRoles []string
	}{
		{
			name:      "empty slice",
			input:     []Message{},
			wantLen:   0,
			wantRoles: []string{},
		},
		{
			name:      "single user message",
			input:     []Message{{Role: "user", Content: "hello"}},
			wantLen:   1,
			wantRoles: []string{"user"},
		},
		{
			name: "preserves order and roles",
			input: []Message{
				{Role: "user", Content: "first"},
				{Role: "assistant", Content: "second"},
				{Role: "user", Content: "third"},
			},
			wantLen:   3,
			wantRoles: []string{"user", "assistant", "user"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := toAnthropicMessages(tt.input)

			if len(params) != tt.wantLen {
				t.Fatalf("expected %d params, got %d", tt.wantLen, len(params))
			}

			for i, param := range params {
				data, _ := json.Marshal(param)
				var result map[string]any
				json.Unmarshal(data, &result)

				if result["role"] != tt.wantRoles[i] {
					t.Errorf("message %d: expected role %q, got %v", i, tt.wantRoles[i], result["role"])
				}
			}
		})
	}
}

func TestNewAnthropicClient(t *testing.T) {
	tests := []struct {
		name      string
		model     anthropic.Model
		wantModel anthropic.Model
	}{
		{
			name:      "sets sonnet model",
			model:     anthropic.ModelClaudeSonnet4_6,
			wantModel: anthropic.ModelClaudeSonnet4_6,
		},
		{
			name:      "sets opus model",
			model:     anthropic.ModelClaudeOpus4_7,
			wantModel: anthropic.ModelClaudeOpus4_7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewAnthropicClient("test-key", tt.model)

			if client == nil {
				t.Fatal("expected non-nil client")
			}
			if client.model != tt.wantModel {
				t.Errorf("expected model %v, got %v", tt.wantModel, client.model)
			}
		})
	}
}

func newMockServer(body string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
}

func newTestClient(serverURL string) *AnthropicClient {
	return &AnthropicClient{
		client: anthropic.NewClient(
			option.WithAPIKey("test-key"),
			option.WithBaseURL(serverURL),
		),
		model: anthropic.ModelClaudeSonnet4_6,
	}
}

func TestSend(t *testing.T) {
	tests := []struct {
		name       string
		serverBody string
		serverCode int
		wantResult string
		wantErrMsg string
	}{
		{
			name: "returns text on success",
			serverBody: `{
				"id": "msg_123",
				"type": "message",
				"role": "assistant",
				"content": [{"type": "text", "text": "Tell me about arrays."}],
				"model": "claude-sonnet-4-6",
				"stop_reason": "end_turn",
				"usage": {"input_tokens": 20, "output_tokens": 10}
			}`,
			serverCode: http.StatusOK,
			wantResult: "Tell me about arrays.",
		},
		{
			name:       "returns error on API failure",
			serverBody: `{"type":"error","error":{"type":"invalid_request_error","message":"bad request"}}`,
			serverCode: http.StatusBadRequest,
			wantErrMsg: "bad request",
		},
		{
			name: "returns error when no text block in response",
			serverBody: `{
				"id": "msg_123",
				"type": "message",
				"role": "assistant",
				"content": [],
				"model": "claude-sonnet-4-6",
				"stop_reason": "end_turn",
				"usage": {"input_tokens": 10, "output_tokens": 0}
			}`,
			serverCode: http.StatusOK,
			wantErrMsg: "no text in response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(tt.serverBody, tt.serverCode)
			defer server.Close()

			client := newTestClient(server.URL)
			result, err := client.Send(
				context.Background(),
				"You are a technical interviewer.",
				[]Message{{Role: "user", Content: "I'm ready."}},
			)

			if tt.wantErrMsg != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.wantResult {
				t.Errorf("expected %q, got %q", tt.wantResult, result)
			}
		})
	}
}
