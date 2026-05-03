package main

import (
	"os"

	"github.com/NeevBhandari13/leetcoach/internal/api"
	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := NewLLMClient()                   // llm client interface
	chatService := chat.NewChatService(client) // chat service with the reply function
	router := api.NewRouter(chatService)       // creates new router
	router.Run(":8080")                        // router on port 8080

}

func NewLLMClient() llm.Client {
	provider := os.Getenv("LLM_PROVIDER")

	switch provider {
	case "anthropic":
		return llm.NewAnthropicClient(os.Getenv("ANTHROPIC_API_KEY"), anthropic.Model(os.Getenv("LLM_MODEL")))
	default:
		panic("Invalid LLM Provider")

	}
}
