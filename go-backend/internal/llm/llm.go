package llm

import (
	"context"
)

type Message struct {
	Role    string // assistant or user
	Content string // message content
}

type Client interface {
	// system is the system prompt, messages is the conversation
	Send(ctx context.Context, system string, messages []Message) (string, error)
}
