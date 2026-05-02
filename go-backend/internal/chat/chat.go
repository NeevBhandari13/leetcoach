package chat

import (
	"context"

	"github.com/NeevBhandari13/leetcoach/internal/llm"
)

type ChatService struct {
	client llm.Client
}

func NewChatService(client llm.Client) *ChatService {
	return &ChatService{
		client: client,
	}
}

func (c ChatService) Reply(ctx context.Context, system string, messages []llm.Message) (string, error) {
	return c.client.Send(ctx, system, messages)
}
