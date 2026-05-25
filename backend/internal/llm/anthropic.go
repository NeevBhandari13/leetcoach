package llm

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicClient struct {
	client anthropic.Client
	model  anthropic.Model
}

func NewAnthropicClient(apiKey string, anthropicModel anthropic.Model) *AnthropicClient {
	// Go's TLS 1.3 is incompatible with api.anthropic.com; force TLS 1.2.
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{MaxVersion: tls.VersionTLS12},
			ForceAttemptHTTP2: false,
		},
	}
	return &AnthropicClient{
		client: anthropic.NewClient(
			option.WithAPIKey(apiKey),
			option.WithHTTPClient(httpClient),
		),
		model: anthropicModel,
	}
}

func toAnthropicMessages(messages []Message) []anthropic.MessageParam {
	anthropicMessages := make([]anthropic.MessageParam, len(messages))
	for idx, message := range messages {
		anthropicMessages[idx] = toAnthropicMessage(message)
	}

	return anthropicMessages
}

func toAnthropicMessage(message Message) anthropic.MessageParam {
	block := anthropic.NewTextBlock(message.Content)
	if message.Role == "assistant" {
		return anthropic.NewAssistantMessage(block)
	} else {
		return anthropic.NewUserMessage(block)
	}
}

func (c *AnthropicClient) Send(ctx context.Context, system string, messages []Message) (string, error) {
	response, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages:  toAnthropicMessages(messages),
		System: []anthropic.TextBlockParam{{
			Text: system,
		}},
		Model: c.model,
	})

	if err != nil {
		return "", err
	}

	for _, block := range response.Content {
		switch v := block.AsAny().(type) {
		case anthropic.TextBlock:
			return v.Text, nil
		}
	}

	return "", fmt.Errorf("no text in response")

}
